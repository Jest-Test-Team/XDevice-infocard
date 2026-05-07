# Android Bridge Stub

This stub defines the Android side of the Nearby method channel contract.

## Method Channel

- Channel name: `xdevice/nearby_bridge`

## Outbound Methods (Dart -> Android)

- `initialize()`
- `startAdvertising(localDisplayName: String)`
- `stopAdvertising()`
- `startDiscovery()`
- `stopDiscovery()`
- `requestConnection(peerId: String)`
- `acceptConnection(peerId: String)`
- `rejectConnection(peerId: String)`
- `sendPayload(peerId: String, transferId: String, fileName: String, totalBytes: int)`

## Inbound Events (Android -> Dart)

- `onPeerDiscovered(peerId: String, displayName: String, medium: String)`
- `onConnectionRequested(peerId: String)`
- `onConnectionStateChanged(peerId: String, state: String)`
- `onTransferProgress(transferId: String, bytesTransferred: int, totalBytes: int)`
- `onTransferCompleted(transferId: String)`
- `onError(code: String, message: String)`

## Notes

- Android implementation should map Google Nearby Connections callbacks into the event list above.
- Event payloads should remain schema-compatible with iOS bridge payloads.
