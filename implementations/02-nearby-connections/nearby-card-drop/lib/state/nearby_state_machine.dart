import 'nearby_event.dart';
import 'nearby_state.dart';

/// Reducer-style state machine skeleton for Nearby transfer flow.
class NearbyStateMachine {
  NearbyStateMachine() : _context = const NearbyContext(state: NearbyState.idle);

  NearbyContext _context;

  NearbyContext get context => _context;

  NearbyContext dispatch(NearbyEvent event) {
    final current = _context;

    switch (event) {
      case InitializeRequested():
        _context = current.copyWith(state: NearbyState.initializing, error: null);
      case StartDiscoveryRequested():
        _context = current.copyWith(state: NearbyState.discovering, error: null);
      case PeerSelected(:final peerId):
        _context = current.copyWith(state: NearbyState.connecting, activePeerId: peerId, error: null);
      case TransferRequested(:final transferId):
        _context = current.copyWith(
          state: NearbyState.transferring,
          activeTransferId: transferId,
          error: null,
        );
      case TransferCompleted():
        _context = current.copyWith(state: NearbyState.completed);
      case FailureObserved(:final message):
        _context = current.copyWith(state: NearbyState.error, error: message);
    }

    return _context;
  }
}
