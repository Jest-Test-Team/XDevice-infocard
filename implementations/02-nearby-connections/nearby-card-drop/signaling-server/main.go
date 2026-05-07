package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

type offerRequest struct {
	SessionID string `json:"sessionId"`
	Sdp       string `json:"sdp"`
	FromPeer  string `json:"fromPeer"`
}

type answerRequest struct {
	SessionID string `json:"sessionId"`
	Sdp       string `json:"sdp"`
	FromPeer  string `json:"fromPeer"`
}

type signalRecord struct {
	OfferSdp   string
	AnswerSdp  string
	OfferPeer  string
	AnswerPeer string
}

var (
	signalsMu sync.Mutex
	signals   = map[string]signalRecord{}
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := healthResponse{
		Status:    "ok",
		Service:   "nearby-signaling-server",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode health response: %v", err)
	}
}

func postOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func offerHandler(w http.ResponseWriter, r *http.Request) {
	var req offerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.SessionID) == "" || strings.TrimSpace(req.Sdp) == "" {
		http.Error(w, "sessionId and sdp are required", http.StatusBadRequest)
		return
	}

	signalsMu.Lock()
	record := signals[req.SessionID]
	record.OfferSdp = req.Sdp
	record.OfferPeer = req.FromPeer
	signals[req.SessionID] = record
	signalsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":    "offer-recorded",
		"sessionId": req.SessionID,
	})
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	var req answerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.SessionID) == "" || strings.TrimSpace(req.Sdp) == "" {
		http.Error(w, "sessionId and sdp are required", http.StatusBadRequest)
		return
	}

	signalsMu.Lock()
	record, ok := signals[req.SessionID]
	if !ok || strings.TrimSpace(record.OfferSdp) == "" {
		signalsMu.Unlock()
		http.Error(w, "offer not found for session", http.StatusNotFound)
		return
	}
	record.AnswerSdp = req.Sdp
	record.AnswerPeer = req.FromPeer
	signals[req.SessionID] = record
	signalsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":    "answer-recorded",
		"sessionId": req.SessionID,
	})
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/signal/offer", postOnly(offerHandler))
	mux.HandleFunc("/signal/answer", postOnly(answerHandler))
	return mux
}

func main() {
	addr := ":8080"
	log.Printf("signaling server listening on %s", addr)
	if err := http.ListenAndServe(addr, newMux()); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
