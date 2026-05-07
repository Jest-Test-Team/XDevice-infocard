# UniFFI Integration Checklist

This checklist covers the in-repo steps for generating bindings and integrating `core-dsp` with iOS and Android clients.

## 1. Generate bindings

1. Install generator CLI:
   - `cargo install uniffi_bindgen_cli`
2. From `implementations/01-ultrasonic-data-transfer/core-dsp` run:
   - `./scripts/generate-uniffi.sh all`
3. Verify generated outputs exist:
   - `uniffi/generated/swift/`
   - `uniffi/generated/kotlin/`

## 2. iOS integration (Swift)

1. Build Rust staticlib/XCFramework from `core-dsp`.
2. Add generated Swift sources from `uniffi/generated/swift/` to Xcode target.
3. Add Rust binary artifact to app/framework linking settings.
4. Validate calls:
   - `encode_payload(payload: String) -> Data`
   - `decode_payload(frame: Data) throws -> String`

## 3. Android integration (Kotlin/JNI)

1. Build Rust shared libs for Android ABIs (arm64-v8a, armeabi-v7a, x86_64).
2. Add generated Kotlin/JNI sources from `uniffi/generated/kotlin/`.
3. Package `.so` artifacts under `src/main/jniLibs/<abi>/`.
4. Validate calls:
   - `encode_payload(payload: String): ByteArray`
   - `decode_payload(frame: ByteArray): String`

## 4. External validation (still required)

1. Device-level audio playback/capture validation in real acoustic environments.
2. End-to-end latency and BER runs on multiple phones/OS versions.
3. Store/release hardening (signing, symbol stripping, CI artifact packaging).
