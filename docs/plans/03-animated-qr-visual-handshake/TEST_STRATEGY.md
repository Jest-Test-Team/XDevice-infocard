# Test Strategy (Plan 03 Animated QR Visual Handshake)

## Objectives
- Prove deterministic codec correctness.
- Validate end-to-end decode reliability under realistic optical and motion conditions.
- Enforce measurable performance and quality gates before release.

## Test Levels
1. Unit tests (Rust codec)
- Symbol serialization/deserialization.
- CRC generation/validation.
- Deterministic seed -> degree/index mapping.
- Duplicate/invalid symbol handling.

2. Property and fuzz tests
- Random payloads and transfer IDs.
- Random symbol order, duplicates, and loss injection.
- Corrupted header/payload bit flips.

3. Integration tests
- Sender frame generator -> scanner parser -> decoder loop.
- Session reset and retry flows.
- Completion marker handling.

4. Device/system tests
- Real camera-based scan in target distance/light/motion matrix.
- iOS + Android cross-device interop.

5. Performance/soak tests
- Sustained transfer loop for 30 minutes.
- CPU/memory budgets during repeated sessions.

## Core Test Matrix
Dimensions:
- Payload sizes: `4KB`, `16KB`, `32KB`, `64KB`.
- Frame loss: `0%`, `10%`, `20%`, `30%` random; burst loss length `4` and `8`.
- Distance: `20cm`, `40cm`, `70cm`.
- Lighting: `low`, `office`, `bright daylight`.
- Motion: `static`, `mild shake`, `aggressive shake`.
- Devices: minimum `3 iOS + 3 Android` models.

## Pass/Fail Quality Gates
1. Correctness gate
- 100% pass on deterministic golden vectors.
- 0 false-success decode (hash mismatch accepted as success) in 10,000 randomized trials.

2. Reliability gate
- >= 99% completion at <= 20% random loss for 16KB and 32KB payloads (office light, mild shake, 40cm).
- >= 95% completion at 30% random loss for 16KB payload (same condition).

3. Performance gate
- Good lighting throughput >= 50KB/s at 30 FPS profile.
- P95 completion time:
  - 16KB <= 1.8s
  - 32KB <= 2.8s
  - 64KB <= 5.5s
- First progress feedback latency < 1.0s.

4. Resource gate
- iOS scan CPU average <= 35%; Android <= 40%.
- Memory peak during 64KB decode <= 32MB incremental working set.
- No crash/leak in 30-minute soak run.

## Validation Procedure
1. Offline codec validation
- Run unit/property/fuzz suites.
- Produce summary with trial counts and failure categories.

2. Replay benchmark validation
- Use recorded frame streams with controlled synthetic loss.
- Compare decode rank progression versus expected curves.

3. Optical test execution
- Run full distance/light/motion matrix on each device pair.
- Record completion rate, median and P95 completion time, and error counts.

4. Regression validation
- Re-run golden vectors and 16KB/32KB critical scenarios for each protocol or scanner change.

## Required Test Artifacts
- `golden_vectors.json` (fixed transfer sessions and expected symbols/hash).
- Benchmark CSV/JSON containing:
  - `device_model`, `os_version`, `payload_size`, `fps`, `redundancy_factor`, `loss_profile`, `completion`, `completion_time_ms`, `throughput_kbps`, `cpu_avg`, `mem_peak_mb`.
- Failure logs with categorized reason:
  - `crc_error`, `unsupported_version`, `session_mismatch`, `decoder_stall`, `hash_mismatch`.

## CI and Release Gates
- CI required checks:
  - Rust unit/property tests.
  - Determinism vectors.
  - Replay loss simulation suite.
- Release-blocking conditions:
  - Any correctness gate failure.
  - Reliability below target on two consecutive runs.
  - Untriaged high-severity defect in decode or integrity path.

## Sign-off Criteria
Release sign-off requires all:
1. All quality gates pass with stored artifacts.
2. No open P0/P1 issues in risk register without explicit acceptance.
3. Documented default runtime profile (`fps`, `ecc`, `redundancy`) per device tier.
