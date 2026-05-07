# Risk Register (Plan 03 Animated QR Visual Handshake)

## Scoring Model
- Likelihood: `1 (low)` to `5 (high)`
- Impact: `1 (low)` to `5 (high)`
- Risk score: `Likelihood x Impact`
- Severity bands:
  - `P0`: score >= 20
  - `P1`: 12-19
  - `P2`: 6-11
  - `P3`: <= 5

## Active Risks
1. Decoder false completion or data corruption acceptance
- ID: `R-001`
- Score: `2 x 5 = 10 (P2)`
- Trigger: Completion event emitted without strict payload hash verification.
- Impact: Silent contact-card corruption.
- Mitigation:
  - Enforce SHA-256 match as hard completion gate.
  - Add 10,000-trial randomized no-false-success test gate.
- Validation: 0 false success in randomized/fuzz suites.
- Contingency: Block release and revert protocol change introducing regression.

2. High frame-loss environments cause decode stalls
- ID: `R-002`
- Score: `4 x 4 = 16 (P1)`
- Trigger: Burst loss or motion blur reduces unique symbol throughput.
- Impact: Transfer timeout and poor UX.
- Mitigation:
  - Default redundancy 1.5x with profile bump to 1.75x on weak conditions.
  - Completion markers + replay loop every 12 frames.
  - Adaptive scanner frame skip/backpressure tuning.
- Validation: >=95% completion at 30% random loss in target condition.
- Contingency: Auto-switch to low-FPS high-redundancy profile.

3. Cross-platform nondeterminism in symbol generation
- ID: `R-003`
- Score: `3 x 5 = 15 (P1)`
- Trigger: Different PRNG or index sampling behavior between Rust and mobile bindings.
- Impact: Decoder incompatibility across devices.
- Mitigation:
  - Rust codec as canonical implementation.
  - Golden vectors for first 256 symbols across platforms.
  - Ban floating-point math in selection path.
- Validation: Exact vector match in CI for iOS/Android bindings.
- Contingency: Disable noncanonical path and ship Rust FFI-only mode.

4. QR payload density too high for real-world camera decode
- ID: `R-004`
- Score: `4 x 4 = 16 (P1)`
- Trigger: Oversized frame payload causes scanner decode failures in low light.
- Impact: Throughput collapse, high failure rate.
- Mitigation:
  - Cap frame payload size per QR version/ECC profile.
  - Maintain per-device profile table for version + ECC + FPS.
  - Validate with full optical matrix before defaults are locked.
- Validation: Completion/throughput gates met at 20/40/70cm.
- Contingency: Reduce symbol bytes/frame and increase redundancy.

5. CPU/thermal throttling on mid-tier devices
- ID: `R-005`
- Score: `3 x 4 = 12 (P1)`
- Trigger: Sustained scanning at high FPS overheats device, dropping decode quality.
- Impact: Latency spikes, session failures.
- Mitigation:
  - Backpressure-aware frame processing.
  - Tiered FPS profile (24 FPS fallback).
  - 30-minute soak tests with thermal telemetry.
- Validation: CPU/memory gates plus no crash in soak tests.
- Contingency: Auto-throttle FPS when processing lag exceeds threshold.

6. Session collision or stale-frame contamination
- ID: `R-006`
- Score: `2 x 4 = 8 (P2)`
- Trigger: Reused transfer IDs or late frames from previous session.
- Impact: Decode confusion and potential hash mismatch.
- Mitigation:
  - 128-bit random transfer IDs.
  - Strict session filtering and idle timeout reset.
- Validation: Session-mismatch tests and replay contamination tests pass.
- Contingency: Force new session generation and clear decode buffers.

7. Security/privacy leak via debug frame storage
- ID: `R-007`
- Score: `2 x 5 = 10 (P2)`
- Trigger: Raw camera frames stored unintentionally.
- Impact: Exposure of sensitive card data.
- Mitigation:
  - Debug frame capture off by default.
  - Explicit opt-in with auto-expiry and local-only storage.
  - Ensure payload encrypted pre-encoding.
- Validation: Storage audit in QA, no frame persistence in default mode.
- Contingency: Remove debug storage feature from release build.

8. Incomplete observability blocks root-cause analysis
- ID: `R-008`
- Score: `3 x 3 = 9 (P2)`
- Trigger: Missing telemetry for CRC/session/hash/decode-rank events.
- Impact: Slow incident response and tuning.
- Mitigation:
  - Standard metric schema and error counters.
  - Persist benchmark artifacts for all gate runs.
- Validation: Required metrics present in test artifacts.
- Contingency: Halt tuning decisions until telemetry baseline is restored.

## Risk Review Cadence
- Weekly risk review during implementation.
- Immediate review on any P0/P1 incident.
- Update risk scores after each major benchmark cycle.

## Release Risk Policy
- No open P0 risks.
- P1 risks require mitigation evidence and explicit acceptance owner.
- All unresolved risks must include contingency + rollback notes.
