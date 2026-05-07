enum NearbyState {
  idle,
  initializing,
  discovering,
  advertising,
  connecting,
  connected,
  transferring,
  completed,
  error,
}

class NearbyContext {
  const NearbyContext({
    required this.state,
    this.activePeerId,
    this.activeTransferId,
    this.error,
  });

  final NearbyState state;
  final String? activePeerId;
  final String? activeTransferId;
  final String? error;

  NearbyContext copyWith({
    NearbyState? state,
    String? activePeerId,
    String? activeTransferId,
    String? error,
    bool clearActivePeerId = false,
    bool clearActiveTransferId = false,
    bool clearError = false,
  }) {
    return NearbyContext(
      state: state ?? this.state,
      activePeerId: clearActivePeerId ? null : (activePeerId ?? this.activePeerId),
      activeTransferId: clearActiveTransferId ? null : (activeTransferId ?? this.activeTransferId),
      error: clearError ? null : (error ?? this.error),
    );
  }
}
