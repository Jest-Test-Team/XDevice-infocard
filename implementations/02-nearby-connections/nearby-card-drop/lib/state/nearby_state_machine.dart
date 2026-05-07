import 'nearby_event.dart';
import 'nearby_state.dart';

/// Reducer-style state machine skeleton for Nearby transfer flow.
class NearbyStateMachine {
  NearbyStateMachine() : _context = const NearbyContext(state: NearbyState.idle);

  NearbyContext _context;

  NearbyContext get context => _context;

  NearbyContext dispatch(NearbyEvent event) {
    final current = _context;

    if (event is InitializeRequested) {
      _context = current.copyWith(state: NearbyState.initializing, error: null);
      return _context;
    }

    if (event is StartDiscoveryRequested) {
      if (current.state == NearbyState.initializing || current.state == NearbyState.idle) {
        _context = current.copyWith(state: NearbyState.discovering, error: null);
      }
      return _context;
    }

    if (event is PeerSelected) {
      if (current.state == NearbyState.discovering) {
        _context = current.copyWith(
          state: NearbyState.connecting,
          activePeerId: event.peerId,
          error: null,
        );
      }
      return _context;
    }

    if (event is TransferRequested) {
      if (current.state == NearbyState.connecting || current.state == NearbyState.connected) {
        _context = current.copyWith(
          state: NearbyState.transferring,
          activeTransferId: event.transferId,
          error: null,
        );
      }
      return _context;
    }

    if (event is TransferCompleted) {
      if (current.state == NearbyState.transferring && current.activeTransferId == event.transferId) {
        _context = current.copyWith(state: NearbyState.completed);
      }
      return _context;
    }

    if (event is FailureObserved) {
      _context = current.copyWith(state: NearbyState.error, error: event.message);
    }

    return _context;
  }
}
