# Android Bridge Contract

Defines Android-side contract for the Flutter method/event channels.

## Channels

- Method channel: `xdevice/nearby_bridge/methods`
- Event channel: `xdevice/nearby_bridge/events`

## Method Calls (Dart -> Android)

All method calls are invoked on the method channel with the name below and the JSON argument object.

### `initialize`

```json
{}
```

Returns:

```json
{"ok": true}
```

### `startAdvertising`

```json
{"localDisplayName": "alice-phone"}
```

Returns:

```json
{"ok": true}
```

### `stopAdvertising`

```json
{}
```

### `startDiscovery`

```json
{}
```

### `stopDiscovery`

```json
{}
```

### `requestConnection`

```json
{"peerId": "peer-123"}
```

### `acceptConnection`

```json
{"peerId": "peer-123"}
```

### `rejectConnection`

```json
{"peerId": "peer-123"}
```

### `sendPayload`

```json
{
  "peerId": "peer-123",
  "transferId": "tx-001",
  "fileName": "card.vcf",
  "totalBytes": 2048
}
```

## Event Payloads (Android -> Dart)

Event channel emits JSON objects with a top-level `event` and `data`.

```json
{
  "event": "onPeerDiscovered",
  "data": {
    "peerId": "peer-123",
    "displayName": "Bob",
    "medium": "nearby"
  }
}
```

Supported `event` values and `data` schema:

### `onPeerDiscovered`

```json
{"peerId": "peer-123", "displayName": "Bob", "medium": "nearby"}
```

### `onConnectionRequested`

```json
{"peerId": "peer-123"}
```

### `onConnectionStateChanged`

```json
{"peerId": "peer-123", "state": "connecting"}
```

`state` enum: `connecting | connected | rejected | disconnected | failed`

### `onTransferProgress`

```json
{"transferId": "tx-001", "bytesTransferred": 1024, "totalBytes": 2048}
```

### `onTransferCompleted`

```json
{"transferId": "tx-001"}
```

### `onError`

```json
{"code": "PERMISSION_DENIED", "message": "Nearby permission missing"}
```

## Notes

- Android implementation maps Google Nearby Connections callbacks into these event payloads.
- iOS must keep identical keys to preserve cross-platform Dart deserialization.
