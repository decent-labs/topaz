pragma solidity 0.4.25;

import "./Ownable.sol";

contract ClientCapture is Ownable {
  event HashCaptured(
    bytes32 digest, uint8 hashFunction, uint8 size
  );

  function store(
    bytes32 digest, uint8 hashFunction, uint8 size
  )
    public
    onlyOwner
  {
    emit HashCaptured(digest, hashFunction, size);
  }
}
