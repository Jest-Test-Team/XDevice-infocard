// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Test.sol";
import "../src/SBTProfile.sol";

contract SBTProfileTest is Test {
    SBTProfile internal profile;
    address internal admin = address(0xA11CE);
    address internal alice = address(0xB0B);
    address internal bob = address(0xCAFE);

    function setUp() public {
        vm.prank(admin);
        profile = new SBTProfile();
    }

    function testMintProfileStoresData() public {
        vm.prank(admin);
        uint256 tokenId = profile.mintProfile(alice, "alice", "ipfs://alice");

        assertEq(tokenId, 1);
        assertEq(profile.tokenOf(alice), 1);
        assertEq(profile.ownerOf(1), alice);

        (string memory handle, string memory metadataURI, bool exists) = profile.profileOf(1);
        assertEq(handle, "alice");
        assertEq(metadataURI, "ipfs://alice");
        assertTrue(exists);
    }

    function testMintProfileRevertsIfAddressAlreadyHasProfile() public {
        vm.startPrank(admin);
        profile.mintProfile(alice, "alice", "ipfs://alice");
        vm.expectRevert(SBTProfile.AlreadyHasProfile.selector);
        profile.mintProfile(alice, "alice2", "ipfs://alice2");
        vm.stopPrank();
    }

    function testUpdateProfileOnlyOwner() public {
        vm.prank(admin);
        uint256 tokenId = profile.mintProfile(alice, "alice", "ipfs://alice");

        vm.prank(bob);
        vm.expectRevert(SBTProfile.NotOwner.selector);
        profile.updateProfile(tokenId, "alice-updated", "ipfs://alice-updated");

        vm.prank(alice);
        profile.updateProfile(tokenId, "alice-updated", "ipfs://alice-updated");

        (string memory handle, string memory metadataURI, bool exists) = profile.profileOf(tokenId);
        assertEq(handle, "alice-updated");
        assertEq(metadataURI, "ipfs://alice-updated");
        assertTrue(exists);
    }

    function testSoulboundTransfersAlwaysRevert() public {
        vm.prank(admin);
        uint256 tokenId = profile.mintProfile(alice, "alice", "ipfs://alice");

        vm.expectRevert(SBTProfile.Soulbound.selector);
        profile.transferFrom(alice, bob, tokenId);

        vm.expectRevert(SBTProfile.Soulbound.selector);
        profile.safeTransferFrom(alice, bob, tokenId);

        vm.expectRevert(SBTProfile.Soulbound.selector);
        profile.safeTransferFrom(alice, bob, tokenId, "");
    }
}
