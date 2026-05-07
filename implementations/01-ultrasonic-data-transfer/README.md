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

## Current Modem Path (Plan 01)

- `core-dsp` now uses a lightweight FSK-like bit modem:
- bit `0` maps to a tone near `18.0 kHz`
- bit `1` maps to a tone near `19.5 kHz`
- modulation emits a sine tone per bit-symbol at configured `sample_rate_hz / symbol_rate_bps`
- demodulation compares per-symbol energy at both tones and selects the stronger bin

This keeps the existing Rust API signatures unchanged while replacing sign-level modulation/detection.

## Known Limits

- No preamble/sync word yet; decode assumes perfect symbol alignment.
- No forward error correction; corruption is observable as raw bit errors.
- No clock drift correction or adaptive thresholding.
- Frequencies are fixed constants and not yet configurable by API.
- Carrier sensing, AGC, and channel estimation are not implemented.

## Next implementation milestones

1. Add preamble plus sync detection to support frame start recovery.
2. Add CRC/FEC path and error-rate metrics in tests.
3. Add configurable tone profiles and symbol timing guard options.
4. Add UniFFI generation and platform binding pipelines.
5. Add waveform playback/capture harnesses in iOS and Android clients.
