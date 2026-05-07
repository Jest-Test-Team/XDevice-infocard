# Plan 01: Ultrasonic Data Transfer (SonicCard Exchange)

## 1. Objective
Build a near-field business-card exchange protocol over ultrasonic audio (18kHz-20kHz) with shared Rust DSP core for iOS and Android, deterministic behavior across devices, and measurable reliability under real-world acoustic noise.

## 2. Scope
- In scope:
  - End-to-end send/receive of vCard payloads up to 8KB.
  - Rust-based modulation/demodulation with FFI bindings to Swift/Kotlin.
  - Mobile-native audio I/O pipelines and permissions.
  - Basic session protocol: discover, handshake, transmit, verify checksum, ack.
  - QA benchmarks for decode success rate and time-to-complete.
- Out of scope (Phase 1):
  - Long-range broadcast (>2m).
  - Multi-peer mesh transfers.
  - Background-mode autonomous transfer.

## 3. Proposed Architecture
- `core-dsp` (Rust):
  - Framing layer (preamble, packet id, sequence, payload, CRC32).
  - FSK/PSK modulation pipeline.
  - FFT-based demodulation and symbol detection.
  - Error detection and optional Reed-Solomon extension point.
- `clients/ios` (Swift + AVFoundation):
  - Capture/playback graph setup.
  - Real-time buffer handoff to Rust.
  - UI state machine and permission flow.
- `clients/android` (Kotlin + AudioRecord/AudioTrack):
  - Audio thread lifecycle and low-latency mode.
  - Buffer conversion and JNI bridge.
  - Foreground transfer UX.

## 4. Repository Blueprint
```text
sonic-card-exchange/
├── core-dsp/
│   ├── src/
│   │   ├── framing.rs
│   │   ├── modulator.rs
│   │   ├── demodulator.rs
│   │   ├── fft.rs
│   │   ├── protocol.rs
│   │   └── lib.rs
│   ├── benches/
│   ├── tests/
│   ├── uniffi.udl
│   ├── build.rs
│   └── Cargo.toml
├── clients/
│   ├── ios/
│   └── android/
├── docs/
│   ├── benchmarks/
│   ├── frequency-profiles/
│   └── qa-test-matrix.md
└── README.md
```

## 5. Execution Phases
1. Protocol and DSP baseline (Week 1-2)
- Define on-wire frame format and checksum rules.
- Implement Rust modulator with synthetic sample output.
- Build demodulator against known-good fixtures.

2. Cross-platform binding integration (Week 2-3)
- Generate UniFFI Swift/Kotlin bindings.
- Integrate Rust artifacts (XCFramework, JNI/AAR).
- Validate byte-for-byte parity in encode/decode tests across platforms.

3. Real device audio integration (Week 3-4)
- Tune buffer size, sample rate (44.1k/48k), gain normalization.
- Add adaptive detection thresholds.
- Handle interrupts (calls, audio route changes).

4. Reliability hardening and QA (Week 4-5)
- Add packet retransmit/ack timeout policy.
- Test in noisy indoor/outdoor environments.
- Publish benchmark report and pass/fail matrix.

## 6. Technical Decisions
- Language: Rust for deterministic DSP and memory safety.
- FFI: UniFFI for lower integration overhead and consistency.
- Symbol scheme: Start with robust low-bitrate FSK, then evaluate PSK as optional optimization.
- Integrity: CRC32 minimum; optional FEC extension for high-loss conditions.

## 7. API Contract (Core)
- `encode(payload: Vec<u8>) -> Vec<f32>`: bytes to PCM samples.
- `ingest(samples: Vec<f32>) -> DecodeState`: streaming decode result.
- `poll_frame() -> Option<Vec<u8>>`: returns completed payload frame.
- `reset_session(session_id)`: reset state for retries.

## 8. QA and Benchmarks
- Functional:
  - 1KB/4KB/8KB payload round-trip decode accuracy >= 99% at 0.5m.
- Performance:
  - End-to-end transfer time for 4KB <= 3.5s median on mid-tier devices.
- Robustness:
  - Maintain >= 95% success with background chatter and moderate fan noise.

## 9. Security and Privacy
- Payload encrypted before audio serialization (AES-GCM with ephemeral key exchange).
- Session nonce to prevent replay.
- No raw contact data written to logs.

## 10. Delivery Artifacts
- Rust core crate with tests + benches.
- iOS and Android demo clients.
- QA report and reproducible benchmark scripts.
- Production-readiness checklist.
