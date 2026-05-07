# App Scaffold Notes

`nearby-card-drop` is now structured as a minimal Flutter app skeleton:

- `pubspec.yaml` provides runnable Flutter metadata.
- `lib/main.dart` provides a minimal Material UI for state-machine transitions.
- `lib/state/` contains event/state/reducer flow.
- `lib/core/adapters/method_channel_nearby_bridge.dart` defines the native bridge adapter stub.

External generation still required:

- `flutter create .` must be run in this directory to generate full platform folders (`android/`, `ios/`, `macos/`, etc.) if a standalone app target is needed.
