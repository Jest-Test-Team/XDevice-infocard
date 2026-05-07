# Implementation Dashboard

Last updated: 2026-05-07
Coordinator: Main agent
Mode: Parallel multi-agent implementation

## Objective
Implement executable scaffolds for all 4 plans under `implementations/` with non-conflicting ownership.

## Worker Ownership
- Worker A: `implementations/01-ultrasonic-data-transfer`
- Worker B: `implementations/02-nearby-connections`
- Worker C: `implementations/03-animated-qr-visual-handshake`
- Worker D: `implementations/04-web3-sbt-contacts`

## Status
- Worker A: Completed (Rust core-dsp scaffold + tests)
- Worker B: Completed (Flutter/Dart + Go signaling scaffold)
- Worker C: Completed (Rust scaffold + tests)
- Worker D: Completed (Solidity + Go + infra scaffold)

## Completion Criteria
- Each folder contains a coherent project skeleton.
- Core runtime stubs exist (Rust/Go/Solidity/Dart where applicable).
- Each implementation has a local README with setup instructions.

## Verification Sweep (2026-05-07)
- `implementations/01-ultrasonic-data-transfer/core-dsp`
  - `cargo check`: PASS
  - `cargo test`: PASS (1 passed, 0 failed)
- `implementations/03-animated-qr-visual-handshake/fountain-core`
  - `cargo check`: PASS
  - `cargo test`: PASS (1 passed, 0 failed)
- `implementations/02-nearby-connections/nearby-card-drop/signaling-server`
  - `go build ./...`: PASS
- `implementations/04-web3-sbt-contacts/backend-relayer`
  - `go build ./...`: PASS

## Continuation Batch (2026-05-07)
- Plan 01 (`core-dsp`): upgraded frame format with protocol version + checksum validation, improved demodulator rounding behavior.
- Plan 02 (`nearby-card-drop`): implemented guarded Dart state transitions in `nearby_state_machine.dart` (removed invalid switch fallthrough pattern).
- Plan 03 (`fountain-core`): added `transfer_id` to symbols and stricter decode validation (consistent transfer/total/sequence checks).
- Plan 04 (`backend-relayer`): replaced placeholder request flow with in-memory prepared-request tracking, nonce generation, input validation, expiry checks, and deterministic tx hash generation.

### Verification
- `implementations/01-ultrasonic-data-transfer/core-dsp`: `cargo test` PASS
- `implementations/03-animated-qr-visual-handshake/fountain-core`: `cargo test` PASS
- `implementations/04-web3-sbt-contacts/backend-relayer`: `go build ./...` PASS

## Completion Reality Check
- Fully production-complete for all 4 plans: **Not yet**.
- Current status: scaffolds + MVP core logic increments are complete; platform-specific native integrations, full protocol/security hardening, contract tests, and end-to-end cross-device QA are still pending.

## Continuation Batch 2 (2026-05-07)
- Plan 01: added framing negative tests (checksum/version/length) and full-byte modulate-demodulate roundtrip test.
- Plan 02: added signaling server endpoint method guards + tests; expanded Dart state machine with connected/reset transitions and clearer context reset semantics.
- Plan 03: hardened decoder duplicate handling and added coverage for incomplete/mixed/out-of-order/duplicate cases.
- Plan 04: added backend-relayer API tests (happy path + error matrix) and contract test strategy doc under `contracts/test/README.md`.

### Verification
- `implementations/01-ultrasonic-data-transfer/core-dsp`: `cargo test` PASS (5 passed)
- `implementations/03-animated-qr-visual-handshake/fountain-core`: `cargo test` PASS (6 passed)
- `implementations/02-nearby-connections/nearby-card-drop/signaling-server`: `go test ./...` PASS
- `implementations/04-web3-sbt-contacts/backend-relayer`: `go test ./...` PASS

## Remaining To Reach Full Completion
- Plan 01: real DSP modulation/demodulation (FSK/FFT), UniFFI generation/integration validation on iOS/Android.
- Plan 02: actual Flutter app wiring and native Nearby SDK bridge implementation (Android/iOS code, not just contracts/docs).
- Plan 03: true fountain/LT coding implementation and QR frame generator/scanner runtime integration.
- Plan 04: executable Solidity tests, deployment scripts, real chain RPC relaying and signature verification.

## Continuation Batch 3 (2026-05-07)
- Plan 01 (`core-dsp`): added binary payload APIs (`encode_payload_bytes` / `decode_payload_bytes`) and roundtrip binary test coverage.
- Plan 03 (`fountain-core`): introduced parity symbol model to recover one missing data chunk; decoder upgraded for parity-aware reconstruction.
- Plan 04 (`backend-relayer`): added transaction status tracking map and `GET /v1/tx/{hash}` endpoint with tests.

### Verification
- `implementations/01-ultrasonic-data-transfer/core-dsp`: `cargo test` PASS (6 passed)
- `implementations/03-animated-qr-visual-handshake/fountain-core`: `cargo test` PASS (7 passed)
- `implementations/04-web3-sbt-contacts/backend-relayer`: `go test ./...` PASS

### Notes
- Go test in sandbox requires writable cache path, executed with `GOCACHE=/tmp/go-build-cache`.

## Continuation Batch 4 (2026-05-07)
- Plan 01 (`core-dsp`): moved from byte-to-sample mapping to bit-symbol modulation/demodulation using configurable samples-per-bit derived from `sample_rate_hz/symbol_rate_bps`.
- Plan 03 (`fountain-core`): added explicit `payload_len` metadata in symbols and decoder now truncates reconstructed payload by declared length (instead of trimming trailing zero bytes), preserving legitimate binary endings.
- Plan 04 (`backend-relayer`): transaction status now records operation and exposes it via `GET /v1/tx/{hash}`.

### Verification
- `implementations/01-ultrasonic-data-transfer/core-dsp`: `cargo test` PASS (6 passed)
- `implementations/03-animated-qr-visual-handshake/fountain-core`: `cargo test` PASS (7 passed)
- `implementations/04-web3-sbt-contacts/backend-relayer`: `go test ./...` PASS

## Continuation Batch 5 (2026-05-07)
- Plan 04 contracts: added executable Foundry test files and config:
  - `contracts/foundry.toml`
  - `contracts/test/SBTProfile.t.sol`
  - `contracts/test/ConnectionGraph.t.sol`
- Plan 04 relayer runtime: made bind address configurable via `RELAYER_ADDR` (default `127.0.0.1:18080`).
- Plan 04 ops: added local E2E script `scripts/relayer_e2e.sh` that builds, starts, runs prepare/submit/status flow, and validates output.

### Verification
- `implementations/04-web3-sbt-contacts/backend-relayer`: `go test ./...` PASS
- `implementations/04-web3-sbt-contacts/scripts/relayer_e2e.sh`: cannot complete in this sandbox due to listener restriction (`bind: operation not permitted`).
- `forge test`: not executed in this environment (`forge` not installed).

## Continuation Batch 6 (2026-05-07)
- Repository hygiene: expanded `.gitignore` to exclude generated build artifacts for Rust/Go outputs across implementations.
- Plan 04 relayer: added replay/idempotency test ensuring `requestId` is single-use (`second submit` fails with not found/already submitted).

### Verification
- `implementations/04-web3-sbt-contacts/backend-relayer`: `go test ./...` PASS

## Continuation Batch 7 (2026-05-07)
- Plan 02 signaling server: upgraded from bare endpoint stubs to in-memory session signaling flow.
  - `/signal/offer` now validates JSON and stores `sessionId` offer SDP.
  - `/signal/answer` now validates JSON and requires existing offer for the same `sessionId`.
- Plan 02 tests: added end-to-end offer/answer happy path and answer-without-offer negative test.
- Plan 01 protocol tests: added decode invalid-frame and non-UTF8 payload rejection tests.

### Verification
- `implementations/02-nearby-connections/nearby-card-drop/signaling-server`: `go test ./...` PASS
- `implementations/01-ultrasonic-data-transfer/core-dsp`: `cargo test` PASS (8 passed)
