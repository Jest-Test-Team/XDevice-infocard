# Failure Modes (Scaffold)

## Connectivity

- Bluetooth disabled/unavailable.
- Wi-Fi transport unavailable or blocked by system policy.
- Peer becomes unreachable during negotiation.

## Permission and Platform

- Runtime permission denied (Bluetooth, nearby devices, local network).
- Background execution limits interrupt discovery or transfer.
- Native API incompatibility across OS versions.

## Protocol and Session

- Connection handshake timeout.
- Peer rejects connection request.
- Payload checksum/length mismatch.
- Duplicate or stale transfer ID observed.

## Operational

- Signaling server unavailable (future rendezvous mode).
- Unexpected app process kill during active transfer.
- Partial file writes on interrupted receive.

## Handling Strategy (Initial)

- Map each failure to a typed error code for UI and retry policy.
- Push state machine into `error` with context for diagnostics.
- Add bounded retry for transient network conditions.
- Require explicit user action before resuming failed transfer.
