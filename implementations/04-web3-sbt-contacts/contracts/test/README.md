# Foundry contract test scaffolding

This folder is reserved for Foundry tests that validate `SBTProfile` and `ConnectionGraph` behavior.

## Scope

- Unit tests for each contract in isolation.
- Integration tests for cross-contract flows (mint identity SBT, then create and update graph links).
- Revert-path tests for authorization and invalid state transitions.
- Event emission checks for key write operations.
- Gas snapshot checks on main state-changing functions.

## Suggested layout

- `contracts/test/SBTProfile.t.sol`
- `contracts/test/ConnectionGraph.t.sol`
- `contracts/test/Integration.t.sol`
- `contracts/test/helpers/Fixtures.sol` (shared setup/utilities)

## Concrete test plan

1. `SBTProfile`
- Mint succeeds for a new address and stores expected metadata.
- Mint fails when token already exists for address (soulbound uniqueness).
- Transfer-related methods revert (soulbound behavior enforcement).
- Metadata updates are restricted to authorized actor(s).
- Expected events are emitted with exact indexed parameters.

2. `ConnectionGraph`
- Create connection edge succeeds with valid identities.
- Duplicate edge creation reverts or is ignored (assert intended behavior).
- Remove/revoke edge updates state and emits event.
- Unauthorized actor cannot mutate another user's graph.
- Read methods return deterministic adjacency/state.

3. Integration
- Mint two SBT profiles, create a graph connection, and verify end-to-end state.
- Attempt graph mutation for an address without SBT and assert revert.
- Assert compatibility assumptions between contract interfaces.

4. Negative and fuzz
- Fuzz addresses and IDs for uniqueness and invariant checks.
- Fuzz ordering of add/remove operations for graph consistency.
- Invariant: no self-loop if protocol forbids it.

5. Gas checks
- Add `forge snapshot` baseline and fail CI on regressions above threshold.

## Example test case skeleton

```markdown
File: contracts/test/ConnectionGraph.t.sol

Contract: ConnectionGraphTest is Test

setUp()
- Deploy SBTProfile
- Deploy ConnectionGraph with SBTProfile dependency
- Mint profile for alice and bob

function test_CreateConnection_HappyPath() public
- vm.prank(alice)
- call createConnection(bob, "coworker")
- assertTrue(connectionGraph.isConnected(alice, bob))
- assertEq(connectionGraph.connectionType(alice, bob), "coworker")

function test_CreateConnection_RevertWhen_NotMinted() public
- vm.prank(charlie)
- vm.expectRevert(ConnectionGraph.ProfileNotFound.selector)
- connectionGraph.createConnection(bob, "coworker")
```

## Minimal commands

- Run all tests: `forge test -vvv`
- Run a single file: `forge test --match-path test/ConnectionGraph.t.sol -vvv`
- Run a single case: `forge test --match-test test_CreateConnection_HappyPath -vvv`
- Gas snapshot: `forge snapshot`

## CI recommendation

- Add a CI job that runs:
- `forge fmt --check`
- `forge test -vvv`
- `forge snapshot --check` (or compare against committed gas snapshot)
