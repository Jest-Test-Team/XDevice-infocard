#!/usr/bin/env python3
import argparse
import base64
import hashlib
import json
import sys
from typing import Iterable


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Chunk input and emit QR frame payload strings with metadata."
    )
    source = parser.add_mutually_exclusive_group(required=True)
    source.add_argument("--text", help="Plain UTF-8 text input")
    source.add_argument(
        "--base64", dest="base64_input", help="Base64-encoded binary input"
    )
    parser.add_argument(
        "--frame-bytes", type=int, default=64, help="Data bytes per frame (default: 64)"
    )
    parser.add_argument(
        "--transfer-id",
        help="Optional transfer id override (default derives from SHA-256 digest)",
    )
    parser.add_argument(
        "--format",
        choices=["jsonl", "plain"],
        default="jsonl",
        help="Output format (default: jsonl)",
    )
    return parser.parse_args()


def chunk_bytes(buf: bytes, width: int) -> Iterable[bytes]:
    for i in range(0, len(buf), width):
        yield buf[i : i + width]


def derive_transfer_id(payload: bytes) -> str:
    digest = hashlib.sha256(payload).hexdigest()
    return digest[:16]


def main() -> int:
    args = parse_args()
    if args.frame_bytes <= 0:
        print("frame size must be > 0", file=sys.stderr)
        return 2

    if args.text is not None:
        payload = args.text.encode("utf-8")
        input_mode = "text"
    else:
        try:
            payload = base64.b64decode(args.base64_input, validate=True)
        except Exception as exc:  # noqa: BLE001
            print(f"invalid base64: {exc}", file=sys.stderr)
            return 2
        input_mode = "base64"

    transfer_id = args.transfer_id or derive_transfer_id(payload)
    frames = list(chunk_bytes(payload, args.frame_bytes))
    total = len(frames)

    for seq, frame in enumerate(frames):
        frame_b64 = base64.b64encode(frame).decode("ascii")
        payload_string = f"xhv1:{transfer_id}:{len(payload)}:{seq}:{total}:{frame_b64}"
        if args.format == "plain":
            print(payload_string)
            continue

        row = {
            "transfer_id": transfer_id,
            "input_mode": input_mode,
            "payload_len": len(payload),
            "frame_bytes": args.frame_bytes,
            "sequence": seq,
            "total": total,
            "chunk_b64": frame_b64,
            "payload_string": payload_string,
        }
        print(json.dumps(row, separators=(",", ":")))

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
