#!/usr/bin/env bash
set -euo pipefail

# Template deployment script for RPC-backed relayer rollout.
# Fill in values or export env vars before use.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELAYER_DIR="$ROOT_DIR/backend-relayer"

: "${RELAYER_ADDR:=0.0.0.0:18080}"
: "${RELAYER_RPC_MODE:=ethereum}"
: "${ETH_RPC_URL:=https://example-rpc.invalid}"
: "${RELAYER_CHAIN_ID:=1}"
: "${RELAYER_SIGNER_PRIVATE_KEY:=0xREPLACE_ME}"
: "${CONTACT_CONTRACT_ADDRESS:=0x0000000000000000000000000000000000000000}"

echo "[1/3] build relayer"
(
  cd "$RELAYER_DIR"
  env GOCACHE=/tmp/go-build-cache go build ./...
)

echo "[2/3] starting relayer in RPC mode"
(
  cd "$RELAYER_DIR"
  env \
    GOCACHE=/tmp/go-build-cache \
    RELAYER_ADDR="$RELAYER_ADDR" \
    RELAYER_RPC_MODE="$RELAYER_RPC_MODE" \
    ETH_RPC_URL="$ETH_RPC_URL" \
    RELAYER_CHAIN_ID="$RELAYER_CHAIN_ID" \
    RELAYER_SIGNER_PRIVATE_KEY="$RELAYER_SIGNER_PRIVATE_KEY" \
    CONTACT_CONTRACT_ADDRESS="$CONTACT_CONTRACT_ADDRESS" \
    go run .
)

echo "[3/3] relayer exited"
