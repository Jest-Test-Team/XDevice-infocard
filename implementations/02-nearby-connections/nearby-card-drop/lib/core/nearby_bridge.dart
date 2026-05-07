import 'nearby_types.dart';

/// Platform bridge interface implemented through method channels.
abstract class NearbyBridge {
  Future<void> initialize();

  Future<void> startAdvertising({required String localDisplayName});
  Future<void> stopAdvertising();

  Future<void> startDiscovery();
  Future<void> stopDiscovery();

  Stream<PeerInfo> discoveredPeers();

  Future<void> requestConnection({required String peerId});
  Future<void> acceptConnection({required String peerId});
  Future<void> rejectConnection({required String peerId});

  Future<void> sendPayload({
    required String peerId,
    required TransferPayload payload,
  });

  Stream<TransferProgressEvent> transferProgress();
}

class TransferProgressEvent {
  const TransferProgressEvent({
    required this.transferId,
    required this.bytesTransferred,
    required this.totalBytes,
  });

  final String transferId;
  final int bytesTransferred;
  final int totalBytes;
}
