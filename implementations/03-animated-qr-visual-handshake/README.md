# Plan 03: Animated QR Visual Handshake

Scaffold for the animated QR visual handshake implementation.

## Structure

- `fountain-core/`: Rust fountain codec crate with deterministic mixed-symbol repair decoding.
- `clients/ios-scanner/`: iOS scanner client notes.
- `clients/android-scanner/`: Android scanner client notes.
- `tools/qr-frame-generator/`: QR frame generation and payload parser notes.
- `tools/replay-benchmark/`: Replay benchmark tool notes.
- `docs/scanner-player-contract.md`: Lightweight scanner/player integration contract.

## Progress

Completed:
- Multi-symbol repair codec in `fountain-core` (deterministic LT-like mixed XOR symbols + peeling decoder).
- Decoder/test coverage for duplicate detection, incomplete sets, mixed metadata rejection, and multi-loss recovery.
- Scanner/player payload contract for frame-level integration.
- Sample payload parser utility for scanner-side ingestion validation.
- Replay benchmark expanded to recovery-rate sweeps across multiple drop percentages.

Uncompleted:
- Runtime scanner clients (iOS/Android) direct decode integration wiring.
- Live animated playback loop integration with capture timing telemetry.
- End-to-end device-to-device runtime validation across mixed camera conditions.
