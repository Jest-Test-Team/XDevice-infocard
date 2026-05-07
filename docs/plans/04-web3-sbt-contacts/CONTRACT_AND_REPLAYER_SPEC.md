# Contract And Relayer Spec (Plan 04)

## 1. Scope
This document defines:
- on-chain contract interfaces/events for `SBTProfile` and `ConnectionGraph`
- relayer REST APIs
- end-to-end nonce + signature + submission flow

Phase 1 excludes cross-chain messaging and zk privacy overlays.

## 2. Chain and Signature Standards
- Token standard: ERC-5192 (minimal soulbound NFT interface)
- Signature standard: EIP-712 typed data
- Signature algo: secp256k1 (`eth_signTypedData_v4` compatible)
- Replay boundary: domain separator + per-handshake nonce + expiry

## 3. Smart Contract API

## 3.1 `SBTProfile` (Proposed)

### State
- `mapping(address => uint256) public profileTokenIdOf`
- `mapping(uint256 => string) private _tokenCid`
- `uint256 public nextTokenId`

### Roles
- `DEFAULT_ADMIN_ROLE`: admin controls (pause, role grants)
- `MINTER_ROLE`: allowlisted minter (optional if self-mint only)

### External Methods
- `function mintProfile(address to, string calldata cid) external returns (uint256 tokenId)`
- `function updateProfileCID(uint256 tokenId, string calldata cid) external`
- `function tokenURI(uint256 tokenId) public view returns (string memory)`
- `function locked(uint256 tokenId) external view returns (bool)`

### Behavior Rules
- one active profile token per wallet in Phase 1.
- transfer/approval operations MUST revert (soulbound).
- `locked(tokenId)` MUST always return `true` once minted.

### Events
- `event ProfileMinted(address indexed owner, uint256 indexed tokenId, string cid)`
- `event ProfileCIDUpdated(uint256 indexed tokenId, string oldCid, string newCid)`

### Custom Errors
- `error ProfileAlreadyExists(address owner)`
- `error Unauthorized()`
- `error InvalidCID()`

## 3.2 `ConnectionGraph` (Proposed)

### State
- `mapping(bytes32 => bool) public nonceUsed`
- `mapping(bytes32 => bool) public connectionRecorded`
  - key: `keccak256(abi.encodePacked(min(a,b), max(a,b), encounterHash))`

### Typed Data Struct
`HandshakeAttestation`:
- `address initiator`
- `address responder`
- `bytes32 encounterHash`
- `bytes32 nonce`
- `uint256 deadline`

Type hash:
- `keccak256("HandshakeAttestation(address initiator,address responder,bytes32 encounterHash,bytes32 nonce,uint256 deadline)")`

### External Methods
- `function recordConnection(HandshakeAttestation calldata a, bytes calldata sigInitiator, bytes calldata sigResponder) external`

### Validation Rules
- `block.timestamp <= deadline`
- `nonceUsed[a.nonce] == false`
- recovered signer from `sigInitiator` equals `a.initiator`
- recovered signer from `sigResponder` equals `a.responder`
- `a.initiator != a.responder`
- both addresses must own profile SBT tokens (Phase 1 gating)
- duplicate pair+encounter writes MUST revert

### State Effects
- mark nonce consumed
- mark connection key consumed
- emit connection event

### Events
- `event ConnectionRecorded(address indexed initiator, address indexed responder, bytes32 indexed nonce, bytes32 encounterHash, uint256 timestamp)`

### Custom Errors
- `error NonceAlreadyUsed(bytes32 nonce)`
- `error SignatureInvalid(address expected)`
- `error SignatureExpired(uint256 deadline, uint256 nowTs)`
- `error SelfConnectionNotAllowed()`
- `error DuplicateConnection(bytes32 key)`
- `error ProfileRequired(address wallet)`

## 4. Relayer API Specification

Base path: `/v1`
Format: JSON over HTTPS
Auth (Phase 1): wallet signatures only; optional API key for partner channels.

## 4.1 `POST /handshake/prepare`
Create canonical typed-data payload and server nonce.

Request:
```json
{
  "initiator": "0x...",
  "responder": "0x...",
  "encounter_hash": "0x...",
  "chain_id": 11155111
}
```

Response (200):
```json
{
  "handshake_id": "hs_01J...",
  "nonce": "0x...32bytes",
  "deadline": 1760000000,
  "typed_data": {
    "domain": {
      "name": "XDeviceConnectionGraph",
      "version": "1",
      "chainId": 11155111,
      "verifyingContract": "0xConnectionGraph"
    },
    "types": {"HandshakeAttestation": ["..."]},
    "primaryType": "HandshakeAttestation",
    "message": {"...": "..."}
  }
}
```

Validation and constraints:
- verify address format + checksum normalization.
- reject if initiator == responder.
- issue cryptographically random nonce (32 bytes).
- default deadline: now + 5 minutes.
- persist state: `issued` with TTL.

## 4.2 `POST /handshake/submit`
Accept both signatures, verify, and relay transaction.

Request:
```json
{
  "handshake_id": "hs_01J...",
  "attestation": {
    "initiator": "0x...",
    "responder": "0x...",
    "encounter_hash": "0x...",
    "nonce": "0x...",
    "deadline": 1760000000
  },
  "sig_initiator": "0x...",
  "sig_responder": "0x...",
  "idempotency_key": "7f0e..."
}
```

Response (202):
```json
{
  "relay_id": "rl_01J...",
  "tx_hash": "0x...",
  "status": "pending"
}
```

Server checks before broadcast:
- handshake exists and not expired/consumed.
- attestation payload exact-match with issued typed-data.
- both signatures recover expected signers.
- nonce not already used in local store.
- rate-limit and abuse checks pass.

Post-broadcast behavior:
- mark handshake `consumed` atomically with relay enqueue.
- record tx hash and confirmation target.
- retry replacement tx if pending beyond threshold.

## 4.3 `GET /tx/:hash`
Response:
```json
{
  "tx_hash": "0x...",
  "status": "pending|confirmed|failed|replaced",
  "block_number": 12345,
  "confirmations": 8,
  "revert_reason": null
}
```

## 4.4 `GET /profile/:address`
Response:
```json
{
  "address": "0x...",
  "token_id": "42",
  "cid": "bafy...",
  "minted_at_block": 120000
}
```

## 4.5 `GET /connections/:address`
Query: `cursor`, `limit` (max 100)

Response:
```json
{
  "items": [
    {
      "counterparty": "0x...",
      "encounter_hash": "0x...",
      "timestamp": 1760000000,
      "tx_hash": "0x..."
    }
  ],
  "next_cursor": "..."
}
```

## 5. Nonce and Signature Lifecycle
1. Client A/B agree encounter context off-chain (QR/BLE), derive `encounter_hash`.
2. App calls `POST /handshake/prepare` with A/B addresses + encounter hash.
3. Relayer returns canonical typed data with unique nonce and short deadline.
4. Both wallets sign the exact same EIP-712 message.
5. App submits payload + two signatures to `POST /handshake/submit`.
6. Relayer re-verifies signatures, nonce freshness, and anti-abuse constraints.
7. Relayer sends `recordConnection(...)` transaction.
8. On receipt/event, indexer writes canonical connection row.
9. Any replay of same nonce is rejected off-chain and on-chain.

## 6. Security Requirements
- Nonce entropy: CSPRNG, 256-bit.
- Nonce TTL: <= 5 minutes (configurable, default 300s).
- Clock skew tolerance: <= 30s server-side.
- Idempotency: required for client retries on submit.
- Domain pinning: chainId + verifyingContract must match configured environment.
- Reorg policy: only mark final after `N` confirmations (default 8 on EVM L2, 12 on L1).
- Relayer key management: KMS/HSM-backed signing, quarterly rotation minimum.

## 7. Failure Modes and Responses
- Invalid signature: `400 SIGNATURE_INVALID`
- Expired nonce/deadline: `400 NONCE_EXPIRED`
- Duplicate submit/idempotent replay: `200/202` with existing relay record
- Nonce already consumed on-chain: `409 NONCE_USED`
- Tx reverted: `422 TX_REVERTED`
- Chain unavailable: `503 RPC_UNAVAILABLE`

## 8. Observability Requirements
- Metrics:
  - `handshake_prepare_total`, `handshake_submit_total`
  - `signature_verify_fail_total`
  - `nonce_expired_total`, `nonce_replay_blocked_total`
  - `relay_tx_pending_seconds`, `relay_tx_fail_total`
- Structured logs include:
  - `handshake_id`, `relay_id`, `tx_hash`, `chain_id`, `error_code`
- Tracing across API -> relay queue -> chain submit -> indexer consume.
