package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func resetRequests() {
	requestsMu.Lock()
	defer requestsMu.Unlock()
	requests = map[string]preparedRequest{}
	txMu.Lock()
	defer txMu.Unlock()
	txStatus = map[string]txStatusResponse{}
}

func TestPrepareAndSubmitHappyPath(t *testing.T) {
	resetRequests()

	prepareBody := `{"operation":"connect","payload":{"from":"did:example:a","to":"did:example:b"}}`
	prepareReq := httptest.NewRequest(http.MethodPost, "/v1/prepare", strings.NewReader(prepareBody))
	prepareRec := httptest.NewRecorder()

	prepareHandler(prepareRec, prepareReq)

	if prepareRec.Code != http.StatusOK {
		t.Fatalf("prepare status = %d, want %d", prepareRec.Code, http.StatusOK)
	}

	var prepareResp prepareResponse
	if err := json.NewDecoder(prepareRec.Body).Decode(&prepareResp); err != nil {
		t.Fatalf("decode prepare response: %v", err)
	}

	if prepareResp.Status != "prepared" {
		t.Fatalf("prepare status field = %q, want prepared", prepareResp.Status)
	}
	if prepareResp.RequestID == "" || prepareResp.Nonce == "" {
		t.Fatalf("prepare response missing requestId/nonce: %+v", prepareResp)
	}

	submitPayload := submitRequest{
		RequestID: prepareResp.RequestID,
		Signature: "0x1234abcd90",
	}
	submitJSON, err := json.Marshal(submitPayload)
	if err != nil {
		t.Fatalf("marshal submit payload: %v", err)
	}

	submitReq := httptest.NewRequest(http.MethodPost, "/v1/submit", bytes.NewReader(submitJSON))
	submitRec := httptest.NewRecorder()
	submitHandler(submitRec, submitReq)

	if submitRec.Code != http.StatusOK {
		t.Fatalf("submit status = %d, want %d", submitRec.Code, http.StatusOK)
	}

	var submitResp submitResponse
	if err := json.NewDecoder(submitRec.Body).Decode(&submitResp); err != nil {
		t.Fatalf("decode submit response: %v", err)
	}
	if submitResp.Status != "submitted" {
		t.Fatalf("submit status field = %q, want submitted", submitResp.Status)
	}
	if !strings.HasPrefix(submitResp.TxHash, "0x") || len(submitResp.TxHash) != 66 {
		t.Fatalf("invalid tx hash format: %q", submitResp.TxHash)
	}

	txReq := httptest.NewRequest(http.MethodGet, "/v1/tx/"+submitResp.TxHash, nil)
	txRec := httptest.NewRecorder()
	txStatusHandler(txRec, txReq)
	if txRec.Code != http.StatusOK {
		t.Fatalf("tx status = %d, want %d", txRec.Code, http.StatusOK)
	}
	var txResp txStatusResponse
	if err := json.NewDecoder(txRec.Body).Decode(&txResp); err != nil {
		t.Fatalf("decode tx response: %v", err)
	}
	if txResp.TxHash != submitResp.TxHash || txResp.RequestID != prepareResp.RequestID {
		t.Fatalf("tx status response mismatch: %+v", txResp)
	}
	if txResp.Operation != "connect" {
		t.Fatalf("tx operation = %q, want connect", txResp.Operation)
	}
}

func TestPrepareHandlerErrors(t *testing.T) {
	resetRequests()

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantErr:    "method not allowed",
		},
		{
			name:       "invalid json",
			method:     http.MethodPost,
			body:       "{",
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid json",
		},
		{
			name:       "missing operation",
			method:     http.MethodPost,
			body:       `{"payload":{"x":1}}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "operation is required",
		},
		{
			name:       "unsupported operation",
			method:     http.MethodPost,
			body:       `{"operation":"delete","payload":{"x":1}}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "unsupported operation",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/v1/prepare", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			prepareHandler(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			var errResp map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if errResp["error"] != tc.wantErr {
				t.Fatalf("error = %q, want %q", errResp["error"], tc.wantErr)
			}
		})
	}
}

func TestSubmitHandlerErrors(t *testing.T) {
	resetRequests()

	t.Run("method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/submit", nil)
		rec := httptest.NewRecorder()
		submitHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusMethodNotAllowed, "method not allowed")
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/submit", strings.NewReader("{"))
		rec := httptest.NewRecorder()
		submitHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusBadRequest, "invalid json")
	})

	t.Run("missing requestId or signature", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/submit", strings.NewReader(`{"requestId":"","signature":""}`))
		rec := httptest.NewRecorder()
		submitHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusBadRequest, "requestId and signature are required")
	})

	t.Run("invalid signature format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/submit", strings.NewReader(`{"requestId":"abc","signature":"1234"}`))
		rec := httptest.NewRecorder()
		submitHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusBadRequest, "signature must be a hex string starting with 0x")
	})

	t.Run("invalid signature hex", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/submit", strings.NewReader(`{"requestId":"abc","signature":"0xzzzzzzzzzz"}`))
		rec := httptest.NewRecorder()
		submitHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusBadRequest, "signature is not valid hex")
	})

	t.Run("request not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/submit", strings.NewReader(`{"requestId":"missing","signature":"0x1234abcd90"}`))
		rec := httptest.NewRecorder()
		submitHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusNotFound, "request not found or already submitted")
	})

	t.Run("expired request", func(t *testing.T) {
		resetRequests()
		requestsMu.Lock()
		requests["expired-id"] = preparedRequest{
			Operation: "connect",
			Nonce:     "00112233",
			CreatedAt: time.Now().UTC().Add(-11 * time.Minute),
		}
		requestsMu.Unlock()

		req := httptest.NewRequest(http.MethodPost, "/v1/submit", strings.NewReader(`{"requestId":"expired-id","signature":"0x1234abcd90"}`))
		rec := httptest.NewRecorder()
		submitHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusBadRequest, "prepared request expired")
	})
}

func TestTxStatusErrors(t *testing.T) {
	resetRequests()
	t.Run("method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/tx/0xabc", nil)
		rec := httptest.NewRecorder()
		txStatusHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusMethodNotAllowed, "method not allowed")
	})
	t.Run("missing tx hash", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/tx/", nil)
		rec := httptest.NewRecorder()
		txStatusHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusBadRequest, "tx hash is required")
	})
	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/tx/0xdeadbeef", nil)
		rec := httptest.NewRecorder()
		txStatusHandler(rec, req)
		assertErrorResponse(t, rec, http.StatusNotFound, "tx not found")
	})
}

func TestSubmitRequestIsSingleUse(t *testing.T) {
	resetRequests()

	prepareReq := httptest.NewRequest(
		http.MethodPost,
		"/v1/prepare",
		strings.NewReader(`{"operation":"connect","payload":{"a":1}}`),
	)
	prepareRec := httptest.NewRecorder()
	prepareHandler(prepareRec, prepareReq)
	if prepareRec.Code != http.StatusOK {
		t.Fatalf("prepare status = %d, want %d", prepareRec.Code, http.StatusOK)
	}

	var prepareResp prepareResponse
	if err := json.NewDecoder(prepareRec.Body).Decode(&prepareResp); err != nil {
		t.Fatalf("decode prepare response: %v", err)
	}

	submitBody := `{"requestId":"` + prepareResp.RequestID + `","signature":"0x1234abcd90"}`

	firstReq := httptest.NewRequest(http.MethodPost, "/v1/submit", strings.NewReader(submitBody))
	firstRec := httptest.NewRecorder()
	submitHandler(firstRec, firstReq)
	if firstRec.Code != http.StatusOK {
		t.Fatalf("first submit status = %d, want %d", firstRec.Code, http.StatusOK)
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/v1/submit", strings.NewReader(submitBody))
	secondRec := httptest.NewRecorder()
	submitHandler(secondRec, secondReq)
	assertErrorResponse(t, secondRec, http.StatusNotFound, "request not found or already submitted")
}

func assertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantErr string) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d", rec.Code, wantStatus)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["error"] != wantErr {
		t.Fatalf("error = %q, want %q", resp["error"], wantErr)
	}
}
