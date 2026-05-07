package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetSignals() {
	signalsMu.Lock()
	defer signalsMu.Unlock()
	signals = map[string]signalRecord{}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	newMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content-type application/json, got %q", got)
	}

	var payload healthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if payload.Status != "ok" {
		t.Fatalf("expected status field 'ok', got %q", payload.Status)
	}
	if payload.Service != "nearby-signaling-server" {
		t.Fatalf("unexpected service field %q", payload.Service)
	}
	if payload.Timestamp == "" {
		t.Fatal("expected non-empty timestamp")
	}
}

func TestPostOnlyEndpoints_MethodHandling(t *testing.T) {
	resetSignals()
	mux := newMux()
	endpoints := []string{"/signal/offer", "/signal/answer"}

	for _, endpoint := range endpoints {
		t.Run(endpoint+"_GET_not_allowed", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
			}

			if got := rec.Header().Get("Allow"); got != http.MethodPost {
				t.Fatalf("expected Allow header %q, got %q", http.MethodPost, got)
			}
		})

		t.Run(endpoint+"_POST_accepted", func(t *testing.T) {
			body := `{"sessionId":"s1","sdp":"v=0","fromPeer":"peer-a"}`
			if endpoint == "/signal/answer" {
				offerReq := httptest.NewRequest(http.MethodPost, "/signal/offer", strings.NewReader(body))
				offerRec := httptest.NewRecorder()
				mux.ServeHTTP(offerRec, offerReq)
				if offerRec.Code != http.StatusAccepted {
					t.Fatalf("offer precondition status %d", offerRec.Code)
				}
			}
			req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(body))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusAccepted {
				t.Fatalf("expected status %d, got %d", http.StatusAccepted, rec.Code)
			}
		})
	}
}

func TestOfferAnswerFlow(t *testing.T) {
	resetSignals()
	mux := newMux()

	offer := `{"sessionId":"sess-1","sdp":"offer-sdp","fromPeer":"alice"}`
	offerReq := httptest.NewRequest(http.MethodPost, "/signal/offer", bytes.NewBufferString(offer))
	offerRec := httptest.NewRecorder()
	mux.ServeHTTP(offerRec, offerReq)
	if offerRec.Code != http.StatusAccepted {
		t.Fatalf("offer status = %d, want %d", offerRec.Code, http.StatusAccepted)
	}

	answer := `{"sessionId":"sess-1","sdp":"answer-sdp","fromPeer":"bob"}`
	answerReq := httptest.NewRequest(http.MethodPost, "/signal/answer", bytes.NewBufferString(answer))
	answerRec := httptest.NewRecorder()
	mux.ServeHTTP(answerRec, answerReq)
	if answerRec.Code != http.StatusAccepted {
		t.Fatalf("answer status = %d, want %d", answerRec.Code, http.StatusAccepted)
	}
}

func TestAnswerWithoutOfferFails(t *testing.T) {
	resetSignals()
	mux := newMux()

	answer := `{"sessionId":"missing","sdp":"answer-sdp","fromPeer":"bob"}`
	req := httptest.NewRequest(http.MethodPost, "/signal/answer", bytes.NewBufferString(answer))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestSessionPolling(t *testing.T) {
	resetSignals()
	mux := newMux()

	offer := `{"sessionId":"sess-poll","sdp":"offer-sdp","fromPeer":"alice"}`
	offerReq := httptest.NewRequest(http.MethodPost, "/signal/offer", bytes.NewBufferString(offer))
	offerRec := httptest.NewRecorder()
	mux.ServeHTTP(offerRec, offerReq)
	if offerRec.Code != http.StatusAccepted {
		t.Fatalf("offer status = %d, want %d", offerRec.Code, http.StatusAccepted)
	}

	answer := `{"sessionId":"sess-poll","sdp":"answer-sdp","fromPeer":"bob"}`
	answerReq := httptest.NewRequest(http.MethodPost, "/signal/answer", bytes.NewBufferString(answer))
	answerRec := httptest.NewRecorder()
	mux.ServeHTTP(answerRec, answerReq)
	if answerRec.Code != http.StatusAccepted {
		t.Fatalf("answer status = %d, want %d", answerRec.Code, http.StatusAccepted)
	}

	pollReq := httptest.NewRequest(http.MethodGet, "/signal/session/sess-poll", nil)
	pollRec := httptest.NewRecorder()
	mux.ServeHTTP(pollRec, pollReq)
	if pollRec.Code != http.StatusOK {
		t.Fatalf("poll status = %d, want %d", pollRec.Code, http.StatusOK)
	}

	var resp sessionStateResponse
	if err := json.NewDecoder(pollRec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode poll response: %v", err)
	}
	if !resp.HasOffer || !resp.HasAnswer {
		t.Fatalf("expected offer+answer true, got %+v", resp)
	}
	if resp.OfferPeer != "alice" || resp.AnswerPeer != "bob" {
		t.Fatalf("unexpected peers: %+v", resp)
	}
}

func TestSessionPollingErrors(t *testing.T) {
	resetSignals()
	mux := newMux()

	t.Run("method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/signal/session/s1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
		}
		if allow := rec.Header().Get("Allow"); allow != "GET, DELETE" {
			t.Fatalf("allow = %q, want %q", allow, "GET, DELETE")
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/signal/session/none", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestSessionDeleteLifecycle(t *testing.T) {
	resetSignals()
	mux := newMux()

	offer := `{"sessionId":"sess-del","sdp":"offer-sdp","fromPeer":"alice"}`
	offerReq := httptest.NewRequest(http.MethodPost, "/signal/offer", bytes.NewBufferString(offer))
	offerRec := httptest.NewRecorder()
	mux.ServeHTTP(offerRec, offerReq)
	if offerRec.Code != http.StatusAccepted {
		t.Fatalf("offer status = %d, want %d", offerRec.Code, http.StatusAccepted)
	}

	delReq := httptest.NewRequest(http.MethodDelete, "/signal/session/sess-del", nil)
	delRec := httptest.NewRecorder()
	mux.ServeHTTP(delRec, delReq)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", delRec.Code, http.StatusNoContent)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/signal/session/sess-del", nil)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusNotFound {
		t.Fatalf("post-delete get status = %d, want %d", getRec.Code, http.StatusNotFound)
	}

	delAgainReq := httptest.NewRequest(http.MethodDelete, "/signal/session/sess-del", nil)
	delAgainRec := httptest.NewRecorder()
	mux.ServeHTTP(delAgainRec, delAgainReq)
	if delAgainRec.Code != http.StatusNotFound {
		t.Fatalf("second delete status = %d, want %d", delAgainRec.Code, http.StatusNotFound)
	}
}

func TestSessionSdpEndpoint(t *testing.T) {
	resetSignals()
	mux := newMux()

	offer := `{"sessionId":"sess-sdp","sdp":"offer-sdp","fromPeer":"alice"}`
	offerReq := httptest.NewRequest(http.MethodPost, "/signal/offer", bytes.NewBufferString(offer))
	offerRec := httptest.NewRecorder()
	mux.ServeHTTP(offerRec, offerReq)
	if offerRec.Code != http.StatusAccepted {
		t.Fatalf("offer status = %d, want %d", offerRec.Code, http.StatusAccepted)
	}

	answer := `{"sessionId":"sess-sdp","sdp":"answer-sdp","fromPeer":"bob"}`
	answerReq := httptest.NewRequest(http.MethodPost, "/signal/answer", bytes.NewBufferString(answer))
	answerRec := httptest.NewRecorder()
	mux.ServeHTTP(answerRec, answerReq)
	if answerRec.Code != http.StatusAccepted {
		t.Fatalf("answer status = %d, want %d", answerRec.Code, http.StatusAccepted)
	}

	req := httptest.NewRequest(http.MethodGet, "/signal/sdp/sess-sdp", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("sdp status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp sessionSdpResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode sdp response: %v", err)
	}
	if resp.SessionID != "sess-sdp" {
		t.Fatalf("sessionId = %q, want %q", resp.SessionID, "sess-sdp")
	}
	if resp.OfferSdp != "offer-sdp" || resp.AnswerSdp != "answer-sdp" {
		t.Fatalf("unexpected sdp payloads: %+v", resp)
	}
}

func TestSessionSdpEndpointErrors(t *testing.T) {
	resetSignals()
	mux := newMux()

	t.Run("method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/signal/sdp/s1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
		}
		if allow := rec.Header().Get("Allow"); allow != http.MethodGet {
			t.Fatalf("allow = %q, want %q", allow, http.MethodGet)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/signal/sdp/none", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("session exists but no sdp", func(t *testing.T) {
		signalsMu.Lock()
		signals["empty-sdp"] = signalRecord{}
		signalsMu.Unlock()

		req := httptest.NewRequest(http.MethodGet, "/signal/sdp/empty-sdp", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}
