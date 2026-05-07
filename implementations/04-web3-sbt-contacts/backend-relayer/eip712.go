package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

var (
	eip712DomainTypeHash     = keccak256Bytes([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
	contactOperationTypeHash = keccak256Bytes([]byte("ContactOperation(string requestId,string operation,string nonce)"))
)

func eip712DigestHex(data signableTypedData) (string, error) {
	digest, err := eip712Digest(data)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(digest), nil
}

func eip712Digest(data signableTypedData) ([]byte, error) {
	domainSeparator, err := eip712DomainSeparator(data.Domain)
	if err != nil {
		return nil, err
	}
	messageHash := eip712ContactOperationHash(data.Message)
	prefix := []byte{0x19, 0x01}
	return keccak256Bytes(prefix, domainSeparator, messageHash), nil
}

func eip712DomainSeparator(domain signableTypedDataDomain) ([]byte, error) {
	chainID := new(big.Int)
	if _, ok := chainID.SetString(strings.TrimSpace(domain.ChainID), 10); !ok {
		return nil, fmt.Errorf("invalid domain.chainId")
	}
	if len(strings.TrimSpace(domain.VerifyingContract)) == 0 {
		return nil, fmt.Errorf("invalid domain.verifyingContract")
	}
	addrWord, err := hexAddressToWord(domain.VerifyingContract)
	if err != nil {
		return nil, fmt.Errorf("invalid domain.verifyingContract")
	}
	return keccak256Bytes(
		eip712DomainTypeHash,
		keccak256Bytes([]byte(domain.Name)),
		keccak256Bytes([]byte(domain.Version)),
		leftPad32(chainID.Bytes()),
		addrWord,
	), nil
}

func eip712ContactOperationHash(msg signableMessage) []byte {
	return keccak256Bytes(
		contactOperationTypeHash,
		keccak256Bytes([]byte(msg.RequestID)),
		keccak256Bytes([]byte(msg.Operation)),
		keccak256Bytes([]byte(msg.Nonce)),
	)
}

func keccak256Bytes(parts ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, p := range parts {
		h.Write(p)
	}
	return h.Sum(nil)
}

func hexAddressToWord(addr string) ([]byte, error) {
	raw := strings.TrimPrefix(strings.TrimSpace(addr), "0x")
	if len(raw) != 40 {
		return nil, fmt.Errorf("address must be 20 bytes")
	}
	decoded, err := hex.DecodeString(raw)
	if err != nil {
		return nil, err
	}
	return leftPad32(decoded), nil
}

func leftPad32(b []byte) []byte {
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}
