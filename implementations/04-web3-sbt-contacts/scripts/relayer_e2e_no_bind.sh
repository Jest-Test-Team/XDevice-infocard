#!/usr/bin/env bash
set -euo pipefail

# Sandbox-safe fallback: no local server bind required.
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELAYER_DIR="$ROOT_DIR/backend-relayer"

echo "Running in-process relayer verification via Go tests"
(
  cd "$RELAYER_DIR"
  env GOCACHE=/tmp/go-build-cache go test ./... -run 'TestPrepareAndSubmitHappyPath|TestSubmitRequestIsSingleUse|TestTxStatusErrors' -v
)
