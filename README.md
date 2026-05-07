# XDevice-infocard

Cross-platform business-card exchange initiatives with plan docs and implementation scaffolds.

## Status
- Multi-agent implementation round completed on 2026-05-07.
- Coordinator verification passed for all Rust/Go modules under `implementations/*` that are runnable in this environment.
- CI/CD pipeline added to run full verification on GitHub runners, including Foundry Solidity tests and bind-based relayer E2E.

## CI/CD
- Workflow: `.github/workflows/verify-implementations.yml`
- Pipeline documentation: `docs/CI_CD_PIPELINE.md`

## Plan Docs
- Canonical plans: `docs/plans`
- Main monitoring dashboard: `docs/plans/IMPLEMENTATION_DASHBOARD.md`
- Progress tracker: `docs/plans/PROGRESS_TRACKER.md`

## Coordinator Docs
- `docs/AGENT_EXECUTION_STATUS.md`
- `docs/COMPLETION_REPORT.md`
- `docs/CI_CD_PIPELINE.md`

## Implementations
- `implementations/01-ultrasonic-data-transfer`
- `implementations/02-nearby-connections`
- `implementations/03-animated-qr-visual-handshake`
- `implementations/04-web3-sbt-contacts`

## Source
Initial strategy and expectations are defined in `spec&expectations.md`.
