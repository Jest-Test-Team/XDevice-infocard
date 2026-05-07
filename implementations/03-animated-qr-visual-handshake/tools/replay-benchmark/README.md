# Replay Benchmark

Companion benchmark utility that replays encoded symbols with deterministic frame drops and evaluates decode success via `fountain-core` APIs.

## Usage

```bash
cargo run --manifest-path tools/replay-benchmark/Cargo.toml -- "visual-handshake" 6 25
```

Arguments:
1. `payload` (default: `animated-qr-visual-handshake`)
2. `frame_size` in bytes (default: `12`)
3. `drop_percent` 0-99 (default: `20`)

The command prints encoded/received symbol counts and `decode_success` / `decoded_matches`.
