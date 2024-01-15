// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.3.2 (token/ERC20/presets/ERC20PresetMinterPauser.sol)

pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract ERC20NonTransferable is ERC20, AccessControl {
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    constructor(string memory name, string memory symbol, address minter) ERC20(name, symbol) {
        _setupRole(ADMIN_ROLE, minter);
    }

    function mint(address to, uint256 amount) public {
        require(hasRole(ADMIN_ROLE, _msgSender()), "must have admin role to mint");
        _mint(to, amount);
    }

    function transfer(address recipient, uint256 amount) public override returns (bool) {
        revert("non-transferable");
    }

    function transferFrom(address sender, address recipient, uint256 amount) public override returns (bool) {
        revert("non-transferable");
    }

    function burn(uint256 amount) public {
        _burn(_msgSender(), amount);
    }

    function burnFrom(address from, uint256 amount) public {
        require(hasRole(ADMIN_ROLE, _msgSender()), "must have admin role to burn");
        _burn(from, amount);
    }
}