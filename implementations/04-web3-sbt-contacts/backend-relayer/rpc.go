package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type RelayTxRequest struct {
	RequestID     string            `json:"requestId"`
	Operation     string            `json:"operation"`
	Signature     string            `json:"signature"`
	TypedData     signableTypedData `json:"typedData"`
	TypedDataHash string            `json:"typedDataHash"`
}

type RelayResult struct {
	TxHash string `json:"txHash"`
	Status string `json:"status"`
}

type ChainRPC interface {
	RelayContactOperation(ctx context.Context, req RelayTxRequest) (RelayResult, error)
}

type InMemoryChainRPC struct{}

func (m *InMemoryChainRPC) RelayContactOperation(_ context.Context, req RelayTxRequest) (RelayResult, error) {
	blob, err := json.Marshal(req)
	if err != nil {
		return RelayResult{}, fmt.Errorf("marshal relay payload: %w", err)
	}
	digest := sha256.Sum256(blob)
	return RelayResult{TxHash: "0x" + hex.EncodeToString(digest[:]), Status: "submitted"}, nil
}

func buildChainRPC() ChainRPC {
	// Template hook: swap this implementation with a JSON-RPC-backed client.
	return &InMemoryChainRPC{}
}
