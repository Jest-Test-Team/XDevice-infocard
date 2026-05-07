package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			req := httptest.NewRequest(http.MethodPost, endpoint, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusAccepted {
				t.Fatalf("expected status %d, got %d", http.StatusAccepted, rec.Code)
			}
		})
	}
}
