# iOS Integration Notes

## Goal
Integrate `core-dsp` via UniFFI-generated Swift bindings for encode/decode flows.

## In-repo helper path

1. Run `core-dsp/scripts/generate-uniffi.sh swift` to generate Swift bindings.
2. Follow `docs/uniffi-integration-checklist.md` for packaging/linking steps.
3. Package Rust static library or XCFramework for iOS consumption.
4. Wire encode/decode API into audio capture/playback pipeline.

## Expected API surface

- `encodePayload(payload: String) -> Data`
- `decodePayload(frame: Data) throws -> String`

## Status
Generation helper and checklist are committed. Platform packaging and app wiring are still pending.
