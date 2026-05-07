import '../../core/nearby_types.dart';

abstract class TransferRepository {
  Future<void> queueOutgoing({
    required String peerId,
    required TransferPayload payload,
  });

  Stream<TransferStatus> watchTransferStatus(String transferId);
}

enum TransferStage {
  queued,
  negotiating,
  inProgress,
  completed,
  failed,
  canceled,
}

class TransferStatus {
  const TransferStatus({
    required this.transferId,
    required this.stage,
    required this.bytesTransferred,
    required this.totalBytes,
    this.errorCode,
  });

  final String transferId;
  final TransferStage stage;
  final int bytesTransferred;
  final int totalBytes;
  final String? errorCode;
}
