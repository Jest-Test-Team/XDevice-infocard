package main

import "testing"

func TestEIP712DigestDeterministic(t *testing.T) {
	signable := buildSignableTypedData("rid-123", "connect", "nonce-xyz")
	signable.Domain.ChainID = "31337"
	signable.Domain.VerifyingContract = "0x1111111111111111111111111111111111111111"

	d1, err := eip712DigestHex(signable)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	d2, err := eip712DigestHex(signable)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	if d1 != d2 {
		t.Fatalf("digest not deterministic: %q vs %q", d1, d2)
	}
	if len(d1) != 66 || d1[:2] != "0x" {
		t.Fatalf("unexpected digest format: %q", d1)
	}
}

func TestEIP712DigestChangesOnMessageMutation(t *testing.T) {
	base := buildSignableTypedData("rid-123", "connect", "nonce-xyz")
	base.Domain.ChainID = "31337"
	base.Domain.VerifyingContract = "0x1111111111111111111111111111111111111111"

	d1, err := eip712DigestHex(base)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	mutated := base
	mutated.Message.Nonce = "nonce-other"
	d2, err := eip712DigestHex(mutated)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	if d1 == d2 {
		t.Fatalf("digest should change when message changes")
	}
}
