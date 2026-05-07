#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CRATE_DIR="$ROOT_DIR"
UDL_FILE="$CRATE_DIR/uniffi/uniffi.udl"
OUT_DIR="${OUT_DIR:-$CRATE_DIR/uniffi/generated}"
LANG="${1:-all}"

if ! command -v uniffi-bindgen >/dev/null 2>&1; then
  echo "uniffi-bindgen is required. Install with: cargo install uniffi_bindgen_cli" >&2
  exit 1
fi

mkdir -p "$OUT_DIR"

case "$LANG" in
  all)
    uniffi-bindgen generate "$UDL_FILE" --language swift --out-dir "$OUT_DIR/swift"
    uniffi-bindgen generate "$UDL_FILE" --language kotlin --out-dir "$OUT_DIR/kotlin"
    ;;
  swift|kotlin)
    uniffi-bindgen generate "$UDL_FILE" --language "$LANG" --out-dir "$OUT_DIR/$LANG"
    ;;
  *)
    echo "Usage: $0 [all|swift|kotlin]" >&2
    exit 2
    ;;
esac

echo "UniFFI bindings generated at: $OUT_DIR"
