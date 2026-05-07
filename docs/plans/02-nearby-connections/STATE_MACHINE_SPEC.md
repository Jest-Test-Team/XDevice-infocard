# State Machine Spec

## Purpose
Define deterministic session behavior for Nearby Card Drop including transition rules, errors, and retries.

## Session Identity
- `sessionId`: UUID generated at transfer initiation.
- `endpointId`: Nearby endpoint identifier (peer-scoped).
- State machine key: `sessionId`.

## States
- `IDLE`: no active nearby operation.
- `DISCOVERING`: discovery/advertising active.
- `CONNECTING`: connection request/accept in progress.
- `AUTHENTICATING`: 4-digit short code confirmation pending.
- `TRANSFERRING`: chunk send/receive active.
- `VERIFYING`: digest/checksum validation in progress.
- `COMPLETED`: transfer succeeded; terminal.
- `FAILED`: transfer failed; terminal with retry path.
- `CANCELLED`: user cancelled; terminal.

## Events
- `START_DISCOVERY`
- `STOP_DISCOVERY`
- `ENDPOINT_FOUND(endpointId)`
- `CONNECT_REQUESTED(endpointId)`
- `CONNECTION_ACCEPTED(endpointId)`
- `CONNECTION_REJECTED(reason)`
- `CONNECTION_LOST(reason)`
- `AUTH_CODE_RECEIVED(code)`
- `AUTH_CONFIRMED`
- `AUTH_DECLINED`
- `TRANSFER_STARTED(totalChunks, digest)`
- `CHUNK_SENT(index)`
- `CHUNK_ACK_TIMEOUT(index)`
- `CHUNK_RECEIVED(index)`
- `TRANSFER_FINISHED`
- `VERIFY_OK`
- `VERIFY_FAIL(reason)`
- `SESSION_TIMEOUT`
- `USER_CANCEL`
- `RETRY_REQUESTED`
- `FATAL_ERROR(code, detail)`

## Transition Table
1. `IDLE` + `START_DISCOVERY` -> `DISCOVERING`
2. `DISCOVERING` + `ENDPOINT_FOUND` -> `DISCOVERING` (cache endpoint)
3. `DISCOVERING` + `CONNECT_REQUESTED` -> `CONNECTING`
4. `DISCOVERING` + `STOP_DISCOVERY` -> `IDLE`
5. `CONNECTING` + `CONNECTION_ACCEPTED` -> `AUTHENTICATING`
6. `CONNECTING` + `CONNECTION_REJECTED` -> `FAILED`
7. `CONNECTING` + `SESSION_TIMEOUT` -> `FAILED`
8. `AUTHENTICATING` + `AUTH_CONFIRMED` -> `TRANSFERRING`
9. `AUTHENTICATING` + `AUTH_DECLINED` -> `FAILED`
10. `AUTHENTICATING` + `SESSION_TIMEOUT` -> `FAILED`
11. `TRANSFERRING` + `TRANSFER_FINISHED` -> `VERIFYING`
12. `TRANSFERRING` + `CONNECTION_LOST` -> `FAILED` (retry-eligible)
13. `TRANSFERRING` + `SESSION_TIMEOUT` -> `FAILED`
14. `VERIFYING` + `VERIFY_OK` -> `COMPLETED`
15. `VERIFYING` + `VERIFY_FAIL` -> `FAILED`
16. `*` + `USER_CANCEL` -> `CANCELLED`
17. `*` + `FATAL_ERROR` -> `FAILED`
18. `FAILED` + `RETRY_REQUESTED` -> `DISCOVERING` (if retry budget remains)
19. `FAILED` + `RETRY_REQUESTED` -> `IDLE` (if retry budget exhausted)

## Guard Conditions
- `CONNECT_REQUESTED` allowed only when endpoint has been seen within last 10s.
- `AUTH_CONFIRMED` allowed only if both sides show identical 4-digit code.
- `TRANSFER_FINISHED` valid only if `receivedChunkCount == totalChunks`.
- `VERIFY_OK` valid only if local digest equals expected digest from handshake.

## Retry Rules

### Connection retry
- Budget: `CONNECTION_RETRY_LIMIT = 2`
- Backoff: `1s`, `2s` (linear)
- Trigger: `CONNECTION_REJECTED` with retryable reason or `CONNECTION_LOST` before transfer start.
- Non-retryable: permission denied, unsupported protocol version.

### Chunk retry
- Budget per chunk: `CHUNK_RETRY_LIMIT = 3`
- Timeout: `ACK_TIMEOUT_MS = 2000`
- Backoff: `250ms`, `500ms`, `1000ms`
- On exhaustion: emit `FATAL_ERROR(CHUNK_RETRY_EXHAUSTED)`.

### Session retry
- Budget: `SESSION_RETRY_LIMIT = 1`
- Eligible terminal states: `FAILED` caused by transient transport errors only.

## Error States and Codes
- `ERR_PERMISSION_DENIED`: required OS permissions missing.
- `ERR_ADAPTER_UNAVAILABLE`: Nearby service unavailable.
- `ERR_PROTOCOL_MISMATCH`: incompatible `protocolVersion`.
- `ERR_AUTH_TIMEOUT`: auth code not confirmed in timeout window.
- `ERR_AUTH_MISMATCH`: user codes differ.
- `ERR_CHUNK_TIMEOUT`: chunk ack timeout.
- `ERR_CHUNK_RETRY_EXHAUSTED`: repeated chunk failures.
- `ERR_DIGEST_MISMATCH`: payload integrity failure.
- `ERR_SESSION_TIMEOUT`: session exceeded wall-clock timeout.
- `ERR_USER_CANCELLED`: user cancelled operation.

## Timeouts
- Discovery passive timeout: 30s (UI prompts refresh/continue).
- Connect timeout: 10s.
- Authenticate timeout: 20s.
- Transfer idle timeout: 8s without chunk progress.
- Session hard timeout: 45s.

## Side Effects by State
- `DISCOVERING`: start discovery + advertising; clear stale endpoints.
- `CONNECTING`: freeze endpoint list; show connection progress UI.
- `AUTHENTICATING`: render shared 4-digit code, require explicit user action.
- `TRANSFERRING`: persist per-chunk progress and retry counters.
- `VERIFYING`: compute SHA-256 digest and CRC verification summary.
- `COMPLETED`: emit success metric + receipt; disconnect cleanly.
- `FAILED`: emit error metric + reason; disconnect and cleanup.
- `CANCELLED`: disconnect and purge session buffers.

## Instrumentation Requirements
Each transition log line must include:
- `sessionId`
- `fromState`
- `toState`
- `event`
- `endpointId`
- `errorCode` (nullable)
- `retryCount`
- `timestampMs`
