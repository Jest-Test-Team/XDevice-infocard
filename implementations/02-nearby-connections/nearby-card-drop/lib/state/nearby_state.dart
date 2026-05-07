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
  }) {
    return NearbyContext(
      state: state ?? this.state,
      activePeerId: activePeerId ?? this.activePeerId,
      activeTransferId: activeTransferId ?? this.activeTransferId,
      error: error ?? this.error,
    );
  }
}
