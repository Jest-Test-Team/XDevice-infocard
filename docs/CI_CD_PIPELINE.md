# CI/CD Pipeline

## Goal
Remove local-environment blockers by shifting verification to deterministic CI jobs.

## Workflow
- File: `.github/workflows/verify-implementations.yml`
- Triggers:
  - `push` to `main`
  - `pull_request`
  - `workflow_dispatch`

## Jobs
1. `plan01-rust`
- Runs `cargo test` for `implementations/01-ultrasonic-data-transfer/core-dsp`.

2. `plan02-go`
- Runs `go test ./...` for signaling server.

3. `plan03-rust`
- Runs `cargo test` for `fountain-core`.
- Builds replay benchmark tool.
- Runs QR frame generator smoke check.

4. `plan04-go`
- Runs `go test ./...` for backend relayer.
- Runs bind-based relayer E2E script (`scripts/relayer_e2e.sh`) in GitHub runner.

5. `plan04-solidity`
- Installs Foundry with `foundry-rs/foundry-toolchain@v1`.
- Runs `forge build` and `forge test` for contracts.

## Local Fallback Scripts
- Solidity without local Foundry install:
  - `implementations/04-web3-sbt-contacts/scripts/solidity_test_via_docker.sh`
- Relayer E2E without port binding:
  - `implementations/04-web3-sbt-contacts/scripts/relayer_e2e_no_bind.sh`

## Blockers Addressed
- `forge/solc` missing locally: covered by CI Foundry job and Docker script.
- sandbox port-bind restriction locally: covered by CI bind-capable E2E and local no-bind fallback script.
