# QR Frame Generator

Runnable utility to chunk input into frame payload strings suitable for animated QR transport.

## Generator

```bash
python3 tools/qr-frame-generator/generate_frames.py --text "hello" --frame-bytes 6 --format jsonl
```

Output format `jsonl` emits rows with metadata + `payload_string` in this shape:

`xhv1:<transfer_id>:<payload_len>:<sequence>:<total>:<neighbors_csv>:<chunk_b64>`

## Parser Utility

```bash
python3 tools/qr-frame-generator/parse_payload.py "xhv1:123:5:0:2:0:aGVsbG8="
```

Parses the payload into structured JSON, including decoded chunk bytes (`chunk_hex`).
