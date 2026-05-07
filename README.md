# XDevice-infocard

This repository organizes four implementation directions for cross-platform business-card exchange.

## Canonical Plan Location

All plans and execution artifacts are maintained under:
- `docs/plans`

Progress monitoring document:
- `docs/plans/PROGRESS_TRACKER.md`

## Plan Folders

1. `docs/plans/01-ultrasonic-data-transfer`
- Focus: Near-field ultrasonic transfer with shared Rust DSP core.
- Core docs:
  - `IMPLEMENTATION_PLAN.md`
  - `EXECUTION_BREAKDOWN.md`
  - `TECH_SPEC.md`
  - `TEST_STRATEGY.md`
  - `RISK_REGISTER.md`

2. `docs/plans/02-nearby-connections`
- Focus: Google Nearby Connections with Flutter + native bridges.
- Core docs:
  - `IMPLEMENTATION_PLAN.md`
  - `EXECUTION_BREAKDOWN.md`
  - `STATE_MACHINE_SPEC.md`
  - `TEST_STRATEGY.md`
  - `RISK_REGISTER.md`

3. `docs/plans/03-animated-qr-visual-handshake`
- Focus: Animated QR streaming with Rust fountain-code recovery.
- Core docs:
  - `IMPLEMENTATION_PLAN.md`
  - `EXECUTION_BREAKDOWN.md`
  - `CODEC_SPEC.md`
  - `TEST_STRATEGY.md`
  - `RISK_REGISTER.md`

4. `docs/plans/04-web3-sbt-contacts`
- Focus: Web3 SBT identity/contact graph with Go relayer.
- Core docs:
  - `IMPLEMENTATION_PLAN.md`
  - `EXECUTION_BREAKDOWN.md`
  - `CONTRACT_AND_REPLAYER_SPEC.md`
  - `TEST_STRATEGY.md`
  - `RISK_REGISTER.md`

## Source Note
The initial plan content is derived from `spec&expectations.md`, then expanded into implementation-ready execution documents.
