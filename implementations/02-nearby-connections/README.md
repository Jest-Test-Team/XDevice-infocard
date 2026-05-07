# Plan 02: Nearby Connections Implementation Status

## Completion Matrix

| Scope | Status | Notes |
| --- | --- | --- |
| Flutter minimal runnable skeleton | Completed (in-repo scaffold) | Added `pubspec.yaml`, existing `lib/main.dart` state UI, and app scaffold docs. |
| UI screens/state flow docs | Completed | Added `docs/ui-screens-state-flow.md` and retained existing state-machine docs. |
| Android/iOS bridge contracts (machine-readable) | Completed | Added `contracts/nearby_bridge.schema.json` and `contracts/bridge_contract_matrix.md`. |
| Dart bridge adapter stubs | Completed | Added `lib/core/adapters/method_channel_nearby_bridge.dart` to map method/event channels. |
| Signaling server expiry + cleanup | Completed | Added TTL-based session expiry, background cleanup loop, and `/signal/cleanup` endpoint. |
| Signaling server tests | Completed | Added cleanup/expiry tests and method validation tests for cleanup endpoint. |
| `go test` | Passing | `go test ./...` passes for `signaling-server`. |

## Layout

- `nearby-card-drop/app` - scaffold notes and app generation guidance
- `nearby-card-drop/contracts` - bridge schema and contract matrix
- `nearby-card-drop/lib` - app entrypoint, state machine, contracts, and bridge adapter stubs
- `nearby-card-drop/signaling-server` - Go signaling server with expiry/cleanup and tests
- `nearby-card-drop/docs` - state machine, failure modes, and UI/state flow

## Local Verification

1. `cd implementations/02-nearby-connections/nearby-card-drop/signaling-server`
2. `go test ./...`
3. `go run .`
4. `curl http://localhost:8080/health`
5. `curl -X POST http://localhost:8080/signal/cleanup`

## Remaining External Tasks

- Run `flutter create .` in `nearby-card-drop` to generate full platform runner metadata for direct `flutter run`.
- Implement Android native bridge bindings (Google Nearby Connections) against `contracts/nearby_bridge.schema.json`.
- Implement iOS native bridge bindings (MultipeerConnectivity) against `contracts/nearby_bridge.schema.json`.
- Validate real-device cross-platform transfer and permission/error behavior on physical devices.
