// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Test.sol";
import "../src/ConnectionGraph.sol";

contract ConnectionGraphTest is Test {
    ConnectionGraph internal graph;
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);
    address internal charlie = address(0xCAFE);

    function setUp() public {
        graph = new ConnectionGraph();
    }

    function testRequestConnectionEmitsEvent() public {
        vm.prank(alice);
        vm.expectEmit(true, true, false, false);
        emit ConnectionGraph.ConnectionRequested(alice, bob);
        graph.requestConnection(bob);
    }

    function testAcceptConnectionSetsSymmetricEdges() public {
        vm.prank(bob);
        graph.acceptConnection(alice);

        assertTrue(graph.isConnected(alice, bob));
        assertTrue(graph.isConnected(bob, alice));
    }

    function testRemoveConnectionClearsSymmetricEdges() public {
        vm.prank(bob);
        graph.acceptConnection(alice);

        vm.prank(alice);
        graph.removeConnection(bob);

        assertFalse(graph.isConnected(alice, bob));
        assertFalse(graph.isConnected(bob, alice));
    }

    function testSelfConnectionReverts() public {
        vm.prank(alice);
        vm.expectRevert(ConnectionGraph.SelfConnection.selector);
        graph.requestConnection(alice);

        vm.prank(alice);
        vm.expectRevert(ConnectionGraph.SelfConnection.selector);
        graph.acceptConnection(alice);
    }

    function testRemoveConnectionRevertsWhenNotConnected() public {
        vm.prank(charlie);
        vm.expectRevert(ConnectionGraph.NotConnected.selector);
        graph.removeConnection(alice);
    }
}
