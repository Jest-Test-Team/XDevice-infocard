package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type prepareRequest struct {
    Operation string `json:"operation"`
    Payload   any    `json:"payload"`
}

type prepareResponse struct {
    RequestID string `json:"requestId"`
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

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/v1/prepare", prepareHandler)
    mux.HandleFunc("/v1/submit", submitHandler)

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
        "status": "ok",
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

    _ = req // placeholder for future validation and payload generation

    writeJSON(w, http.StatusOK, prepareResponse{
        RequestID: "placeholder-request-id",
        Status:    "prepared",
        Message:   "prepare placeholder: implement calldata/gas simulation",
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

    _ = req // placeholder for future signature verification and chain submission

    writeJSON(w, http.StatusOK, submitResponse{
        TxHash:  "0xplaceholder",
        Status:  "submitted",
        Message: "submit placeholder: implement signer + RPC broadcast",
    })
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}
