# Execution Breakdown (Plan 04)

## Goal
Deliver a production-ready Phase 1 Web3 contact protocol with:
- ERC-5192 profile SBT contracts
- connection attestation contract
- Go relayer + indexer
- mobile signing/submit integration
- staged rollout with measurable guardrails

## Workstreams and Owners
- Smart Contracts: Solidity engineer + security reviewer
- Relayer/API: Go backend engineer
- Indexer/Data: backend + data engineer
- Mobile dApp: React Native engineer
- DevSecOps: SRE/security engineer
- Program Control: tech lead / PM

## Milestones (6 Weeks)

### Week 1: Specification Freeze
- Finalize contract interfaces, events, errors, and role model.
- Finalize EIP-712 domain and typed-data schema.
- Define relayer API request/response contracts and error codes.
- Define chain targets (staging + prod), RPC providers, and confirmation policy.
- Exit criteria:
  - `CONTRACT_AND_REPLAYER_SPEC.md` approved by backend + mobile + contracts owners.
  - Threat model review completed (replay, frontrun, spam, key compromise).

### Week 2: Contract MVP + Unit Tests
- Implement `SBTProfile.sol` with ERC-5192 locked semantics.
- Implement `ConnectionGraph.sol` write path with nonce binding and duplicate prevention.
- Add Foundry tests for minting rules, lock behavior, signature validation, replay protection.
- Exit criteria:
  - Unit test coverage >= 90% of contract branches.
  - Gas report generated for core methods.

### Week 3: Relayer MVP + Integration Harness
- Implement `POST /v1/handshake/prepare` and `POST /v1/handshake/submit`.
- Server-side EIP-712 verification for both participants.
- Relayer transaction submission, pending tracking, and retry/backoff.
- Add Redis/Postgres nonce state (`issued`, `consumed`, `expired`).
- Exit criteria:
  - Deterministic replay rejection verified.
  - 99% success in synthetic submission tests (excluding chain outage).

### Week 4: Indexer + Query APIs
- Implement event consumers for profile and connection events.
- Persist canonical projections to Postgres.
- Implement `GET /v1/tx/:hash`, `GET /v1/profile/:address`, `GET /v1/connections/:address`.
- Add backfill/reorg handling and lag metrics.
- Exit criteria:
  - Reorg-safe indexing validated in testnet reorg simulation.
  - Query APIs satisfy pagination/consistency requirements.

### Week 5: Mobile Integration + Staging Soak
- Integrate WalletConnect flow and EIP-712 signing.
- Wire handshake prepare/submit/status UX states.
- Run staging soak tests with synthetic high-rate handshakes.
- Configure rate limits, abuse detection, and alerting.
- Exit criteria:
  - End-to-end handshake success >= 98% in staging.
  - p95 relayer prepare/submit latency < 300ms (chain confirmation excluded).

### Week 6: Security Hardening + Controlled Launch
- Internal security review + external audit fixes.
- Finalize runbook: relayer key rotation, incident response, and rollback.
- Progressive rollout gates (1%, 10%, 50%, 100%).
- Exit criteria:
  - All critical/high audit findings resolved or formally accepted.
  - Go/No-Go decision documented with KPI evidence.

## Detailed Execution Tasks

### Contracts
- Implement:
  - profile mint/update policy
  - non-transferability/locked checks
  - connection attestation write function
  - strict events + custom errors
- Security tasks:
  - domain separator chain binding
  - nonce uniqueness and expiry validation
  - input normalization (address ordering for connection pairs)

### Relayer
- Implement:
  - nonce issuance and TTL
  - signature validation for initiator/responder
  - sponsorship policy (who qualifies for gas subsidy)
  - tx broadcasting, confirmation, replacement policy
- Security tasks:
  - per-IP and per-wallet throttling
  - idempotency key support on submit endpoint
  - hot key isolation (KMS/HSM signer)

### Indexer
- Implement:
  - finalized and pending event processing modes
  - dedupe by `(tx_hash, log_index)`
  - reorg rollback depth and replay
- Reliability tasks:
  - dead-letter queue for malformed events
  - checkpointing and resumable consumers

### Mobile
- Implement:
  - deterministic typed-data rendering before sign
  - dual-sign handshake flow
  - explicit pending/confirmed/failed states
- Security tasks:
  - verify domain fields match expected chain + contracts
  - prevent blind signing by displaying decoded fields

## Rollout Plan
- Stage A (Internal): team wallets only, capped at 1k handshakes/day.
- Stage B (Pilot): allowlist cohorts, capped sponsorship budget.
- Stage C (Public): open registration with dynamic rate limits.
- Rollback triggers:
  - replay bypass observed
  - subsidy drain anomaly
  - chain/index mismatch > defined threshold

## Operational Readiness Checklist
- Secrets in managed KMS; no plaintext private keys in env files.
- On-call alerting configured for:
  - API error rate
  - tx stuck pending
  - indexer lag
  - sponsorship spend spikes
- Dashboards: relayer SLO, chain confirmation times, failed signature ratios.
- Runbooks tested via game-day drills.

## Definition of Done
- Contracts deployed to staging and production with verified source.
- Relayer and indexer passing CI, load tests, and security checks.
- Mobile flow validated on iOS/Android with target wallets.
- Documentation complete: API, incident response, rollback, and risk register.
