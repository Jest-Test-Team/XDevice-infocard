import 'package:flutter/material.dart';

import 'state/nearby_event.dart';
import 'state/nearby_state_machine.dart';

void main() {
  runApp(const NearbyCardDropApp());
}

class NearbyCardDropApp extends StatefulWidget {
  const NearbyCardDropApp({super.key});

  @override
  State<NearbyCardDropApp> createState() => _NearbyCardDropAppState();
}

class _NearbyCardDropAppState extends State<NearbyCardDropApp> {
  final NearbyStateMachine _machine = NearbyStateMachine();

  void _dispatch(NearbyEvent event) {
    setState(() {
      _machine.dispatch(event);
    });
  }

  @override
  Widget build(BuildContext context) {
    final contextState = _machine.context;

    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(title: const Text('Nearby Card Drop Scaffold')),
        body: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('state: ${contextState.state.name}'),
              Text('activePeerId: ${contextState.activePeerId ?? '-'}'),
              Text('activeTransferId: ${contextState.activeTransferId ?? '-'}'),
              Text('error: ${contextState.error ?? '-'}'),
              const SizedBox(height: 16),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: [
                  ElevatedButton(
                    onPressed: () => _dispatch(const InitializeRequested()),
                    child: const Text('Initialize'),
                  ),
                  ElevatedButton(
                    onPressed: () => _dispatch(const StartDiscoveryRequested()),
                    child: const Text('Discover'),
                  ),
                  ElevatedButton(
                    onPressed: () => _dispatch(const PeerSelected('peer-001')),
                    child: const Text('Select Peer'),
                  ),
                  ElevatedButton(
                    onPressed: () => _dispatch(const ConnectedEstablished()),
                    child: const Text('Connected'),
                  ),
                  ElevatedButton(
                    onPressed: () => _dispatch(const TransferRequested('tx-001')),
                    child: const Text('Transfer'),
                  ),
                  ElevatedButton(
                    onPressed: () => _dispatch(const TransferCompleted('tx-001')),
                    child: const Text('Complete'),
                  ),
                  ElevatedButton(
                    onPressed: () => _dispatch(const FailureObserved('mock failure')),
                    child: const Text('Fail'),
                  ),
                  OutlinedButton(
                    onPressed: () => _dispatch(const ResetRequested()),
                    child: const Text('Reset'),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
