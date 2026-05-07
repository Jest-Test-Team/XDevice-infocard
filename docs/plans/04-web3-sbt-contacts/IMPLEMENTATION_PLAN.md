# Plan 04: Web3 Soulbound Token Contact Protocol

## 1. Objective
Build a verifiable contact-exchange protocol where business-card identity is anchored to Soulbound Tokens (ERC-5192), relationship events are written on-chain, and mobile UX is simplified via a Go relayer handling gas abstraction.

## 2. Scope
- In scope:
  - Solidity contracts for identity SBT and connection graph events.
  - Go relayer/indexer for submission and query APIs.
  - Mobile dApp client for wallet-based signing and handshake flow.
  - IPFS-backed metadata references.
- Out of scope (Phase 1):
  - Cross-chain support.
  - DAO governance mechanics.
  - Advanced zero-knowledge privacy overlays.

## 3. Proposed Architecture
- `contracts` (Solidity + Foundry/Hardhat):
  - `SBTProfile.sol` for identity mint and metadata CID reference.
  - `ConnectionGraph.sol` for signed physical-encounter attestations.
- `backend-relayer` (Go + go-ethereum):
  - Signature verification and nonce replay protection.
  - Sponsored transaction submission.
  - Event indexer into relational cache.
- `dApp-client` (React Native + WalletConnect/ethers):
  - Wallet auth and signing workflows.
  - Handshake UX driven by QR/BLE trigger.

## 4. Repository Blueprint
```text
web3-sbt-contacts/
├── contracts/
│   ├── src/
│   │   ├── SBTProfile.sol
│   │   └── ConnectionGraph.sol
│   ├── test/
│   ├── script/
│   └── foundry.toml
├── backend-relayer/
│   ├── cmd/server/main.go
│   ├── internal/api/
│   ├── internal/chain/
│   ├── internal/indexer/
│   ├── internal/store/
│   └── go.mod
├── dApp-client/
│   ├── src/
│   └── package.json
├── infra/
│   ├── docker-compose.yml
│   └── migrations/
└── README.md
```

## 5. Execution Phases
1. Contract design and tests (Week 1-2)
- Finalize event schema and access control rules.
- Implement SBT mint/lock semantics (ERC-5192 compliance).
- Implement connection write flow and unit tests.

2. Relayer MVP (Week 2-3)
- Build API endpoints for submission and status query.
- Verify EIP-712 signatures and enforce nonce freshness.
- Submit and monitor transactions on target chain.

3. Indexer and query layer (Week 3-4)
- Consume contract events and persist to Postgres.
- Serve connection history and profile lookup APIs.

4. Mobile dApp integration (Week 4-5)
- Wallet connect, sign payload, and submit via relayer.
- Handle pending/success/failure transaction states in UX.

5. Security hardening and staging launch (Week 5-6)
- Contract audits (internal + external).
- Abuse controls (rate limits, signature replay guards).
- Staging soak test with synthetic contact events.

## 6. On-chain Data Model
- `SBTProfile`:
  - owner address
  - metadata CID
  - minted timestamp
  - locked state
- `ConnectionGraph` event:
  - initiator address
  - responder address
  - shared nonce
  - encounter hash (time/location blinded hash)
  - block timestamp

## 7. API Surface (Relayer)
- `POST /v1/handshake/prepare`: create nonce and typed-data payload.
- `POST /v1/handshake/submit`: validate signatures and relay tx.
- `GET /v1/tx/:hash`: transaction status.
- `GET /v1/profile/:address`: indexed profile summary.
- `GET /v1/connections/:address`: paginated connection history.

## 8. Security Controls
- EIP-712 typed signatures.
- Nonce expiration and single-use enforcement.
- Rate limiting by wallet and IP.
- Optional server-side sanctions/blocklist checks.
- Encrypted off-chain metadata pointer strategy.

## 9. SLO / KPIs
- Relayer API p95 latency < 300ms (excluding chain confirmation).
- Successful relay submission >= 99.5%.
- Indexer lag < 2 blocks on average.

## 10. Delivery Artifacts
- Audited Solidity contracts + deployment scripts.
- Go relayer/indexer service + API docs.
- React Native dApp prototype.
- Ops runbook (key management, incident response, rollback strategy).
