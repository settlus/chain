# Abstract

The `x/nftownership` module enables the Settlus blockchain to easily get the owner of an NFT regardless of the chain the NFT is on.

# Contents
1. **[Concepts](#concepts)**
2. **[Queries](#queries)**
3. **[Events](#events)**
4. **[Parameters](#parameters)**

# Concepts

Nft ownership module allows users to easily query the owner of an NFT.
This module can query not only the NFTs on the Settlus chain, but also the NFTs on other chains by utilizing the [`x/oracle`](../../oracle/) module.

## Allowed Chain IDs
NFT ownership module has a set of allowed chain IDs.
Governance can modify the `allowed_chain_ids` parameter to add or remove chain IDs.

## Querying NFT Owner
Process for querying the owner of an NFT is different depending on the chain the NFT is on.

#### NFTs on Settlus Chain
`x/nftownership` module simply calls ERC721 contract's `ownerOf` function to query the owner of an NFT.

#### NFTs on Other Chains
`x/nftownership` module queries the owner of an NFT on other chains by utilizing the [`x/oracle`](../oracle/README.md) module.
[`x/oracle`](../oracle/README.md) module keeps track of the latest block number of the other chains.

# Queries

## Get NFT Owner
`/nftownership/get_nft_owner/{chain_id}/{contract_address}/{token_id}`

Request
```protobuf
message QueryGetNftOwnerRequest {
  string chain_id = 1;
  string contract_address = 2;
  uint64 token_id = 3;
}
```

Response

Returns hex address (e.g. `0xfoobar`) of the owner of the NFT.
```protobuf
message QueryGetNftOwnerResponse {
  string owner_address = 1;
}
```

# Events

TODO

The `x/nftownership` module emits the following events:

# Parameters

The `x/nftownership` module contains the following parameters:

| Key               | Type     | Example                             |
|-------------------|----------|-------------------------------------|
| allowed_chain_ids | []string | ["settlus_5371-1", "ethereum", ...] |

## Validations
- `allowed_chain_ids` can't have an empty string or duplicate chain IDs.
