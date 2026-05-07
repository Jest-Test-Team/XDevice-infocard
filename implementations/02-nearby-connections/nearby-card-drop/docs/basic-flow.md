# Basic App Flow (Scaffold)

This is the minimal UI-to-state-machine wiring used by `lib/main.dart`.

## Event Sequence

1. `InitializeRequested` -> state `initializing`
2. `StartDiscoveryRequested` -> state `discovering`
3. `PeerSelected(peerId)` -> state `connecting`
4. `ConnectedEstablished` -> state `connected`
5. `TransferRequested(transferId)` -> state `transferring`
6. `TransferCompleted(transferId)` -> state `completed`
7. `ResetRequested` -> state `idle`

## Failure Path

- Any step can emit `FailureObserved(message)` -> state `error`
- UI can recover by dispatching `ResetRequested`

## Notes

- This scaffold is intentional: it validates reducer transitions without requiring native Nearby bridge integration yet.
- Native bridge callbacks should later dispatch the same events into `NearbyStateMachine`.
