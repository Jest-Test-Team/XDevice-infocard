# Technical Specification - Ultrasonic SonicCard Transfer

## 1. Scope
Defines the implementation-level protocol, frame format, module interfaces, data flow, and error handling for Plan 01.

## 2. Audio and Signal Parameters
- Carrier band: 18,000-20,000 Hz.
- Supported sample rates: 44,100 Hz and 48,000 Hz.
- Channel: mono.
- Sample format across FFI boundary: `f32` normalized to [-1.0, 1.0].
- Initial modulation: binary FSK with two tones per symbol.
- Symbol duration target: 10 ms (configurable by profile).

## 3. Session Protocol
- States:
  - `IDLE`
  - `DISCOVER`
  - `HANDSHAKE`
  - `TRANSMIT`
  - `VERIFY`
  - `ACK_WAIT`
  - `COMPLETE`
  - `FAILED`
- Sequence:
  1. Sender emits `DISCOVER` frame repeatedly for `discover_window_ms`.
  2. Receiver responds with `HANDSHAKE_ACK` containing session nonce.
  3. Sender transmits data frames in order (`seq` ascending).
  4. Receiver validates frame CRC and sequence; sends ACK/NACK bitmap.
  5. Sender retransmits missing frames until complete or retry limit reached.
  6. Receiver emits `TRANSFER_COMPLETE` with aggregate checksum.

## 4. Frame Format (On-Wire)
All multi-byte fields are little-endian.

| Field | Size (bytes) | Description |
|---|---:|---|
| Preamble | 4 | Fixed sync bytes `0xAA55AA55` |
| Version | 1 | Protocol version (`0x01`) |
| Frame Type | 1 | `DISCOVER`, `HANDSHAKE`, `DATA`, `ACK`, `NACK`, `COMPLETE`, `ERROR` |
| Session ID | 8 | Random per transfer |
| Packet ID | 4 | Random per packet stream |
| Sequence | 2 | Frame index starting at 0 |
| Total Frames | 2 | Total frame count for payload |
| Payload Length | 2 | Bytes in payload section |
| Flags | 1 | Bits: encrypted, final-frame, reserved |
| Payload | 0..512 | Frame body (chunked payload or control body) |
| CRC32 | 4 | CRC over `Version..Payload` |

Frame size target:
- Max encoded payload per `DATA` frame: 512 bytes.
- Max physical frame bytes (excluding modulation overhead): 541 bytes.

## 5. Payload Chunking Rules
- App payload max: 8KB.
- Split into chunks of up to 512 bytes.
- `Total Frames = ceil(payload_len / 512)`.
- Last chunk may be shorter; exact length in `Payload Length`.
- Missing or duplicate sequence numbers are handled by ACK/NACK bitmap.

## 6. Module Interfaces

### `framing.rs`
- Responsibility:
  - Serialize/deserialize frame structs.
  - CRC32 compute/verify.
- Core APIs:
  - `fn encode_frame(frame: &Frame) -> Result<Vec<u8>, FrameError>`
  - `fn decode_frame(bytes: &[u8]) -> Result<Frame, FrameError>`

### `protocol.rs`
- Responsibility:
  - Session state machine and retry policy.
  - Sequence bookkeeping.
- Core APIs:
  - `fn on_event(event: ProtocolEvent) -> Vec<ProtocolAction>`
  - `fn pending_retransmits(now_ms: u64) -> Vec<u16>`

### `modulator.rs`
- Responsibility:
  - Convert frame bytes into PCM samples using FSK symbols.
- Core APIs:
  - `fn modulate(frame_bytes: &[u8], cfg: ModemConfig) -> Result<Vec<f32>, ModemError>`

### `demodulator.rs`
- Responsibility:
  - Stream ingest PCM -> detect symbols -> reconstruct frame bytes.
- Core APIs:
  - `fn ingest(samples: &[f32]) -> DemodStatus`
  - `fn poll_frame() -> Option<Vec<u8>>`

### `lib.rs` (FFI boundary)
- Responsibility:
  - Stable external API for Swift/Kotlin clients.
- Required API:
  - `encode(payload: Vec<u8>) -> Vec<f32>`
  - `ingest(samples: Vec<f32>) -> DecodeState`
  - `poll_frame() -> Option<Vec<u8>>`
  - `reset_session(session_id: u64)`

## 7. End-to-End Data Flow
1. Mobile sender app serializes vCard bytes and encrypts payload (AES-GCM external module).
2. Rust protocol module chunks payload into `DATA` frames.
3. Framing module builds wire frame + CRC.
4. Modulator converts bytes to ultrasonic PCM.
5. Mobile playback pushes PCM to speaker.
6. Receiver microphone captures PCM, forwards buffers to Rust demodulator.
7. Demodulator reconstructs bytes, framing validates CRC.
8. Protocol module assembles chunks and emits ACK/NACK.
9. On completion, receiver returns reassembled payload to app layer.

## 8. Error Handling
- Error classes:
  - `FrameError`: bad preamble, unsupported version, CRC mismatch, invalid field.
  - `ModemError`: unsupported sample rate, symbol decode underflow, tone detection ambiguity.
  - `ProtocolError`: timeout, retry exhausted, session mismatch, replay detected.
  - `IOBridgeError`: audio interruption, permission denial, FFI marshaling mismatch.
- Handling policy:
  - CRC mismatch -> NACK specific sequence.
  - Sequence gap -> buffer out-of-order window and request retransmit.
  - Repeated demod sync loss -> trigger re-handshake.
  - Retry exhaustion -> fail session with error code and telemetry marker.
  - Permission/route interruption -> pause transfer and resume from last ACKed frame.

## 9. Logging and Privacy Constraints
- Never log raw vCard payload bytes.
- Log only session IDs (hashed), frame counts, retry counts, and error codes.
- Redact nonce and encryption metadata from user-visible logs.

## 10. Extension Points (Post-Phase 1)
- Optional PSK profile behind feature flag.
- Reed-Solomon FEC stage between framing and modulation.
- Device-specific EQ profile table for problematic microphones/speakers.
