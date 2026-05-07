# Plan 02: Google Nearby Connections (Nearby Card Drop)

## 1. Objective
Deliver a production-ready cross-platform card-exchange app using Google Nearby Connections with robust connection-state management, predictable UX, and optional Go signaling fallback for difficult network topologies.

## 2. Scope
- In scope:
  - Flutter app with unified UX for discovery, pairing, and transfer.
  - Android/iOS native bridge for Nearby APIs.
  - State machine, retries, and transfer confirmations.
  - Optional lightweight Go signaling service.
- Out of scope (Phase 1):
  - Persistent social graph.
  - Desktop platform support.
  - Offline analytics ingestion.

## 3. Proposed Architecture
- `app/lib` (Flutter):
  - Domain state machine and session orchestration.
  - UI for discovery list, 4-digit code verification, and transfer status.
- `app/android` + `app/ios`:
  - Native wrappers around Nearby discovery/advertising/connection callbacks.
  - Event bridge into Dart stream layer.
- `signaling-server` (Go, optional):
  - WebSocket channel for constrained environments.
  - Session hints and fallback peer hint distribution.

## 4. Repository Blueprint
```text
nearby-card-drop/
├── app/
│   ├── lib/
│   │   ├── core/
│   │   ├── features/discovery/
│   │   ├── features/transfer/
│   │   └── state/
│   ├── android/
│   └── ios/
├── signaling-server/
│   ├── cmd/server/
│   ├── internal/ws/
│   ├── internal/session/
│   └── go.mod
├── docs/
│   ├── state-machine.md
│   ├── failure-modes.md
│   └── rollout-checklist.md
└── README.md
```

## 5. Execution Phases
1. Flutter app skeleton and contract definitions (Week 1)
- Define Dart interfaces for discovery/advertising/send/receive.
- Implement app state machine and mock adapter for local testing.

2. Native bridge implementation (Week 2-3)
- Android Nearby SDK integration with callback translation.
- iOS counterpart integration and permission handling.
- Method channel streams with clear event schemas.

3. Transfer protocol and resilience (Week 3-4)
- Payload chunking, checksum validation, and reassembly.
- Timeout policy and reconnect strategy.
- UI progression from discover -> authenticate -> transfer -> done.

4. Optional signaling service (Week 4)
- Build Go WebSocket service for fallback session assist.
- Add feature toggle to enable/disable signaling mode.

5. QA and release readiness (Week 5)
- Cross-device matrix tests.
- Battery and background behavior checks.
- Crash-free and success-rate targets.

## 6. State Machine
- `IDLE`: no active discovery.
- `DISCOVERING`: scanning nearby advertisers.
- `CONNECTING`: endpoint negotiation in progress.
- `AUTHENTICATING`: user confirms short auth code.
- `TRANSFERRING`: chunk stream in progress.
- `VERIFYING`: checksum and payload integrity validation.
- `COMPLETED` / `FAILED`: terminal states with retry action.

## 7. Data Contract
- `HandshakePayload`: app version, protocol version, endpoint id, nonce.
- `TransferChunk`: session id, chunk index, total chunks, bytes, crc.
- `TransferReceipt`: session id, status, error code, digest.

## 8. Observability
- Structured logs by session id.
- Metrics:
  - discovery-to-connect latency
  - transfer completion rate
  - retry count per completed transfer
- Crash diagnostics with event breadcrumbs.

## 9. Security and Compliance
- Short auth code verification required before transfer.
- Application-level payload encryption (X25519 + AES-GCM).
- Data minimization for logs/analytics.

## 10. Delivery Artifacts
- Flutter app with native bridges.
- Optional Go signaling service.
- End-to-end test scripts and QA matrix.
- Operational runbook for release.
