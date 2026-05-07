package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type prepareRequest struct {
	Operation string `json:"operation"`
	Payload   any    `json:"payload"`
}

type prepareResponse struct {
	RequestID string            `json:"requestId"`
	Nonce     string            `json:"nonce"`
	Status    string            `json:"status"`
	Message   string            `json:"message"`
	Signable  signableTypedData `json:"signable"`
}

type submitRequest struct {
	RequestID string            `json:"requestId"`
	Signature string            `json:"signature"`
	Signable  signableTypedData `json:"signable"`
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
	Signable  signableTypedData
	CreatedAt time.Time
}

type signableTypedData struct {
	PrimaryType string                      `json:"primaryType"`
	Domain      signableTypedDataDomain     `json:"domain"`
	Types       map[string][]signableMember `json:"types"`
	Message     signableMessage             `json:"message"`
}

type signableTypedDataDomain struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainID           string `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

type signableMember struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type signableMessage struct {
	RequestID string `json:"requestId"`
	Operation string `json:"operation"`
	Nonce     string `json:"nonce"`
}

var (
	requestsMu sync.Mutex
	requests   = map[string]preparedRequest{}
	txMu       sync.Mutex
	txStatus   = map[string]txStatusResponse{}
)

func main() {
	addr := os.Getenv("RELAYER_ADDR")
	if strings.TrimSpace(addr) == "" {
		addr = "127.0.0.1:18080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/v1/prepare", prepareHandler)
	mux.HandleFunc("/v1/submit", submitHandler)
	mux.HandleFunc("/v1/tx/", txStatusHandler)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend-relayer listening on %s", addr)
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
	if !isAllowedOperation(req.Operation) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported operation"})
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
	signable := buildSignableTypedData(requestID, req.Operation, nonce)
	requests[requestID] = preparedRequest{
		Operation: req.Operation,
		Nonce:     nonce,
		Signable:  signable,
		CreatedAt: time.Now().UTC(),
	}
	requestsMu.Unlock()

	writeJSON(w, http.StatusOK, prepareResponse{
		RequestID: requestID,
		Nonce:     nonce,
		Status:    "prepared",
		Message:   "request prepared",
		Signable:  signable,
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
	if _, err := hex.DecodeString(strings.TrimPrefix(req.Signature, "0x")); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "signature is not valid hex"})
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
	if err := validateSignable(req.Signable, prepared.Signable, req.RequestID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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

func isAllowedOperation(op string) bool {
	switch strings.ToLower(strings.TrimSpace(op)) {
	case "connect", "exchange":
		return true
	default:
		return false
	}
}

func buildSignableTypedData(requestID, operation, nonce string) signableTypedData {
	return signableTypedData{
		PrimaryType: "ContactOperation",
		Domain: signableTypedDataDomain{
			Name:              "Web3SBTContactsRelayer",
			Version:           "1",
			ChainID:           "0",
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Types: map[string][]signableMember{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"ContactOperation": {
				{Name: "requestId", Type: "string"},
				{Name: "operation", Type: "string"},
				{Name: "nonce", Type: "string"},
			},
		},
		Message: signableMessage{
			RequestID: requestID,
			Operation: operation,
			Nonce:     nonce,
		},
	}
}

func validateSignable(actual, expected signableTypedData, requestID string) error {
	if actual.PrimaryType == "" || actual.Domain.Name == "" || len(actual.Types) == 0 || actual.Message.RequestID == "" {
		return fmt.Errorf("signable payload is required")
	}
	if actual.PrimaryType != expected.PrimaryType {
		return fmt.Errorf("signable primaryType mismatch")
	}
	if actual.Domain != expected.Domain {
		return fmt.Errorf("signable domain mismatch")
	}
	if !sameTypes(actual.Types, expected.Types) {
		return fmt.Errorf("signable types mismatch")
	}
	if actual.Message != expected.Message {
		return fmt.Errorf("signable message mismatch")
	}
	if actual.Message.RequestID != requestID {
		return fmt.Errorf("signable requestId mismatch")
	}
	return nil
}

func sameTypes(actual, expected map[string][]signableMember) bool {
	if len(actual) != len(expected) {
		return false
	}
	for key, expectedMembers := range expected {
		actualMembers, ok := actual[key]
		if !ok || len(actualMembers) != len(expectedMembers) {
			return false
		}
		for i := range expectedMembers {
			if actualMembers[i] != expectedMembers[i] {
				return false
			}
		}
	}
	return true
}
