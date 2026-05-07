sealed class NearbyEvent {
  const NearbyEvent();
}

class InitializeRequested extends NearbyEvent {
  const InitializeRequested();
}

class StartDiscoveryRequested extends NearbyEvent {
  const StartDiscoveryRequested();
}

class PeerSelected extends NearbyEvent {
  const PeerSelected(this.peerId);

  final String peerId;
}

class ConnectedEstablished extends NearbyEvent {
  const ConnectedEstablished();
}

class TransferRequested extends NearbyEvent {
  const TransferRequested(this.transferId);

  final String transferId;
}

class TransferCompleted extends NearbyEvent {
  const TransferCompleted(this.transferId);

  final String transferId;
}

class FailureObserved extends NearbyEvent {
  const FailureObserved(this.message);

  final String message;
}

class ResetRequested extends NearbyEvent {
  const ResetRequested();
}
