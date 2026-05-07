# Risk Register

## Scale
- Probability: `Low | Medium | High`
- Impact: `Low | Medium | High`
- Exposure score: `1-9` (Probability x Impact mapping)

## Active Risks

| ID | Risk | Probability | Impact | Exposure | Owner | Mitigation | Trigger | Contingency |
|---|---|---|---|---:|---|---|---|---|
| R-01 | Android/iOS Nearby API behavior divergence causes inconsistent session events | High | High | 9 | Mobile Lead | Define strict event envelope + replay tests + platform adapters with normalization layer | Event ordering mismatch in integration tests | Force platform-specific guard branches and ship with narrowed supported flows |
| R-02 | Permission model differences reduce discovery/connect success | High | High | 9 | Android+iOS Engineers | Preflight permission checks, contextual prompts, denied-state UX path, telemetry by permission status | Discovery success < 90% in matrix | Block transfer entry until permission health passes and provide explicit remediation UI |
| R-03 | Chunk loss in noisy environments degrades transfer completion rate | Medium | High | 6 | Mobile Lead | Chunk retry policy, adaptive backoff, transfer timeout tuning, chunk-level metrics | Success rate < 95% for 5 MB tests | Reduce chunk size and enable signaling fallback for affected cohorts |
| R-04 | Integrity check false positives due to serialization mismatch | Medium | High | 6 | Protocol Owner | Canonical byte encoding spec, cross-platform golden vectors, digest tests | Digest mismatch > 0.5% in staging | Hotfix to canonical serializer and enforce protocol version bump |
| R-05 | Signaling fallback increases complexity and introduces regression risk | Medium | Medium | 4 | Backend Engineer | Feature-flag default off, isolate fallback path, dedicated regression suite | Direct-mode metrics regress after enabling flag | Disable flag remotely and route users to direct mode only |
| R-06 | Battery drain from continuous discovery/advertising harms UX | Medium | Medium | 4 | QA Lead | Discovery session timeout, adaptive scan intervals, background throttling | Battery drain > 8%/20 min test | Ship with conservative scan duty cycle and stricter idle stop |
| R-07 | Security weakness in auth code confirmation leads to accidental mispairing | Low | High | 3 | Security Champion | Mandatory 4-digit code confirmation, mismatch lockout, audit logs without PII | Any field report of wrong-recipient transfer | Disable transfer without explicit dual confirmation and release patch |
| R-08 | Observability gaps delay root-cause analysis in beta | Medium | Medium | 4 | SRE/Platform | Session-scoped structured logs + required transition fields + dashboard before beta | >10% failures with unknown cause | Hold rollout and add missing instrumentation before resume |
| R-09 | CI device matrix instability blocks release confidence | Medium | Medium | 4 | QA Lead | Reserve stable physical devices, nightly retries, flaky-test quarantine policy | >5% non-deterministic e2e failures | Manual certification run with trace capture and targeted infra fixes |
| R-10 | Protocol version drift between app builds causes connection failures | Medium | High | 6 | Mobile Lead | Enforce `protocolVersion` compatibility checks at handshake + release coordination | Spike in `ERR_PROTOCOL_MISMATCH` | Temporary compatibility shim for N-1 version, then force upgrade |

## Risk Review Cadence
- Weekly risk review during execution (every Monday).
- 48-hour review before each milestone gate.
- Daily review during release-candidate week.

## Quantitative Risk Thresholds
- Any risk with exposure >= 6 requires explicit mitigation task in sprint board.
- Any triggered high-impact risk (`Impact=High`) blocks milestone closure until containment is verified.

## Ownership Rules
- Each risk must have one directly accountable owner.
- Mitigation tasks must include due date and measurable outcome.
- Closed risks remain logged with closure evidence.
