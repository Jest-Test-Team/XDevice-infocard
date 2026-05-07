# Bridge Contract Matrix

## Channels

| Type | Name |
| --- | --- |
| MethodChannel | `xdevice/nearby_bridge/methods` |
| EventChannel | `xdevice/nearby_bridge/events` |

## Method Calls (Dart -> Native)

| Method | Required args | Response |
| --- | --- | --- |
| `initialize` | none | `{ "ok": true }` |
| `startAdvertising` | `localDisplayName: string` | `{ "ok": true }` |
| `stopAdvertising` | none | `{ "ok": true }` |
| `startDiscovery` | none | `{ "ok": true }` |
| `stopDiscovery` | none | `{ "ok": true }` |
| `requestConnection` | `peerId: string` | `{ "ok": true }` |
| `acceptConnection` | `peerId: string` | `{ "ok": true }` |
| `rejectConnection` | `peerId: string` | `{ "ok": true }` |
| `sendPayload` | `peerId: string`, `transferId: string`, `fileName: string`, `totalBytes: int` | `{ "ok": true }` |

## Event Payloads (Native -> Dart)

| Event | Required data | Notes |
| --- | --- | --- |
| `onPeerDiscovered` | `peerId`, `displayName`, `medium` | `medium` supports `nearby`, `multipeer`, `bluetooth`, `wifiLan`, `wifiAware` |
| `onConnectionRequested` | `peerId` | incoming request hint |
| `onConnectionStateChanged` | `peerId`, `state` | `state`: `connecting`, `connected`, `rejected`, `disconnected`, `failed` |
| `onTransferProgress` | `transferId`, `bytesTransferred`, `totalBytes` | used for progress UI |
| `onTransferCompleted` | `transferId` | terminal success signal |
| `onError` | `code`, `message` | transport/permission failures |
