import 'dart:async';

import 'package:flutter/services.dart';

import '../nearby_bridge.dart';
import '../nearby_types.dart';

class MethodChannelNearbyBridge implements NearbyBridge {
  MethodChannelNearbyBridge({
    MethodChannel? methodChannel,
    EventChannel? eventChannel,
  })  : _methodChannel = methodChannel ?? const MethodChannel(_methodsChannelName),
        _eventChannel = eventChannel ?? const EventChannel(_eventsChannelName);

  static const String _methodsChannelName = 'xdevice/nearby_bridge/methods';
  static const String _eventsChannelName = 'xdevice/nearby_bridge/events';

  final MethodChannel _methodChannel;
  final EventChannel _eventChannel;

  Stream<Map<String, dynamic>> get _events {
    return _eventChannel
        .receiveBroadcastStream()
        .where((raw) => raw is Map)
        .map((raw) => Map<String, dynamic>.from(raw as Map));
  }

  @override
  Future<void> initialize() => _invokeOk('initialize', const <String, Object?>{});

  @override
  Future<void> startAdvertising({required String localDisplayName}) {
    return _invokeOk('startAdvertising', <String, Object?>{
      'localDisplayName': localDisplayName,
    });
  }

  @override
  Future<void> stopAdvertising() => _invokeOk('stopAdvertising', const <String, Object?>{});

  @override
  Future<void> startDiscovery() => _invokeOk('startDiscovery', const <String, Object?>{});

  @override
  Future<void> stopDiscovery() => _invokeOk('stopDiscovery', const <String, Object?>{});

  @override
  Stream<PeerInfo> discoveredPeers() {
    return _events.where((event) => event['event'] == 'onPeerDiscovered').map((event) {
      final data = Map<String, dynamic>.from((event['data'] as Map?) ?? const <String, dynamic>{});
      return PeerInfo(
        id: (data['peerId'] as String?) ?? '',
        displayName: (data['displayName'] as String?) ?? '',
        medium: _parseMedium((data['medium'] as String?) ?? ''),
      );
    });
  }

  @override
  Future<void> requestConnection({required String peerId}) {
    return _invokeOk('requestConnection', <String, Object?>{'peerId': peerId});
  }

  @override
  Future<void> acceptConnection({required String peerId}) {
    return _invokeOk('acceptConnection', <String, Object?>{'peerId': peerId});
  }

  @override
  Future<void> rejectConnection({required String peerId}) {
    return _invokeOk('rejectConnection', <String, Object?>{'peerId': peerId});
  }

  @override
  Future<void> sendPayload({
    required String peerId,
    required TransferPayload payload,
  }) {
    return _invokeOk('sendPayload', <String, Object?>{
      'peerId': peerId,
      'transferId': payload.transferId,
      'fileName': payload.fileName,
      'totalBytes': payload.totalBytes,
    });
  }

  @override
  Stream<TransferProgressEvent> transferProgress() {
    return _events.where((event) => event['event'] == 'onTransferProgress').map((event) {
      final data = Map<String, dynamic>.from((event['data'] as Map?) ?? const <String, dynamic>{});
      return TransferProgressEvent(
        transferId: (data['transferId'] as String?) ?? '',
        bytesTransferred: (data['bytesTransferred'] as int?) ?? 0,
        totalBytes: (data['totalBytes'] as int?) ?? 0,
      );
    });
  }

  Future<void> _invokeOk(String method, Map<String, Object?> args) async {
    final dynamic raw = await _methodChannel.invokeMethod<dynamic>(method, args);
    final response = (raw is Map) ? Map<String, dynamic>.from(raw) : const <String, dynamic>{};
    if (response['ok'] != true) {
      throw PlatformException(
        code: 'BRIDGE_CALL_FAILED',
        message: 'Method $method did not return {ok:true}',
        details: response,
      );
    }
  }

  ConnectionMedium _parseMedium(String raw) {
    switch (raw) {
      case 'wifiLan':
        return ConnectionMedium.wifiLan;
      case 'wifiAware':
        return ConnectionMedium.wifiAware;
      case 'nearby':
      case 'multipeer':
      case 'bluetooth':
      default:
        return ConnectionMedium.bluetooth;
    }
  }
}
