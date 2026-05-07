# Test Strategy (Plan 04)

## Objectives
- prove correctness of soulbound profile and handshake recording
- prevent replay/signature abuse
- validate relayer reliability under realistic load
- enforce rollout gates with measurable pass/fail criteria

## Test Pyramid
- Unit tests: contract logic, signature parsing, nonce store
- Integration tests: relayer + test chain + DB
- End-to-end tests: mobile signing flow through confirmation
- Non-functional tests: load, reorg resilience, fault injection, security checks

## 1. Contract Test Plan (Foundry)

## 1.1 SBTProfile
- mint success with valid CID
- mint rejects second profile for same owner
- transfer/approve methods revert (soulbound enforcement)
- `locked(tokenId)` returns true
- update CID allowed only for owner/authorized role

## 1.2 ConnectionGraph
- record succeeds with valid dual signatures and active SBTs
- reject when signatures mismatch signer
- reject when deadline expired
- reject self-connection
- reject duplicate pair+encounter hash
- reject nonce replay
- emit event fields exactly as specified

## 1.3 Property/Fuzz Tests
- fuzz address ordering and duplicate-key generation
- fuzz nonce uniqueness and consumed-state transitions
- fuzz deadline boundary (`now-1`, `now`, `now+1`)

## 1.4 Coverage and Quality Gates
- branch coverage >= 90%
- zero high-severity static-analysis issues (Slither/Mythril equivalent)
- gas snapshots checked into CI for regression alerts

## 2. Relayer/API Test Plan (Go)

## 2.1 Unit Tests
- typed-data canonicalization and hash stability
- signature recovery and normalization
- nonce store transitions: `issued -> consumed|expired`
- idempotency key semantics
- error mapping to API status codes

## 2.2 Integration Tests
- `prepare` then `submit` happy path produces tx hash
- duplicate `submit` with same idempotency key is safe/reused
- expired handshake rejected
- chain RPC timeout triggers retry policy
- tx revert path captured and surfaced in `GET /tx/:hash`

## 2.3 Database/Indexer Tests
- event ingest and projection consistency
- dedupe by `(tx_hash, log_index)`
- reorg simulation rollback + replay
- cursor pagination deterministic ordering

## 3. End-to-End Mobile Flow
- wallet connection and chain mismatch handling
- dual sign flow (initiator/responder)
- handshake submit and pending UX state
- confirmation polling and final success rendering
- recoverable failure UX (expired nonce, rate limit, revert)

## 4. Security Test Plan
- replay attack attempts (re-submit same nonce/signatures)
- typed-data tampering between prepare and submit
- malicious domain substitution attempts
- flood/spam tests for rate limiting by IP/wallet
- abuse test for sponsorship budget exhaustion
- secret scanning and dependency vulnerability checks in CI

## 5. Performance and Reliability

## 5.1 Load Targets
- `POST /handshake/prepare`: 200 rps sustained, p95 < 300ms
- `POST /handshake/submit`: 100 rps sustained, p95 < 300ms (without chain confirmation)
- indexer lag average < 2 blocks

## 5.2 Soak Test
- 24-hour continuous synthetic handshakes
- randomized signer pairs and nonce churn
- success rate >= 99.5% relay submission

## 5.3 Fault Injection
- RPC provider outage failover
- DB failover and reconnect recovery
- delayed confirmations and replacement transaction flow

## 6. Environments and Data
- Local: Anvil/Hardhat + ephemeral Postgres/Redis
- CI: deterministic seeded tests + coverage gates
- Staging: public testnet + production-like infra limits
- Production canary: allowlisted wallets only

Test data requirements:
- deterministic wallets for fixtures
- synthetic encounter hashes only (no raw location payloads)
- redact PII from logs and test artifacts

## 7. Release Gates
- Gate 1 (contracts): all unit/fuzz/security checks pass
- Gate 2 (relayer): integration + security replay suite pass
- Gate 3 (staging): 24h soak + KPI thresholds pass
- Gate 4 (production): canary error budget not exceeded for 72h

## 8. CI/CD Execution
- On every PR:
  - Solidity lint, unit, fuzz, gas snapshot diff
  - Go lint, unit, race detector, integration (mocked RPC)
- On merge to main:
  - full integration with testnet
  - security scans and dependency audits
- Nightly:
  - reorg/fault-injection suite
  - staged load benchmark trend report

## 9. Exit Criteria
Plan 04 is test-complete when:
- all release gates pass
- no open critical/high security findings
- performance targets are met in staging
- runbooks validated during at least one incident simulation
