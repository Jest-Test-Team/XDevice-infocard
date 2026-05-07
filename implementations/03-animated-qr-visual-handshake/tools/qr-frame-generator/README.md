# QR Frame Generator

Runnable utility to chunk input into frame payload strings suitable for animated QR transport.

## Usage

```bash
python3 tools/qr-frame-generator/generate_frames.py --text "visual-handshake" --frame-bytes 5
```

```bash
python3 tools/qr-frame-generator/generate_frames.py --base64 "AAECAwQFAA==" --format plain
```

Output format `jsonl` emits rows with metadata + `payload_string` in this shape:

`xhv1:<transfer_id>:<payload_len>:<sequence>:<total>:<chunk_b64>`
