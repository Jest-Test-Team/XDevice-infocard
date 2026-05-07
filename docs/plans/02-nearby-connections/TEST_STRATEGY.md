# Test Strategy

## Objective
Validate Plan 02 delivery quality with traceable coverage from state machine rules to user outcomes.

## Test Scope
- Flutter domain/state machine logic
- Android/iOS native bridge event correctness
- End-to-end transfer behavior
- Failure handling and retry semantics
- Optional signaling fallback behavior

## Test Pyramid
- Unit tests: 60%
- Integration tests: 25%
- End-to-end/device tests: 15%

## Environments
- Local simulator/emulator with mock adapter
- Physical device lab:
  - Android: latest-2 major versions
  - iOS: latest-2 major versions
- Staging backend (if signaling enabled)

## Test Matrix (Minimum)
- Android -> Android (3 device pairs)
- iOS -> iOS (3 device pairs)
- Android -> iOS (4 cross-platform pairs)
- Foreground/Background transitions during transfer
- Permissions states:
  - all granted
  - bluetooth denied
  - location/nearby denied

## Coverage Requirements
- State transitions: 100% transition-path coverage from `STATE_MACHINE_SPEC.md`.
- Retry branches: 100% for connection and chunk retry logic.
- Data contract serialization/deserialization: 100% for required fields.
- Native event decoding: 0 unhandled event types.

## Test Suites

### Unit
- `session_state_machine_test.dart`
  - Valid and invalid transitions
  - Guard condition enforcement
  - Timeout transitions
- `transfer_protocol_test.dart`
  - Chunk split/reassembly
  - CRC32 and SHA-256 verification
  - Retry budget exhaustion
- `contracts_test.dart`
  - DTO schema compatibility and version checks

Pass gate
- >= 90% line coverage in `app/lib/state` and `app/lib/core`.

### Integration
- Flutter <-> native bridge contract tests
  - Event envelope parsing
  - Event ordering constraints
  - Endpoint lifecycle mapping
- Mock-native replay tests
  - Replays of recorded callback sequences
  - Fault injections for dropped callbacks

Pass gate
- 500-event replay run with 0 decode errors and 0 state deadlocks.

### End-to-End
- Happy path card transfer:
  - 1 MB payload, 100 runs, success >= 98%
- Stress path:
  - 5 MB payload, 50 runs, success >= 95%
- Recovery path:
  - forced disconnect mid-transfer, recovery succeeds within retry policy >= 90%

Pass gate
- P95 discovery-to-connect <= 8s
- P95 total transfer time (1 MB) <= 12s

## Non-Functional Tests
- Battery impact: 20-minute repeated discovery/connect cycles.
  - Gate: battery drain <= 8% on reference devices.
- Memory pressure test: repeated 5 MB transfers for 30 cycles.
  - Gate: no OOM, no unbounded memory growth.
- Background behavior:
  - App pause/resume during `TRANSFERRING` and `VERIFYING`.

## Failure Injection Plan
- Inject wrong auth code.
- Drop every Nth chunk ack.
- Corrupt a single chunk CRC.
- Simulate adapter unavailable event.
- Simulate permissions revoked mid-session.

Expected result
- Session terminates into correct `FAILED` code with cleanup complete.

## CI/CD Gating
- PR checks required:
  - unit tests
  - integration tests
  - static analysis/lint
- Nightly checks:
  - extended replay/fault-injection suite
- Release-candidate checks:
  - full physical-device matrix

## Defect Policy
- `P0`: data corruption/security issue/crash loop; blocks release.
- `P1`: transfer failure > gate threshold; blocks release.
- `P2`: non-critical UX inconsistency; can ship with mitigation.

## Exit Criteria
- All quality gates green.
- No open P0/P1 defects.
- Observability metrics validate target SLOs over 7-day beta.
