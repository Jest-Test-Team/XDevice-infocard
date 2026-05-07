# Fountain Codec Specification (Plan 03)

## Purpose
Define the on-frame fountain chunk schema and deterministic encode/decode behavior for animated QR transport.

## Protocol Version
- `protocol_version`: `0x01` (Phase 1)
- Any incompatible wire-format or decode-semantic change must increment this value.

## Terminology
- Source block: Original payload split into fixed-size blocks.
- Symbol: XOR combination of one or more source blocks, identified by `symbol_index` and `symbol_seed`.
- Frame: One QR payload carrying one serialized symbol chunk.

## Source Payload Constraints (Phase 1)
- Payload size: `1KB` to `250KB` (hard cap).
- Recommended profile sizes for QA gates: `4KB`, `16KB`, `32KB`, `64KB`.
- Payload must be encrypted before fountain encoding.

## Fountain Chunk Schema (Binary)
Byte order: network/big-endian.

1. `protocol_version` (`u8`)
2. `flags` (`u8`)
- bit0: `is_completion_marker`
- bit1: `is_system_frame`
- bit2-7: reserved (must be 0 in v1)
3. `transfer_id` (`u128`)
- Random per session; must never repeat for active concurrent sessions.
4. `payload_sha256` (`[u8; 32]`)
- Hash of full plaintext or encrypted payload (must be consistent across sender/receiver policy).
5. `total_source_blocks` (`u16`)
- Number of source blocks in original payload.
6. `source_block_size` (`u16`)
- Default `512` bytes; last block may be shorter.
7. `symbol_index` (`u32`)
- Monotonic at sender, wraps disallowed in v1.
8. `symbol_seed` (`u64`)
- Seed for deterministic degree + block-selection generation.
9. `symbol_degree` (`u8`)
- Number of source blocks XORed for this symbol; must be `>=1`.
10. `symbol_len` (`u16`)
- Symbol payload bytes.
11. `symbol_bytes` (`[u8; symbol_len]`)
12. `symbol_crc16` (`u16`)
- CRC-16/CCITT-FALSE over fields `1..11`.

## Envelope for QR String
- Binary chunk is encoded with base64url (no padding) for portability.
- Final QR string format:
  - `XDI:1:<base64url_chunk>`
- Prefix rules:
  - `XDI` constant magic.
  - `1` is schema envelope version (independent of protocol byte).

## Sender Encode Flow
1. Validate payload size and compute `payload_sha256`.
2. Split payload into source blocks with configured `source_block_size`.
3. Generate base symbol budget:
- `target_symbols = ceil(total_source_blocks * redundancy_factor)`
- Default `redundancy_factor = 1.5`.
4. For each `symbol_index`:
- Derive `symbol_seed` from `(transfer_id, symbol_index)` via deterministic PRNG.
- Sample degree from robust distribution profile.
- Select unique source block indices deterministically from seed.
- XOR selected blocks into `symbol_bytes`.
- Build binary chunk and append CRC16.
- Serialize into envelope and emit as QR frame.
5. After normal symbols, emit completion marker frames (same transfer metadata, `is_completion_marker=1`) every `N=12` frames.

## Receiver Decode Flow
1. Scan QR and parse envelope prefix/version.
2. base64url decode payload; parse binary fields.
3. Validate:
- `protocol_version` supported.
- Reserved bits are zero.
- `transfer_id` matches active session.
- CRC16 valid.
- `symbol_len` sane for active profile.
4. Drop duplicates by `(transfer_id, symbol_index)`.
5. Reconstruct symbol equation from `symbol_seed` and `symbol_degree`.
6. Feed equation into decoder (peeling + fallback Gaussian elimination).
7. Update completion estimate: `recovered_source_blocks / total_source_blocks`.
8. On full recovery:
- Reassemble payload.
- Verify `payload_sha256` exact match.
- Emit transfer success event only after hash pass.

## Completion Criteria
Transfer is complete only when all are true:
1. Decoder solved all source blocks.
2. Reassembled payload length matches expected.
3. SHA-256 matches `payload_sha256` from chunk metadata.
4. No unresolved contradictions in decode graph.

## Failure and Recovery Rules
- CRC failure: drop frame, increment `crc_error_count`.
- Protocol mismatch: reject frame, increment `unsupported_version_count`.
- Transfer mismatch: ignore unless app explicitly supports multi-session.
- Hash mismatch after solved graph: mark transfer failed; restart session.
- Decoder stall timeout: if no rank increase for `1.2s`, request continued playback/replay.

## Interop/Determinism Requirements
- PRNG algorithm and seed derivation must be identical across platforms (Rust reference implementation is canonical).
- `symbol_seed` must fully determine degree and selected indices.
- Floating-point operations are forbidden in degree/index generation.

## Performance Requirements
- Serialize + encode overhead per symbol: <= `0.3ms` median on modern mobile CPU class.
- Decoder ingest path: sustain >= `45` symbols/s on mid-tier device.
- Decoder false completion probability: effectively `0` (guarded by SHA-256 validation).

## Validation Checklist
1. Golden vectors:
- Fixed transfer id + payload -> fixed first 256 symbols.
2. Round-trip tests:
- Out-of-order ingest and duplicate-heavy streams.
3. Corruption tests:
- Bit flips in header and payload must fail CRC/hash checks.
4. Stress tests:
- 10,000 random sessions with configured loss models.
