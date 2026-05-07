// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title SBTProfile
/// @notice Minimal non-transferable profile token scaffold for Web3 contacts.
contract SBTProfile {
    struct Profile {
        string handle;
        string metadataURI;
        bool exists;
    }

    address public owner;
    uint256 public nextTokenId;

    mapping(uint256 => address) private _ownerOf;
    mapping(address => uint256) public tokenOf;
    mapping(uint256 => Profile) public profileOf;

    event ProfileMinted(address indexed to, uint256 indexed tokenId, string handle, string metadataURI);
    event ProfileUpdated(uint256 indexed tokenId, string handle, string metadataURI);

    error NotOwner();
    error Soulbound();
    error AlreadyHasProfile();
    error InvalidToken();

    modifier onlyOwner() {
        if (msg.sender != owner) revert NotOwner();
        _;
    }

    constructor() {
        owner = msg.sender;
        nextTokenId = 1;
    }

    function mintProfile(address to, string calldata handle, string calldata metadataURI) external onlyOwner returns (uint256 tokenId) {
        if (tokenOf[to] != 0) revert AlreadyHasProfile();

        tokenId = nextTokenId++;
        _ownerOf[tokenId] = to;
        tokenOf[to] = tokenId;
        profileOf[tokenId] = Profile({handle: handle, metadataURI: metadataURI, exists: true});

        emit ProfileMinted(to, tokenId, handle, metadataURI);
    }

    function updateProfile(uint256 tokenId, string calldata handle, string calldata metadataURI) external {
        if (_ownerOf[tokenId] == address(0)) revert InvalidToken();
        if (_ownerOf[tokenId] != msg.sender) revert NotOwner();

        profileOf[tokenId] = Profile({handle: handle, metadataURI: metadataURI, exists: true});
        emit ProfileUpdated(tokenId, handle, metadataURI);
    }

    function ownerOf(uint256 tokenId) external view returns (address) {
        address tokenOwner = _ownerOf[tokenId];
        if (tokenOwner == address(0)) revert InvalidToken();
        return tokenOwner;
    }

    function transferFrom(address, address, uint256) external pure {
        revert Soulbound();
    }

    function safeTransferFrom(address, address, uint256) external pure {
        revert Soulbound();
    }

    function safeTransferFrom(address, address, uint256, bytes calldata) external pure {
        revert Soulbound();
    }
}
