#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELAYER_DIR="$ROOT_DIR/backend-relayer"
RELAYER_ADDR="${RELAYER_ADDR:-127.0.0.1:18080}"
BASE_URL="http://$RELAYER_ADDR"

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

echo "[1/5] building relayer"
(
  cd "$RELAYER_DIR"
  env GOCACHE=/tmp/go-build-cache go build ./...
)

echo "[2/5] starting relayer"
(
  cd "$RELAYER_DIR"
  env GOCACHE=/tmp/go-build-cache RELAYER_ADDR="$RELAYER_ADDR" go run . >/tmp/web3-relayer.log 2>&1
) &
SERVER_PID=$!

for _ in {1..30}; do
  if ! kill -0 "$SERVER_PID" 2>/dev/null; then
    echo "relayer exited unexpectedly. log:" >&2
    sed -n '1,120p' /tmp/web3-relayer.log >&2 || true
    exit 1
  fi
  if curl -fsS "$BASE_URL/health" >/dev/null 2>&1; then
    break
  fi
  sleep 0.2
done

HEALTH_JSON="$(curl -fsS "$BASE_URL/health")"
echo "[3/5] health: $HEALTH_JSON"

PREPARE_JSON='{"operation":"connect","payload":{"from":"did:example:a","to":"did:example:b"}}'
PREPARE_RESP="$(curl -fsS -X POST "$BASE_URL/v1/prepare" -H 'Content-Type: application/json' -d "$PREPARE_JSON")"
REQUEST_ID="$(printf '%s' "$PREPARE_RESP" | sed -n 's/.*"requestId":"\([^"]*\)".*/\1/p')"

if [[ -z "$REQUEST_ID" ]]; then
  echo "prepare response missing requestId: $PREPARE_RESP" >&2
  exit 1
fi

echo "[4/5] prepare: $PREPARE_RESP"

SUBMIT_JSON="{\"requestId\":\"$REQUEST_ID\",\"signature\":\"0x1234abcd90\"}"
SUBMIT_RESP="$(curl -fsS -X POST "$BASE_URL/v1/submit" -H 'Content-Type: application/json' -d "$SUBMIT_JSON")"
TX_HASH="$(printf '%s' "$SUBMIT_RESP" | sed -n 's/.*"txHash":"\([^"]*\)".*/\1/p')"

if [[ -z "$TX_HASH" ]]; then
  echo "submit response missing txHash: $SUBMIT_RESP" >&2
  exit 1
fi

echo "[5/5] submit: $SUBMIT_RESP"
STATUS_RESP="$(curl -fsS "$BASE_URL/v1/tx/$TX_HASH")"
echo "status: $STATUS_RESP"

echo "E2E check passed"
