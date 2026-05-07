package main

import (
	"context"
	"testing"
)

func TestInMemoryChainRPCDeterministic(t *testing.T) {
	rpc := &InMemoryChainRPC{}
	req := RelayTxRequest{
		RequestID:     "rid-1",
		Operation:     "connect",
		Signature:     "0x1234abcd90",
		TypedData:     buildSignableTypedData("rid-1", "connect", "nonce-1"),
		TypedDataHash: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	r1, err := rpc.RelayContactOperation(context.Background(), req)
	if err != nil {
		t.Fatalf("relay: %v", err)
	}
	r2, err := rpc.RelayContactOperation(context.Background(), req)
	if err != nil {
		t.Fatalf("relay: %v", err)
	}
	if r1.TxHash != r2.TxHash {
		t.Fatalf("tx hash mismatch: %q vs %q", r1.TxHash, r2.TxHash)
	}
	if r1.Status != "submitted" {
		t.Fatalf("status = %q, want submitted", r1.Status)
	}
}
