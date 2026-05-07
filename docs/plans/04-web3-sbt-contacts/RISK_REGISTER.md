# Risk Register (Plan 04)

## Scale
- Likelihood: Low / Medium / High
- Impact: Low / Medium / High / Critical
- Priority: P0 (highest) to P3 (lowest)

## Active Risks

| ID | Risk | Likelihood | Impact | Priority | Owner | Detection Signal | Mitigation | Contingency / Rollback Trigger |
|---|---|---|---|---|---|---|---|---|
| R-01 | Nonce replay bypass allows duplicate or forged connections | Medium | Critical | P0 | Smart Contract Lead | Repeat nonce seen in logs or chain events | Dual enforcement: off-chain nonce state + on-chain `nonceUsed`; short TTL; strict typed-data match | Immediate relayer pause if replay confirmed; rotate nonce namespace and patch contracts/relayer |
| R-02 | Signature validation mismatch across wallets (EIP-712 incompatibility) | Medium | High | P0 | Mobile Lead | Elevated `SIGNATURE_INVALID` by wallet type/version | Canonical typed-data generator; wallet compatibility matrix; preflight signer checks | Disable affected wallet versions via feature flag; fallback supported wallets only |
| R-03 | Relayer key compromise drains gas sponsorship | Low | Critical | P0 | Security/SRE | Unusual tx volume/spend; unknown destination patterns | KMS/HSM signer, no raw private key, strict IAM, spend caps, anomaly alerts | Revoke signer immediately, rotate key, disable sponsorship until postmortem complete |
| R-04 | Chain congestion causes long pending tx and poor UX | High | High | P1 | Backend Lead | `pending > threshold`, confirmation SLA breach | Dynamic fee strategy, replacement tx policy, multi-RPC routing | Shift to degraded mode with explicit pending notices; temporarily throttle new submissions |
| R-05 | Indexer reorg handling bug leads to incorrect connection history | Medium | High | P1 | Data/Indexer Lead | Divergence between chain and DB projections | Confirmations threshold, checkpoint rollback, idempotent replay by block range | Freeze query endpoints to last known consistent checkpoint and run backfill repair |
| R-06 | Sybil/spam wallets abuse free relay service | High | High | P1 | Abuse Prevention Owner | Rate-limit hits and subsidy spend spike | Per-IP/per-wallet rate limits, allowlist stage, anomaly scoring | Tighten limits, require API key/challenge, suspend sponsorship for offenders |
| R-07 | CID metadata contains sensitive data leakage | Medium | High | P1 | Product Security | Security review finds PII in sample metadata | Enforce schema validation, client-side warnings, encrypted payload pointers | Block updates with policy violations; purge indexed metadata cache |
| R-08 | Smart contract logic bug in soulbound constraints | Low | Critical | P0 | Smart Contract Lead | Audit findings, failing invariants | Internal + external audits, invariant tests, formal checks on transfer lock properties | Halt mint path with pause role; deploy patched version and migrate |
| R-09 | API abuse leads to relayer outage (DoS) | Medium | High | P1 | SRE Lead | Elevated 5xx, CPU saturation, queue growth | WAF, global concurrency limits, circuit breakers, autoscaling | Enable emergency strict limits and queue backpressure mode |
| R-10 | Dependency/RPC provider outage reduces availability | Medium | Medium | P2 | Platform Lead | Provider error rate > threshold | Multi-provider failover, health checks, timeout budgets | Route to backup provider set; if full outage, fail closed with clear status |
| R-11 | Regulatory/privacy concerns for storing encounter-derived hashes | Medium | High | P1 | Compliance Owner | Legal review flags data classification issue | Store only blinded hashes; no raw location/time; retention policy | Disable offending fields via schema version bump; update policy and re-consent |
| R-12 | Inadequate incident response delays containment | Medium | High | P1 | Engineering Manager | Slow MTTR in drills or incidents | On-call rotations, runbooks, game-day drills, clear severity matrix | Freeze rollout progression until runbook gaps are closed |

## Rollout Risk Controls
- Stage-gate progression is blocked unless prior stage KPIs and risk actions are complete.
- Mandatory Go/No-Go review at each gate (1% -> 10% -> 50% -> 100%).
- Error-budget policy:
  - rollback if handshake success < 98% for 30 minutes
  - rollback if replay/security anomaly confirmed
  - rollback if sponsorship spend exceeds daily cap by >20%

## Residual Risks (Accepted for Phase 1)
- No cross-chain portability.
- Limited privacy model (hash blinding only, no zk proofs).
- Wallet ecosystem variability may still produce edge-case signature UX failures.

## Review Cadence
- Weekly risk review during build phase.
- Daily risk check during rollout week.
- Post-incident review within 48 hours for Sev-1/Sev-2 events.
