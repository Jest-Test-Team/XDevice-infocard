# iOS Bridge Stub

This stub defines the iOS side of the Nearby method channel contract.

## Method Channel

- Channel name: `xdevice/nearby_bridge`

## Outbound Methods (Dart -> iOS)

- `initialize()`
- `startAdvertising(localDisplayName: String)`
- `stopAdvertising()`
- `startDiscovery()`
- `stopDiscovery()`
- `requestConnection(peerId: String)`
- `acceptConnection(peerId: String)`
- `rejectConnection(peerId: String)`
- `sendPayload(peerId: String, transferId: String, fileName: String, totalBytes: int)`

## Inbound Events (iOS -> Dart)

- `onPeerDiscovered(peerId: String, displayName: String, medium: String)`
- `onConnectionRequested(peerId: String)`
- `onConnectionStateChanged(peerId: String, state: String)`
- `onTransferProgress(transferId: String, bytesTransferred: int, totalBytes: int)`
- `onTransferCompleted(transferId: String)`
- `onError(code: String, message: String)`

## Notes

- iOS implementation should map MultipeerConnectivity callbacks into the event list above.
- Keep event naming and field keys aligned with Android to simplify Dart deserialization.
