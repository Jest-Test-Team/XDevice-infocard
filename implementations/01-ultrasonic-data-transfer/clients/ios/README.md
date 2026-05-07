# iOS Integration Notes

## Goal
Integrate `core-dsp` via UniFFI-generated Swift bindings for encode/decode flows.

## Planned steps

1. Add UniFFI generation to Rust build pipeline.
2. Generate Swift binding artifacts from `core-dsp/uniffi/uniffi.udl`.
3. Package Rust static library or XCFramework for iOS consumption.
4. Wire encode/decode API into audio capture/playback pipeline.

## Expected API surface

- `encodePayload(payload: String) -> Data`
- `decodePayload(frame: Data) throws -> String`

## Status
Scaffold only. No generated Swift bindings committed yet.
