# Contracts

This folder contains minimal Solidity scaffolding for Plan 04:

- `src/SBTProfile.sol`: soulbound profile mint/update primitive.
- `src/ConnectionGraph.sol`: symmetric connection graph primitive.

## Foundry Notes

Suggested starter commands:

```bash
forge init contracts
cd contracts
forge build
forge test
```

You can keep these `src/*.sol` files and add tests under `test/`.

## Hardhat Notes

Suggested starter commands:

```bash
mkdir hardhat && cd hardhat
npm init -y
npm install --save-dev hardhat
npx hardhat
npx hardhat compile
npx hardhat test
```

Copy or reference these contracts from `contracts/src/` into your Hardhat workspace as needed.
