# Relayer RPC Migration (In-Memory -> RPC-backed)

This plan upgrades the relayer from deterministic in-memory tx hash simulation to a real JSON-RPC-backed broadcast path.

## Current State

- `backend-relayer` verifies deterministic EIP-712 typed-data digest from signable fields.
- Relay submission uses `ChainRPC` interface.
- Default implementation is `InMemoryChainRPC` for local development/tests.

## Migration Steps

1. Implement a JSON-RPC client that satisfies `ChainRPC`:
   - New file suggestion: `backend-relayer/rpc_eth.go`
   - Implement `RelayContactOperation(ctx, req)`.
2. Build and sign on-chain tx payload:
   - Derive calldata for contract method (e.g. contact operation entrypoint).
   - Include EIP-155 chain ID from runtime config.
3. Broadcast and return canonical tx hash from node response.
4. Replace `buildChainRPC()` to select implementation by env:
   - `RELAYER_RPC_MODE=inmemory|ethereum`
5. Keep `InMemoryChainRPC` as default fallback for local no-infra testability.

## Required Environment (Template)

- `RELAYER_RPC_MODE=ethereum`
- `ETH_RPC_URL=https://...`
- `RELAYER_CHAIN_ID=...`
- `RELAYER_SIGNER_PRIVATE_KEY=0x...`
- `CONTACT_CONTRACT_ADDRESS=0x...`

## Operations Checklist

1. Health check: `GET /health`
2. Functional check: `POST /v1/prepare` then `POST /v1/submit`
3. Verify tx status endpoint output is persisted from real hash
4. Rollback path: set `RELAYER_RPC_MODE=inmemory` and restart

## CI Notes

- Unit tests must not hit external chain infra.
- Keep interface-level tests on `mockChainRPC`.
- Add optional integration job for private testnet only when secrets are present.
