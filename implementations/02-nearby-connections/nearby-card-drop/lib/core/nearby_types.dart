/// Shared types for nearby discovery and transfer flows.

enum ConnectionMedium {
  bluetooth,
  wifiLan,
  wifiAware,
}

enum TransferDirection {
  send,
  receive,
}

class PeerInfo {
  const PeerInfo({
    required this.id,
    required this.displayName,
    required this.medium,
  });

  final String id;
  final String displayName;
  final ConnectionMedium medium;
}

class TransferPayload {
  const TransferPayload({
    required this.transferId,
    required this.fileName,
    required this.totalBytes,
  });

  final String transferId;
  final String fileName;
  final int totalBytes;
}
