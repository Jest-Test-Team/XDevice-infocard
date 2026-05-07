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
