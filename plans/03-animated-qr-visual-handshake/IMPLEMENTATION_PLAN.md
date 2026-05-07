# Plan 03: Animated QR + Fountain Code (Visual Handshake)

## 1. Objective
Implement a camera-based animated QR pipeline that reliably transfers medium-size business-card payloads under frame loss and variable lighting, using Rust fountain code core for cross-platform deterministic recovery.

## 2. Scope
- In scope:
  - Fountain-code encoder/decoder core in Rust.
  - QR frame generation and high-FPS playback.
  - iOS and Android scanner integration with live decode.
  - Transfer completion signaling and integrity checks.
- Out of scope (Phase 1):
  - Arbitrary large media transfer (>250KB).
  - Multi-hop relay between devices.
  - Desktop webcam clients.

## 3. Proposed Architecture
- `fountain-core` (Rust):
  - Payload chunking and LT code symbol generation.
  - Decoder with progressive rank recovery.
  - APIs for ingesting out-of-order chunks.
- iOS scanner (Swift + AVCapture + Vision/OpenCV):
  - Frame sampling, QR extraction, chunk forwarding.
- Android scanner (Kotlin + CameraX + ML Kit):
  - Similar extraction pipeline with adaptive frame skip.
- Sender player module:
  - Animated QR stream generation at tunable FPS.

## 4. Repository Blueprint
```text
visual-handshake/
├── fountain-core/
│   ├── src/
│   │   ├── encoder.rs
│   │   ├── decoder.rs
│   │   ├── symbol.rs
│   │   └── lib.rs
│   ├── tests/
│   └── Cargo.toml
├── clients/
│   ├── ios-scanner/
│   ├── android-scanner/
│   └── shared-ui-spec/
├── tools/
│   ├── qr-frame-generator/
│   └── replay-benchmark/
├── docs/
│   ├── frame-loss-model.md
│   ├── optical-test-conditions.md
│   └── benchmark-results/
└── README.md
```

## 5. Execution Phases
1. Fountain core MVP (Week 1-2)
- Implement LT symbol generation and solve matrix-based decode.
- Add deterministic test vectors and corruption/failure tests.

2. QR payload framing and player (Week 2-3)
- Define compact chunk metadata schema (session/chunk/hash).
- Build frame generator with configurable FPS and redundancy.

3. Mobile scanner integration (Week 3-4)
- Implement frame capture and adaptive sampling strategy.
- Parse QR payload, deduplicate chunk ids, stream to decoder.

4. Completion and UX loop (Week 4)
- Add progress indicator by estimated decode completeness.
- Emit transfer success event with payload hash verification.

5. Stress testing and optimization (Week 5)
- Evaluate different lighting, glare, and motion scenarios.
- Tune redundancy factor and frame pacing per device class.

## 6. Data/Codec Decisions
- Chunk metadata fields:
  - protocol version
  - transfer id
  - source payload hash
  - symbol index
  - symbol seed
  - symbol bytes
- Redundancy baseline: 1.5x symbol budget over source block count.
- Integrity: SHA-256 payload hash + per-symbol CRC.

## 7. Performance Targets
- Effective throughput: 50KB-80KB/s at ~30 FPS under good lighting.
- Frame-loss tolerance: recover with up to 30% missing frames.
- First-decode latency: <1.0s to meaningful progress feedback.

## 8. QA Matrix
- Distance: 20cm / 40cm / 70cm.
- Lighting: low light / office / bright daylight.
- Motion: static / mild shake / aggressive shake.
- Devices: at least 3 iOS + 3 Android models.

## 9. Security and Privacy
- Encrypt payload before fountain encoding.
- Rotate transfer id and nonce per session.
- Never store raw frame images unless debug mode is explicitly enabled.

## 10. Delivery Artifacts
- Rust fountain core with coverage tests.
- iOS and Android scanner/playback demos.
- Benchmark harness and published result set.
- Parameter tuning guide for production rollout.
