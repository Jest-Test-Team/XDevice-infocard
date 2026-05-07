# Plan 02: Nearby Connections Scaffold

This directory contains the initial scaffold for `nearby-card-drop`, including:
- app and library structure
- Dart state machine and interface stubs
- Android/iOS bridge contract stubs
- Minimal Go signaling service
- State-machine and failure-mode docs

## Layout

- `nearby-card-drop/app` - app-level entry area placeholder
- `nearby-card-drop/lib/core` - shared core interfaces
- `nearby-card-drop/lib/features/discovery` - peer discovery contracts
- `nearby-card-drop/lib/features/transfer` - transfer contracts
- `nearby-card-drop/lib/state` - state machine skeleton
- `nearby-card-drop/android` - Android bridge contract README
- `nearby-card-drop/ios` - iOS bridge contract README
- `nearby-card-drop/signaling-server` - minimal Go health service
- `nearby-card-drop/docs` - state machine and failure mode docs

## Setup

1. Enter scaffold root:
   - `cd implementations/02-nearby-connections/nearby-card-drop`
2. Review Dart contracts in `lib/` and align with app integration plan.
3. Run signaling server locally:
   - `cd signaling-server`
   - `go run .`
4. Verify health endpoint:
   - `curl http://localhost:8080/health`

## Notes

- Dart files are interface/skeleton stubs only and are not wired to Flutter runtime yet.
- Native bridge readmes define a shared method channel contract for both platforms.
