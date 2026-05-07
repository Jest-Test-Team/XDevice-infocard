# Replay Benchmark

Companion benchmark utility that replays encoded symbols with deterministic frame drops and evaluates decode recovery rates via `fountain-core` APIs.

## Usage

```bash
cargo run --manifest-path tools/replay-benchmark/Cargo.toml -- "visual-handshake" 6 "0,10,20,30,40,50" 20
```

Arguments:
1. `payload` (default: `animated-qr-visual-handshake`)
2. `frame_size` in bytes (default: `12`)
3. `drop_percents_csv` 0-99 list (default: `0,10,20,30,40,50`)
4. `runs_per_drop` (default: `20`)

Output includes a CSV summary per drop percentage:

`drop_percent,recovery_rate,decoded_match_rate,avg_received_symbols`
