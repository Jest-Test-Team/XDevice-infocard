# Contracts

This folder contains minimal Solidity scaffolding for Plan 04:

- `src/SBTProfile.sol`: soulbound profile mint/update primitive.
- `src/ConnectionGraph.sol`: symmetric connection graph primitive.

## Foundry Notes

Once Foundry is installed, run these exact commands from this folder:

```bash
cd implementations/04-web3-sbt-contacts/contracts
forge --version
forge build
forge test -vv
forge test --match-path test/SBTProfile.t.sol -vv
forge test --match-path test/ConnectionGraph.t.sol -vv
```

Expected test files in this repo:

- `test/SBTProfile.t.sol`
- `test/ConnectionGraph.t.sol`
- `test/README.md`

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
