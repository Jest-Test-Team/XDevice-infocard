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
- demodulation supports selectable detection method:
  - correlator (original)
  - Goertzel (DFT-window style single-bin detector)
- adaptive noise-floor thresholding is now included to stabilize ambiguous symbol decisions

This keeps the existing Rust API signatures unchanged while replacing sign-level modulation/detection.

## Known Limits

- No preamble/sync word yet; decode assumes perfect symbol alignment.
- No forward error correction; corruption is observable as raw bit errors.
- No clock drift correction.
- Frequencies are fixed constants and not yet configurable by API.
- Carrier sensing, AGC, and channel estimation are not implemented.

## Plan01 Status

### Completed in repo

1. FSK-like modulation/demodulation path in `core-dsp`.
2. Selectable demod method (`Correlator` and `Goertzel`) by config.
3. Basic adaptive threshold/noise-floor logic for ambiguous symbol handling.
4. UniFFI helper script: `core-dsp/scripts/generate-uniffi.sh`.
5. Cross-platform integration checklist: `docs/uniffi-integration-checklist.md`.

### Still external / pending

1. Preamble + sync detection for frame start recovery.
2. CRC/FEC and BER-oriented robustness metrics.
3. Configurable tone profiles and symbol timing guard options.
4. Full iOS/Android binary packaging and app-level audio pipeline wiring.
5. Device-lab validation for latency/BER across environments.

## Next implementation milestones

1. Add preamble plus sync detection to support frame start recovery.
2. Add CRC/FEC path and error-rate metrics in tests.
3. Add configurable tone profiles and symbol timing guard options.
4. Add waveform playback/capture harnesses in iOS and Android clients.
