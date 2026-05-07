# Android Integration Notes

## Goal
Integrate `core-dsp` via UniFFI-generated Kotlin bindings for encode/decode flows.

## Planned steps

1. Add UniFFI generation to Rust build pipeline.
2. Generate Kotlin/JNI binding artifacts from `core-dsp/uniffi/uniffi.udl`.
3. Package Rust library for Android ABIs (arm64-v8a, armeabi-v7a, x86_64).
4. Wire encode/decode API into microphone/speaker pipeline.

## Expected API surface

- `encodePayload(payload: String): ByteArray`
- `decodePayload(frame: ByteArray): String`

## Status
Scaffold only. No generated Kotlin bindings committed yet.
