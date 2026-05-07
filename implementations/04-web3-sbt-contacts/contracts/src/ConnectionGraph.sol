// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title ConnectionGraph
/// @notice Minimal directed connection graph scaffold for profile relationships.
contract ConnectionGraph {
    mapping(address => mapping(address => bool)) public isConnected;

    event ConnectionRequested(address indexed from, address indexed to);
    event ConnectionAccepted(address indexed from, address indexed to);
    event ConnectionRemoved(address indexed from, address indexed to);

    error SelfConnection();
    error NotConnected();

    function requestConnection(address to) external {
        if (to == msg.sender) revert SelfConnection();
        emit ConnectionRequested(msg.sender, to);
    }

    function acceptConnection(address from) external {
        if (from == msg.sender) revert SelfConnection();
        isConnected[from][msg.sender] = true;
        isConnected[msg.sender][from] = true;
        emit ConnectionAccepted(from, msg.sender);
    }

    function removeConnection(address peer) external {
        if (!isConnected[msg.sender][peer]) revert NotConnected();
        isConnected[msg.sender][peer] = false;
        isConnected[peer][msg.sender] = false;
        emit ConnectionRemoved(msg.sender, peer);
    }
}
