package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type prepareRequest struct {
	Operation string `json:"operation"`
	Payload   any    `json:"payload"`
}

type prepareResponse struct {
	RequestID string `json:"requestId"`
	Nonce     string `json:"nonce"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type submitRequest struct {
	RequestID string `json:"requestId"`
	Signature string `json:"signature"`
}

type submitResponse struct {
	TxHash  string `json:"txHash"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type txStatusResponse struct {
	TxHash    string `json:"txHash"`
	RequestID string `json:"requestId"`
	Operation string `json:"operation"`
	Status    string `json:"status"`
}

type preparedRequest struct {
	Operation string
	Nonce     string
	CreatedAt time.Time
}

var (
	requestsMu sync.Mutex
	requests   = map[string]preparedRequest{}
	txMu       sync.Mutex
	txStatus   = map[string]txStatusResponse{}
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/v1/prepare", prepareHandler)
	mux.HandleFunc("/v1/submit", submitHandler)
	mux.HandleFunc("/v1/tx/", txStatusHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("backend-relayer listening on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "backend-relayer",
	})
}

func prepareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req prepareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if strings.TrimSpace(req.Operation) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "operation is required"})
		return
	}

	requestID, err := randomHex(16)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate request id"})
		return
	}
	nonce, err := randomHex(16)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate nonce"})
		return
	}

	requestsMu.Lock()
	requests[requestID] = preparedRequest{
		Operation: req.Operation,
		Nonce:     nonce,
		CreatedAt: time.Now().UTC(),
	}
	requestsMu.Unlock()

	writeJSON(w, http.StatusOK, prepareResponse{
		RequestID: requestID,
		Nonce:     nonce,
		Status:    "prepared",
		Message:   "request prepared",
	})
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if strings.TrimSpace(req.RequestID) == "" || strings.TrimSpace(req.Signature) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "requestId and signature are required"})
		return
	}
	if !strings.HasPrefix(req.Signature, "0x") || len(req.Signature) < 10 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "signature must be a hex string starting with 0x"})
		return
	}

	requestsMu.Lock()
	prepared, ok := requests[req.RequestID]
	if ok {
		delete(requests, req.RequestID)
	}
	requestsMu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "request not found or already submitted"})
		return
	}

	if time.Since(prepared.CreatedAt) > 10*time.Minute {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "prepared request expired"})
		return
	}

	digest := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", req.RequestID, prepared.Nonce, req.Signature)))
	txHash := "0x" + hex.EncodeToString(digest[:])

	writeJSON(w, http.StatusOK, submitResponse{
		TxHash:  txHash,
		Status:  "submitted",
		Message: "accepted for relay",
	})

	txMu.Lock()
	txStatus[txHash] = txStatusResponse{
		TxHash:    txHash,
		RequestID: req.RequestID,
		Operation: prepared.Operation,
		Status:    "submitted",
	}
	txMu.Unlock()
}

func txStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	txHash := strings.TrimPrefix(r.URL.Path, "/v1/tx/")
	if txHash == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tx hash is required"})
		return
	}
	txMu.Lock()
	status, ok := txStatus[txHash]
	txMu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "tx not found"})
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func randomHex(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
