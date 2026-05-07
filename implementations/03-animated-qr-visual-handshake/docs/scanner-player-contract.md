# Scanner/Player Integration Contract (v1)

This contract defines the lightweight frame payload format and scanner behavior for animated QR playback.

## Frame Payload Format

Each QR frame carries one UTF-8 string:

`xhv1:<transfer_id>:<payload_len>:<sequence>:<total>:<neighbors_csv>:<chunk_b64>`

Fields:
- `xhv1`: protocol/version marker.
- `transfer_id`: unsigned 64-bit decimal transfer identifier.
- `payload_len`: original payload byte length.
- `sequence`: symbol sequence number in `[0, total)`.
- `total`: total symbol count emitted by encoder.
- `neighbors_csv`: comma-separated source shard indices used to build this symbol.
  - Data symbols: one value equal to `sequence` (example: `7`).
  - Mixed repair symbols: multiple values (example: `0,3,5`).
- `chunk_b64`: standard Base64 payload bytes for the symbol.

## Receiver Contract

- Group symbols by `transfer_id`.
- Reject symbols with mixed `payload_len` or `total` inside the same transfer bucket.
- De-duplicate by `sequence`; conflicting duplicates are invalid.
- Decode once enough equations exist to recover all data symbols.
- On decode success, truncate to `payload_len` bytes.

## Player Contract

- Prefer repeating the full symbol loop at least 2-3 times.
- Maintain stable ordering by ascending `sequence` for deterministic scanner behavior.
- Keep frame dwell time stable (for example 60-120 ms) to reduce capture jitter.

## Error Handling

- Invalid format or parse failure: drop frame.
- Conflicting duplicates: discard transfer.
- Decode timeout (implementation-defined): surface retry prompt to user.
