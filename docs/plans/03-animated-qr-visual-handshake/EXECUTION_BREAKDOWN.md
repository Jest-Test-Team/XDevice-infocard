# Plan 03 Execution Breakdown

## Goal
Deliver a production-ready Phase 1 visual handshake path for medium business-card payloads, from sender animation to receiver reconstruction, with deterministic Rust codec behavior and measurable mobile decode reliability.

## Workstreams
1. Codec Core (Rust)
- Build LT-style fountain encoder/decoder with deterministic PRNG per `(transfer_id, symbol_seed)`.
- Implement symbol deduplication, out-of-order ingestion, and completion detection.
- Export stable API surface for mobile bindings.

2. Frame Schema + QR Packaging
- Finalize binary chunk schema and base45/base64url transport envelope.
- Add frame-level integrity (CRC-16) and payload-level SHA-256 validation.
- Define QR size/ECC defaults and payload budget per frame.

3. Sender Animation Pipeline
- Implement frame scheduler (target 24/30 FPS profiles).
- Add redundancy profile controls (`1.25x`, `1.5x`, `1.75x`).
- Emit completion marker frames and restart behavior for late scanners.

4. Mobile Scanner Pipelines
- iOS: AVCapture + Vision/OpenCV decode loop with adaptive frame skip.
- Android: CameraX + ML Kit decode loop with backpressure controls.
- Feed parsed symbols into decoder with de-dup and session filtering.

5. Verification + Benchmarking
- Deterministic vectors, simulated frame-loss tests, and optical scenario tests.
- Collect throughput, completion rate, and latency metrics by device class.
- Publish benchmark summary and recommended default profiles.

## Milestones and Exit Criteria
1. M1: Codec MVP Complete
- Exit criteria:
  - `encode -> shuffle/drop/duplicate -> decode` succeeds at 30% random frame loss for payloads 4KB, 16KB, 64KB.
  - Decoder rejects corrupted symbol CRC and never reports false success.
  - Deterministic vectors stable across 3 consecutive CI runs.

2. M2: Wire Format Frozen
- Exit criteria:
  - `CODEC_SPEC.md` fields implemented exactly and version-gated.
  - Max encoded frame payload size documented and enforced.
  - Backward-incompatible changes blocked unless `protocol_version` increments.

3. M3: Sender + Scanner Loop Operational
- Exit criteria:
  - End-to-end transfer succeeds on at least 1 iOS + 1 Android pair for 16KB payload.
  - Progress feedback appears within 1.0s from first decodable frame.
  - Duplicate-frame rate does not cause decoder stall.

4. M4: Performance Gate
- Exit criteria:
  - Completion success >= 95% under office light + mild shake at 40cm for 16KB payload.
  - Effective throughput >= 50KB/s under good lighting at 30 FPS.
  - P95 completion time for 32KB payload <= 2.8s on mid-tier 2023+ phones.

5. M5: Release Readiness
- Exit criteria:
  - Test matrix in `TEST_STRATEGY.md` executed with logged results.
  - All high-severity risks in `RISK_REGISTER.md` are mitigated or explicitly accepted.
  - Rollout config includes per-device profile defaults and fallback behavior.

## Sequenced Implementation Plan
1. Week 1: Codec foundations
- Implement symbol generation, degree selection, and decode graph/reduction.
- Add symbol/header structs and binary serialization tests.
- Add deterministic seeds and test vector fixtures.

2. Week 2: Codec hardening + schema lock
- Add CRC validation and payload hash completion checks.
- Freeze frame schema and parser behavior (strict/lenient fields).
- Add failure-mode tests (truncated frame, wrong transfer id, stale session).

3. Week 3: Sender frame player
- Implement frame generation with tunable FPS and redundancy.
- Add frame pacing telemetry (`generated_fps`, `display_fps`, drop count).
- Implement completion marker cadence and replay loop.

4. Week 4: Mobile ingestion
- Integrate iOS and Android scanner loops.
- Add decode queue, duplicate suppression cache, and timeout reset logic.
- Surface progress UI from decoder completion ratio estimates.

5. Week 5: Stress/optimization
- Run optical and motion matrix.
- Tune QR version/ECC/FPS per device tier.
- Lock default operational profile and document rollout checklist.

## Quality/Performance Targets
- Decoder correctness: 0 false-positive completion in 10,000 randomized trials.
- Frame-loss resilience: decode success at <= 30% random loss and <= 15% burst loss (burst length up to 8 frames).
- UX latency: first meaningful progress in < 1.0s; final integrity verification in < 150ms after last required symbol.
- Resource budgets:
  - iOS CPU: <= 35% average during scan on A15-class device.
  - Android CPU: <= 40% average during scan on Snapdragon 7-class device.
  - Memory: decoder working set <= 32MB for 64KB payload profile.

## Validation Steps (Release Gate)
1. Run deterministic codec vectors and randomized property tests.
2. Run replay-benchmark tool with loss profiles `0%`, `10%`, `20%`, `30%`, burst `8`.
3. Run device optical matrix: distances `20/40/70cm`, light `low/office/daylight`, motion `static/mild/aggressive`.
4. Verify completion hash checks and mismatch handling UX.
5. Publish benchmark table with pass/fail against targets.

## Ownership and Handoff
- Worker 3 scope: plan docs/spec/testing/risk artifacts only.
- Implementation owners consume these docs as normative references for code and QA.
- Any protocol changes require coordinated update to `CODEC_SPEC.md`, vectors, and test gates in same PR.
