# Plan 04: Web3 SBT Contacts (Scaffold)

This scaffold initializes the project layout for a soulbound-profile contact system.

## Structure

- `contracts/`: Solidity contracts and notes for Foundry/Hardhat setup.
- `backend-relayer/`: Minimal Go HTTP API with `health`, `prepare`, and `submit` placeholders.
- `dApp-client/`: Frontend placeholder and implementation notes.
- `infra/`: Local infrastructure placeholders (`docker-compose.yml`).

## Current Status

- Contracts are intentionally minimal and not production-ready.
- Relayer endpoints return placeholder responses.
- Frontend and infra are scaffolds only.

## Next Steps

1. Add contract tests and access-control refinements.
2. Implement relayer transaction building/signature verification/broadcast.
3. Build dApp flows and wallet/network UX.
4. Add local chain + CI checks.
