# Plan 04: Web3 SBT Contacts

## Structure

- `contracts/`: Solidity contracts and notes for Foundry/Hardhat setup.
- `backend-relayer/`: Go HTTP relayer with in-memory request lifecycle and deterministic typed-data digest checks.
- `dApp-client/`: Frontend placeholder and implementation notes.
- `infra/`: Local infrastructure placeholders (`docker-compose.yml`).
- `docs/`: Operations notes for RPC migration.
- `scripts/`: local test and deployment templates.

## Completion Matrix

| Area | Status | Notes |
|---|---|---|
| Contracts scaffold (`SBTProfile`, `ConnectionGraph`) | Partial | Minimal implementation and tests exist; production hardening pending. |
| Foundry tests/docs | Partial | Unit tests + `forge` command docs present in `contracts/README.md` and `contracts/test/README.md`. |
| Relayer API (`/health`, `/v1/prepare`, `/v1/submit`, `/v1/tx/:hash`) | Done (local scope) | Request lifecycle, expiry, single-use, and tx status tracking implemented. |
| Signable payload validation | Done | Domain/type/message strict validation wired in submit flow. |
| EIP-712 digest computation | Done (deterministic hash verification) | Canonical digest generation from signable fields implemented in `backend-relayer/eip712.go`. |
| Wallet signature recovery | Not done | Deliberately out of scope for this pass; digest verification is implemented. |
| Chain RPC abstraction | Done | `ChainRPC` interface + default `InMemoryChainRPC` adapter implemented in `backend-relayer/rpc.go`. |
| RPC integration tests (no network) | Done | Mock RPC integration covered in relayer tests; no external infra required. |
| RPC-backed deployment path | Partial | Script and ops templates added; concrete JSON-RPC broadcaster implementation pending. |
| dApp client | Not done | Placeholder only. |
| Full infra/production rollout | Not done | Requires external chain infra and secret management. |

## Ops and Migration Docs

1. RPC migration guide: `docs/relayer-rpc-migration.md`
2. Deployment template: `scripts/deploy_relayer_rpc_template.sh`
3. Local relayer E2E checks:
   - `scripts/relayer_e2e.sh`
   - `scripts/relayer_e2e_no_bind.sh`
