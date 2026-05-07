# Risk Register - Plan 01 Ultrasonic Data Transfer

## Rating Model
- Likelihood: Low / Medium / High
- Impact: Low / Medium / High
- Priority heuristic: Highest attention to `High-High` and `Medium-High`.

## Top Risks

| ID | Risk | Likelihood | Impact | Mitigation | Trigger / Early Warning |
|---|---|---|---|---|---|
| R1 | Device speaker/mic frequency response varies, causing decode failures in 18-20kHz band | High | High | Build per-device calibration profiles, adaptive thresholding, and fallback tone spacing profile | Success rate drops below 95% on any target device in quiet environment |
| R2 | Background noise and AGC behavior distort symbols | High | High | Add robust preamble detection, dynamic noise floor estimation, and retransmit policy tuning | Rising CRC mismatch and sync-loss counts in office/fan tests |
| R3 | iOS/Android audio pipeline differences break deterministic behavior | Medium | High | Enforce shared Rust logic for framing/protocol, add byte-parity tests for both bindings, standardize PCM normalization | Parity test mismatch between mobile decode and Rust direct decode |
| R4 | Latency exceeds 3.5s median for 4KB payload | Medium | High | Optimize symbol duration, chunk size, buffering, and retransmit backoff; profile CPU hot paths in demodulation | Median transfer time trend >3.5s in nightly benchmarks |
| R5 | Session retry policy causes stalls or premature failures | Medium | Medium | Validate timeout/retry constants using packet-loss simulation; add bounded exponential backoff | Sessions terminate with retry exhaustion under moderate noise before reaching max logical retries |
| R6 | FFI boundary bugs (memory/ownership/marshaling) create crashes | Medium | High | Keep FFI API minimal, add stress tests for repeated ingest/poll/reset, run sanitizers where available | Crash signatures concentrated in JNI/Swift bridge functions |
| R7 | Replay or tampering concerns if nonce/checks are inconsistent | Low | High | Enforce session nonce checks in protocol layer and verify aggregate checksum on complete | Duplicate/old session IDs accepted by receiver in integration tests |
| R8 | Logging or telemetry leaks sensitive contact payload data | Low | High | Redaction-by-default logging policy, lint checks for forbidden payload log patterns | Any raw payload bytes found in test logs or crash reports |
| R9 | Audio interruptions (calls/route changes) leave session unrecoverable | Medium | Medium | Implement pause/resume from last ACKed frame and deterministic state reset APIs | Transfer fails to recover in interruption test scenarios |
| R10 | Schedule slip due to device lab availability constraints | Medium | Medium | Lock device matrix early, reserve QA slots weekly, parallelize fixture-based testing while waiting for devices | Missed weekly device test window or blocked access to required OS versions |

## Active Monitoring Plan
- Weekly risk review during phase checkpoint.
- Daily automated metrics ingestion:
  - Decode success rate by device and environment.
  - CRC fail rate and retransmit count distribution.
  - Transfer latency median/p95 trends.
- Escalation rule:
  - Any `High impact` risk with active trigger for 2 consecutive days becomes release blocker.

## Contingency Actions
1. If R1/R2 persist at end of Phase 4: ship with constrained supported-device list and documented compatibility table.
2. If R4 persists near release: reduce max payload from 8KB to 4KB for Phase 1 GA and queue optimization follow-up.
3. If R6 or R8 occurs in pre-release: block release until crash/privacy issue is verified fixed and regression-tested.
