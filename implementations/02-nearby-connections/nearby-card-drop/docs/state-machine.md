# Nearby State Machine (Scaffold)

## States

- `idle`
- `initializing`
- `discovering`
- `advertising`
- `connecting`
- `connected`
- `transferring`
- `completed`
- `error`

## Primary Events

- `InitializeRequested`
- `StartDiscoveryRequested`
- `PeerSelected(peerId)`
- `ConnectedEstablished`
- `TransferRequested(transferId)`
- `TransferCompleted(transferId)`
- `FailureObserved(message)`
- `ResetRequested`

## Initial Transition Intent

- `idle -> initializing` when initialization begins.
- `initializing -> discovering` when local scan/advertise starts.
- `discovering -> connecting` when user selects a peer.
- `connecting -> connected` when the session is established.
- `connected -> transferring` when payload send/receive starts.
- `transferring -> completed` on successful completion.
- `* -> error` on transport, permission, or protocol failures.
- `* -> idle` on explicit reset, clearing active peer/transfer/error context.
