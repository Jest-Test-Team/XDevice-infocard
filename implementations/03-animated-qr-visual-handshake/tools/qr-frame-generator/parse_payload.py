#!/usr/bin/env python3
"""Parse xhv1 frame payload strings into structured fields."""

import argparse
import base64
import json
from dataclasses import asdict, dataclass


@dataclass
class ParsedFrame:
    version: str
    transfer_id: int
    payload_len: int
    sequence: int
    total: int
    neighbors: list[int]
    chunk_hex: str
    chunk_b64: str


def parse_payload(payload: str) -> ParsedFrame:
    parts = payload.strip().split(":")
    if len(parts) != 7:
        raise ValueError("expected 7 colon-delimited fields")

    version, transfer_id, payload_len, sequence, total, neighbors_csv, chunk_b64 = parts
    if version != "xhv1":
        raise ValueError("unsupported version")

    neighbors = [int(v) for v in neighbors_csv.split(",") if v]
    if not neighbors:
        raise ValueError("neighbors field must not be empty")

    chunk = base64.b64decode(chunk_b64, validate=True)

    return ParsedFrame(
        version=version,
        transfer_id=int(transfer_id),
        payload_len=int(payload_len),
        sequence=int(sequence),
        total=int(total),
        neighbors=neighbors,
        chunk_hex=chunk.hex(),
        chunk_b64=chunk_b64,
    )


def main() -> None:
    parser = argparse.ArgumentParser(description="Parse xhv1 frame payload string.")
    parser.add_argument("payload", help="payload in xhv1 format")
    args = parser.parse_args()

    parsed = parse_payload(args.payload)
    print(json.dumps(asdict(parsed), ensure_ascii=True, indent=2))


if __name__ == "__main__":
    main()
