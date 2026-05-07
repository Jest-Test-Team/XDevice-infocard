# Execution Breakdown - Plan 01 Ultrasonic Data Transfer

## Planning Assumptions
- Team shape: 1 DSP engineer, 1 iOS engineer, 1 Android engineer, 1 QA engineer (shared).
- Work calendar: 5 working days per week.
- Estimates below are elapsed working days with 20% buffer already included.

## Phase 0 - Setup and Baseline Alignment
- Estimation: 2 days
- Dependencies:
  - Final approval of `IMPLEMENTATION_PLAN.md` scope.
  - Access to 2 iOS and 2 Android physical devices.
- Entry criteria:
  - Plan 01 accepted by engineering + QA leads.
  - Target devices and OS versions listed.
- Tasks:
  - Create `core-dsp` crate skeleton and module stubs (`framing`, `protocol`, `modulator`, `demodulator`).
  - Set fixed coding conventions (sample rate constants, naming, logging redaction policy).
  - Add CI jobs for Rust unit tests and mobile binding compile checks.
  - Prepare initial benchmark fixture folder (`docs/benchmarks/raw-fixtures`).
- Exit criteria:
  - Repository structure created and builds on CI.
  - First “hello-frame” encode/decode smoke test passes in Rust.

## Phase 1 - Frame Protocol and DSP Encode/Decode MVP
- Estimation: 6 days
- Dependencies:
  - Phase 0 complete.
- Entry criteria:
  - Frame field definitions approved.
  - Baseline sample rates selected (44.1kHz and 48kHz).
- Tasks:
  - Implement frame serializer/deserializer with strict validation and CRC32 verification.
  - Implement FSK modulation pipeline with deterministic symbol timing.
  - Implement FFT-based demodulation and symbol classifier for clean-channel input.
  - Create golden fixtures: payload sizes 1KB, 4KB, 8KB across both sample rates.
  - Add property-based tests for frame parser and sequence ordering.
- Exit criteria:
  - Rust-only round-trip pass rate >= 99.5% using synthetic/noiseless fixtures.
  - Decoder rejects malformed headers, bad CRC, invalid sequence values.

## Phase 2 - Session Protocol and Reliability Controls
- Estimation: 4 days
- Dependencies:
  - Phase 1 complete.
- Entry criteria:
  - MVP encode/decode stable in deterministic test fixtures.
- Tasks:
  - Implement session state machine: `DISCOVER -> HANDSHAKE -> TX -> VERIFY -> ACK -> COMPLETE`.
  - Add packet retransmit policy (max retry count, ack timeout backoff).
  - Add session nonce and replay protection checks.
  - Add payload encryption boundary interface (AES-GCM assumed upstream).
  - Implement state reset and duplicate packet suppression.
- Exit criteria:
  - Integration tests confirm successful transfer with injected packet loss up to 10%.
  - Session timeout and retransmit behavior match protocol rules.

## Phase 3 - iOS/Android FFI and Audio Pipeline Integration
- Estimation: 7 days
- Dependencies:
  - Phase 1 and 2 complete.
  - UniFFI generation pipeline available.
- Entry criteria:
  - Stable Rust APIs frozen for first mobile integration (`encode`, `ingest`, `poll_frame`, `reset_session`).
- Tasks:
  - Generate and integrate UniFFI bindings for Swift and Kotlin.
  - Implement iOS AVFoundation capture/playback graph with Rust bridge.
  - Implement Android AudioRecord/AudioTrack thread pipeline with JNI bridge.
  - Add conversion guards (PCM int16 <-> f32 normalization).
  - Implement mobile permission and interruption handling (calls, route changes).
- Exit criteria:
  - Both platforms can send/receive 1KB payload successfully on local loopback and device-to-device tests.
  - Byte-level parity checks pass between Rust direct decode and mobile pipeline decode.

## Phase 4 - Device Tuning and Noise Robustness
- Estimation: 5 days
- Dependencies:
  - Phase 3 complete.
- Entry criteria:
  - End-to-end transfer works on at least one iOS model and one Android model.
- Tasks:
  - Calibrate gain normalization and detection threshold by device profile.
  - Tune buffer sizes for 44.1kHz and 48kHz modes.
  - Add adaptive threshold mode for moderate ambient noise.
  - Run controlled noise tests (quiet room, office chatter, fan noise).
  - Capture failure traces and classify root causes (sync loss, CRC fail, timeout).
- Exit criteria:
  - 1KB/4KB/8KB transfer success >= 99% at 0.5m in quiet conditions.
  - >= 95% success in moderate noise profile.

## Phase 5 - Performance Validation and Release Gate
- Estimation: 3 days
- Dependencies:
  - Phase 4 complete.
- Entry criteria:
  - QA test matrix complete with no unresolved P0/P1 defects.
- Tasks:
  - Run benchmark suite and collect median/p95 transfer times.
  - Verify target: 4KB transfer <= 3.5s median on mid-tier devices.
  - Finalize production checklist (security redaction, crash-free run, retry behavior).
  - Publish benchmark and pass/fail report in `docs/benchmarks`.
- Exit criteria:
  - All acceptance gates in `TEST_STRATEGY.md` pass.
  - Go/No-Go review signed by engineering + QA.

## Cross-Phase Dependency Graph
1. Phase 0 -> required for all other phases.
2. Phase 1 -> blocks Phase 2 and Phase 3.
3. Phase 2 -> required before robust device testing in Phase 4.
4. Phase 3 -> required before Phase 4.
5. Phase 4 -> required before Phase 5.

## Total Delivery Estimate
- Total elapsed estimate: 27 working days.
- Critical path: Phase 0 -> 1 -> 2 -> 3 -> 4 -> 5.
- Recommended checkpoint cadence: end-of-phase demo and risk review.
