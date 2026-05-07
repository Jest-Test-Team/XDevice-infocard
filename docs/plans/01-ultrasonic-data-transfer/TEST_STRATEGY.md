# Test Strategy - Plan 01 Ultrasonic Data Transfer

## 1. Objectives
- Validate protocol correctness, cross-platform determinism, and real-device reliability.
- Gate release with measurable thresholds tied to Plan 01 targets.

## 2. Test Layers

### Unit Tests (Rust Core)
- Scope:
  - Frame encoding/decoding and CRC validation.
  - Protocol state transitions and retry logic.
  - Symbol mapping/demapping for FSK.
- Required cases:
  - Valid frame decode for each frame type.
  - Invalid preamble/version/length rejection.
  - CRC corruption detection.
  - Sequence reordering and duplicate suppression.
- Gate:
  - `>= 95%` line coverage on `framing.rs` and `protocol.rs`.
  - 0 failing unit tests in CI.

### Integration Tests (Rust + FFI)
- Scope:
  - End-to-end encode -> modulate -> demodulate -> decode in process.
  - Swift/Kotlin bindings parity for fixed fixtures.
- Required cases:
  - 1KB, 4KB, 8KB payload transfer with both sample rates.
  - Mixed chunk sizes and final-frame short chunk.
  - Session reset and restart after timeout.
- Gate:
  - Byte-for-byte payload equality across Rust, iOS, Android decode outputs.
  - Pass rate `>= 99%` in synthetic low-noise pipeline over 1,000 runs.

### Device Tests (Real Hardware)
- Scope:
  - Physical speaker/microphone path on target devices.
  - Permission flow, interruptions, and route changes.
- Matrix minimum:
  - iOS: 2 models (latest major iOS and one previous major).
  - Android: 2 models (mid-tier + flagship).
  - Distances: 0.2m, 0.5m, 1.0m.
  - Environments: quiet room, office chatter, fan noise.
- Required cases:
  - First-time permission denied then granted.
  - Incoming call/audio focus interruption mid-transfer.
  - Retransmit behavior under induced packet loss.
- Gate:
  - 0.5m quiet: success `>= 99%` for 1KB/4KB/8KB.
  - 0.5m moderate noise: success `>= 95%`.
  - No unrecovered crashes/ANRs during 500 transfer cycles.

### Performance Tests
- Scope:
  - Transfer latency and CPU impact.
  - Retransmit overhead under noise.
- Metrics:
  - End-to-end transfer latency (median, p95).
  - Decode CPU time per second of audio.
  - Retransmit count distribution.
- Gate:
  - 4KB payload median transfer time `<= 3.5s` on mid-tier device.
  - p95 transfer time `<= 5.0s` at 0.5m quiet.
  - Sustained CPU utilization under test cap defined by platform team.

## 3. Test Data and Fixtures
- Deterministic payload fixtures:
  - `fixture_1kb_vcard.bin`
  - `fixture_4kb_vcard.bin`
  - `fixture_8kb_vcard.bin`
- Audio fixtures:
  - Clean synthetic encoded streams.
  - Noise-mixed variants (SNR buckets: 30dB, 20dB, 10dB).
- Corruption fixtures:
  - Bit flips in header and payload.
  - Dropped and duplicated frame sequences.

## 4. Automation Plan
- CI per pull request:
  - Rust unit + integration tests.
  - FFI generation and compile check for iOS/Android bindings.
- Nightly:
  - Extended randomized fuzz tests for frame parser.
  - Benchmark replay for trend monitoring.
- Pre-release:
  - Full device matrix run with signed report.

## 5. Acceptance Gates by Milestone
1. Milestone A (Core complete)
- All Rust unit/integration tests green.
- Synthetic pass rate `>= 99.5%` noiseless.

2. Milestone B (Mobile integration complete)
- iOS + Android parity tests green.
- Device loopback on each platform stable for 100 cycles.

3. Milestone C (Release candidate)
- Device matrix passes all reliability gates.
- Performance targets met (`4KB <= 3.5s median`).
- No open P0/P1 defects.

## 6. Exit Criteria for Plan 01
- Every gate above is met and documented.
- Benchmark and failure analysis artifacts are published under `docs/benchmarks`.
- Engineering + QA sign-off recorded.
