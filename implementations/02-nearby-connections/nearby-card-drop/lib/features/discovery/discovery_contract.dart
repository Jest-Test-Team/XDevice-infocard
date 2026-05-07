import '../../core/nearby_types.dart';

abstract class DiscoveryRepository {
  Stream<List<PeerInfo>> watchDiscoveredPeers();
  Future<void> beginScan();
  Future<void> endScan();
}

abstract class ConnectionRequestHandler {
  Future<void> connectToPeer(String peerId);
  Future<void> approveIncoming(String peerId);
  Future<void> declineIncoming(String peerId);
}
