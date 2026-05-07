# Android Integration Notes

## Goal
Integrate `core-dsp` via UniFFI-generated Kotlin bindings for encode/decode flows.

## In-repo helper path

1. Run `core-dsp/scripts/generate-uniffi.sh kotlin` to generate Kotlin/JNI bindings.
2. Follow `docs/uniffi-integration-checklist.md` for ABI packaging/linking steps.
3. Package Rust library for Android ABIs (arm64-v8a, armeabi-v7a, x86_64).
4. Wire encode/decode API into microphone/speaker pipeline.

## Expected API surface

- `encodePayload(payload: String): ByteArray`
- `decodePayload(frame: ByteArray): String`

## Status
Generation helper and checklist are committed. ABI packaging and app wiring are still pending.
