# sonic-card-exchange (Plan 01 Scaffold)

Starter repository scaffold for ultrasonic card exchange experiments.

## Layout

- `core-dsp/`: Rust DSP/protocol crate (framing, modulation, demodulation, protocol)
- `clients/ios/`: iOS integration notes for future UniFFI-generated bindings
- `clients/android/`: Android integration notes for future UniFFI-generated bindings
- `docs/`: QA and validation templates

## Build

```bash
cd core-dsp
cargo check
cargo test
```

## Run (current state)

No runtime app is included yet. This scaffold focuses on a compilable DSP core crate and integration docs.

## Next implementation milestones

1. Implement ultrasonic symbol mapping and sync detection in `core-dsp`.
2. Add UniFFI generation and platform binding pipelines.
3. Add waveform playback/capture harnesses in iOS and Android clients.
