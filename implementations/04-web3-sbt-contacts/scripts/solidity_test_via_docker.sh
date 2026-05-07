#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONTRACTS_DIR="$ROOT_DIR/contracts"

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required for this script" >&2
  exit 1
fi

echo "Running forge build/test inside Foundry container"
docker run --rm \
  -v "$CONTRACTS_DIR":/work \
  -w /work \
  ghcr.io/foundry-rs/foundry:latest \
  sh -lc "forge --version && forge build && forge test -vv"
