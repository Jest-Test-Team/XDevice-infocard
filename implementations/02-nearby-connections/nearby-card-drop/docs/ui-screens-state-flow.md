# UI Screens and State Flow

## Screens

- `HomeScaffoldScreen` (`lib/main.dart`): debug controls and state snapshot.
- `DiscoveryScreen` (planned): peer list sourced from `NearbyBridge.discoveredPeers()`.
- `TransferProgressScreen` (planned): progress updates from `NearbyBridge.transferProgress()`.

## Current Runtime Flow

1. App boots into `NearbyState.idle`.
2. `Initialize` button dispatches `InitializeRequested`.
3. `Discover` dispatches `StartDiscoveryRequested`.
4. `Select Peer` dispatches `PeerSelected(peer-001)`.
5. `Connected` dispatches `ConnectedEstablished`.
6. `Transfer` dispatches `TransferRequested(tx-001)`.
7. `Complete` dispatches `TransferCompleted(tx-001)`.
8. `Fail` dispatches `FailureObserved(message)`.
9. `Reset` dispatches `ResetRequested`.

## Integration Notes

- `MethodChannelNearbyBridge` in `lib/core/adapters/` implements the channel contract for both Android and iOS.
- Native layers must emit event payload keys exactly as declared in `contracts/nearby_bridge.schema.json`.
- Planned screens are blocked on external native bridge implementation and real device pairing tests.
