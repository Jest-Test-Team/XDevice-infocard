# Execution Breakdown

## Goal
Ship Plan 02 (Nearby Card Drop) as a production-ready Flutter app with Android/iOS Nearby bridges, deterministic session lifecycle, and an optional Go signaling fallback.

## Delivery Window
- Target duration: 5 weeks
- Start assumption: Week 1 begins immediately after plan approval.
- Release candidate date gate: end of Week 5.

## Workstreams
- WS1 App Core (Flutter domain + UI)
- WS2 Native Bridge (Android/iOS)
- WS3 Transfer Protocol + Resilience
- WS4 Optional Signaling Service (Go)
- WS5 QA, Observability, and Release Readiness

## Milestones and Exit Gates

### Milestone M1: Contracts + App Skeleton (Week 1)
Deliverables
- Flutter module scaffolding under `app/lib`:
  - `core/nearby_adapter.dart`
  - `state/session_state_machine.dart`
  - `features/discovery/`, `features/transfer/`
- Domain contracts for handshake/chunk/receipt.
- Mock adapter for simulator/local deterministic runs.

Actionable API contracts
- `abstract class NearbyAdapter`:
  - `Stream<NearbyEvent> events();`
  - `Future<void> startDiscovery(DiscoveryConfig config);`
  - `Future<void> startAdvertising(AdvertisingConfig config);`
  - `Future<void> requestConnection(String endpointId, HandshakePayload payload);`
  - `Future<void> acceptConnection(String endpointId);`
  - `Future<void> rejectConnection(String endpointId, {String? reason});`
  - `Future<void> sendChunk(String endpointId, TransferChunk chunk);`
  - `Future<void> disconnect(String endpointId);`
  - `Future<void> stopAll();`
- DTOs (JSON-serializable):
  - `HandshakePayload {appVersion, protocolVersion, endpointId, nonce, pubKey}`
  - `TransferChunk {sessionId, chunkIndex, totalChunks, bytesBase64, crc32}`
  - `TransferReceipt {sessionId, status, errorCode, digestSha256}`

Exit gates (must pass)
- 100% compile/build green on Flutter analyzer and unit tests.
- State machine transitions tested for all happy-path states.
- Mock adapter can run a full transfer simulation in CI in < 20s.

### Milestone M2: Native Bridge Alpha (Week 2-3)
Deliverables
- Android Nearby integration:
  - Discovery, advertising, connection lifecycle callbacks.
- iOS Nearby integration parity (capability-gated if API differences).
- Unified event bridge into Dart stream.

Bridge event schema (native -> Dart)
- `NearbyEvent` envelope:
  - `type`: `DISCOVERY_FOUND | DISCOVERY_LOST | CONNECTION_INITIATED | CONNECTION_RESULT | AUTH_CODE | PAYLOAD_PROGRESS | PAYLOAD_RECEIVED | DISCONNECTED | ERROR`
  - `sessionId`
  - `endpointId`
  - `timestampMs`
  - `payload` (event-specific map)

Exit gates
- Android: 3-device physical test, discovery success >= 95% in 30 trials.
- iOS: 3-device physical test, discovery success >= 90% in 30 trials.
- Event schema conformance tests: 0 decoding failures in 500 replayed events.

### Milestone M3: Protocol Hardening (Week 3-4)
Deliverables
- Chunking/reassembly with checksum validation.
- Retry + timeout rules integrated into state machine.
- UI flow finalized: discover -> authenticate -> transferring -> verifying -> complete/fail.

Protocol constants (initial)
- `MAX_CHUNK_BYTES = 32 * 1024`
- `ACK_TIMEOUT_MS = 2000`
- `CHUNK_RETRY_LIMIT = 3`
- `CONNECTION_RETRY_LIMIT = 2`
- `SESSION_TIMEOUT_MS = 45000`

Exit gates
- Transfer completion >= 98% for 1 MB payload in same-room tests (n=100).
- Transfer completion >= 95% for 5 MB payload (n=50).
- Corrupted chunk detection = 100% in fault injection tests.

### Milestone M4: Optional Signaling Fallback (Week 4)
Deliverables
- Go WebSocket signaling server:
  - Session create/join/hint events.
- Feature-flagged client fallback path.

Server API (v1)
- `GET /healthz` -> `200 {"status":"ok"}`
- `GET /readyz` -> `200` when ws loop and in-memory session store ready
- `WS /v1/signal?session_id=<id>&peer_id=<id>`
- WS messages:
  - `offer_hint`, `answer_hint`, `peer_presence`, `session_close`, `error`

Feature flag
- `nearby.signaling.enabled` (default `false`)

Exit gates
- When enabled, fallback connection success improves by >= 20% in constrained-network scenario tests.
- No regression in direct nearby mode metrics.

### Milestone M5: Release Readiness (Week 5)
Deliverables
- QA matrix completion (Android/iOS versions, foreground/background, permission states).
- Observability dashboards and alerts.
- Operational runbook and rollback steps.

Release quality gates
- Crash-free sessions >= 99.5% over 7-day beta window.
- End-to-end transfer success >= 97% across device matrix.
- P95 discovery-to-connect latency <= 8s.
- P95 total transfer time (1 MB) <= 12s.

## RACI (Lean)
- Mobile Lead: WS1/WS3 owner, release sign-off.
- Android Engineer: WS2 Android owner.
- iOS Engineer: WS2 iOS owner.
- Backend Engineer: WS4 owner.
- QA Lead: WS5 owner, gate enforcement.
- PM: milestone governance and scope control.

## Dependency Checklist
- Google Nearby SDK versions locked and documented.
- iOS/Android runtime permissions copy reviewed.
- Feature-flag service/remote config available.
- Analytics pipeline accepts session-scoped metrics.

## Definition of Done
- All milestone exit gates passed.
- All P1/P0 defects closed.
- State machine spec and test strategy fully traceable to implementation.
- Rollback path validated in staging.
