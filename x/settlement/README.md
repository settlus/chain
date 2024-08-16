---
sidebar_position: 1
---

# `x/settlement`

## Abstract

The `x/settlement` module, designed for Cosmos-SDK based blockchains, revolutionizes the management of financial transactions and settlements.

## Contents

* [Concepts](#concepts)
    * [Tenant](#tenant)
    * [Unspent Transaction Record](#unspent-transaction-record-utxr)
    * [Recipients](#recipients)
    * [Payout Period](#payout-period)
    * [Settlement](#settlement)
    * [Payout Method](#payoud-method)
    * [Fixed Fee](#fixed-fee)
* [State](#state)
    * [Unspent Transaction Record](#unspent-transaction-record)
* [Messages](#messages)
    * [MsgRecord](#msgrecord)
    * [MstCancel](#MstCancel)
    * [MsgDepositToTreasury](#msgdeposittotreasury)
* [Events](#events)
    * [EventRecord](#EventRecord)
    * [EventCancel](#EventCancel)
    * [DepositToTreasury](#deposit-to-treasury)
* [Parameters](#parameters)
* [Client](#client)
    * [CLI](#cli)
    * [gRPC](#grpc)
    * [REST](#rest)

## Concepts

### Tenant
In the Settlus blockchain, the concept of a "Tenant" represents an individual platform or service that utilizes the `x/settlement` module.
Each tenant operates independently within the Settlus ecosystem, maintaining their distinct transaction records, revenue streams, and user interactions.
The tenant is identified by a unique `tenant_id`, which serves as a cornerstone for tracking and managing their specific transactions and settlements.

#### Tenant Admin
A crucial aspect of each tenant's structure is the designation of Tenant Admins.
These admins are authorized individuals or entities who have the capability to execute administrative actions specific to their tenant.
This includes managing settlement treasury, recording revenues, cancelling UTXRs, and overseeing settlements.
The list of Tenant Admins is configured and maintained in the `admins` parameter, ensuring secure and authorized access to administrative functions.

### Unspent Transaction Record (UTXR)
The Unspent Transaction Record (UTXR) is a fundamental element of the Settlus blockchain, created whenever a payment is made from a tenant to a recipient.
These records are crucial in tracking the flow of funds and ensuring the accuracy of settlements.
Each UTXR contains details such as the NFT address, recipient's address, the amount of the transaction in the `unspent_records` state, these records are the backbone of the settlement process.

### Recipients
The `Recipients` field contains a list of recipient addresses and their corresponding weights, representing the owners of an NFT. When a record is settled, the amount is split by weight and sent to the recipients. 

- If the NFT is stored in Settlus, we directly determine the recipients of the NFT during the execution of a record transaction.
- If the NFT is stored on an external chain, we postpone determining the owners. The oracle module will later fill in these details through voting from feeders.

### Payout Period
The "Payout Period" in the Settlus blockchain is a critical concept pertaining to the lifecycle of an Unspent Transaction Record (UTXR).
This period is defined as a specific number of blockchain blocks and represents the timeframe during which a UTXR is eligible for a cancel.
The length of the payout period is pre-established and is integral in dictating the conditions under which a transaction can be reversed.

During the payout period, if the involved parties decide that a cancel is necessary, the UTXR can be canceled by simply removing the UTXR from the `unspent_records` state.
This mechanism provides a safety net for both tenants and recipients, allowing for transaction disputes or errors to be rectified within a reasonable time frame.

Once the payout period has elapsed, the UTXR is no longer eligible for a cancel.
At this juncture, the transaction is considered final and the funds are ready to be disbursed.
The recipient will then receive the payment in USDC, the designated settlement currency.
This process ensures a clear and structured settlement timeline, providing certainty and transparency to both tenants and recipients in the transaction process.

#### Payout Period Example
`payout_period` is measured in blocks.

> Suppose `payout_period` is set to 201600 and average block period is 3 seconds.
> A `payout_period` of 201600 (60 * 60 * 24 * 7 / 3) which is about a week, means that the UTXR can be canceled for 201600 blocks after the UTXR is created.
> After 201600 blocks, the UTXR is considered settled and cannot be canceled.

### Settlement
Once the payout period of a UTXR has concluded, the UTXR is eligible for settlement.
This process is the final step in the lifecycle of a transaction and is essential for the actual transfer of funds from tenants to recipients.

At each `BeginBlock`, the `x/settlement` module iterates through the `unspent_records` state and checks if the UTXR's payout period has passed.
The UTXRs that have passed the payout period are considered eligible for settlement.
If the tenant has enough funds in the treasury to settle the UTXR, the UTXR is removed from the `unspent_records` state and the amount is transferred from the tenant's treasury to the recipient's wallet.
If the tenant does not have enough funds in the treasury to settle the UTXR, the settlement will be deferred until the tenant has enough funds.
The UTXR will remain in the `unspent_records` state until the tenant has enough funds to settle the UTXR.

### Payout Method
When creating a Tenant in the Settlement Module, one of two types of Payout Methods must be selected: Native, Mintable Contract

#### Native
When the Payout Method is set to Native, an existing token on the blockchain is used as the settlement currency. There is a Treasury where settlement reserves are held, and after the Payout Period ends, the Settlement Module deletes the UTXR and transfers the tokens from the Treasury to the Recipient. The Settlement Module is the only entity with transfer authority over each Tenant’s Treasury.

The Tenant Admin must ensure that the Treasury has sufficient tokens to avoid any delays in settlement. This can be managed through the deposit_to_treasury function. If the settlement currency is registered as a Native/ERC20 Token Pair via the ERC20 Module, it will ultimately be converted and transferred as an ERC20 token.

#### Mintable Contract
When the Payout Method is set to Mintable Contract, an ERC20 token with a mint(address to, uint256 amount) function is used as the settlement currency. In this case, there is no separate Treasury; instead, after the Payout Period ends, the Settlement Module mints new ERC20 tokens and transfers them to the Recipient.

If a specific Contract address is not provided when creating a Tenant, a new Contract is generated, which defaults to a Soul-Bound Token (Non-Transferable Token). In this case, the Settlement Module holds the exclusive minting authority. If a specific Contract address is provided, the Settlement Module must be granted the authority to execute the mint function.

### Fixed Fee
It is common for the price of a coin to fluctuate significantly due to external factors, regardless of the supply and demand related to its actual use. Such fluctuations are more frequent before the blockchain stabilizes. If such events occur, the cost required to record a transaction could fluctuate significantly, which could be unfavorable for the creators and platform services using Settlus. To avoid this, transactions handled by the Settlement Module are paid with a fixed gas amount and gas price, such as 0.001 USDC

## State

### Unspent Transaction Record

UTXR data structure is designed with the following requirements.
- Fast look up by the oldest payout period.
- Fast look up & deletion by the UTXR ID.
- Fast insert.

To meet the requirements above, we use two stores: `unspent_records` and `utxr_by_request_id`.

`UTXR` contains unspent transaction records for a tenant.

- `UTXR`: `((TenantID)-(UTXRID)) -> UTXR`
```go
struct {
  UtxrId uint64
  RequestId string
  Recipients []*Recipient
  Nft    Nft
  Amount sdk.Coins
  CreatedAt uint64
}
```
The UTXR ID is incremented by 1 for each UTXR. Because the UTXRs are created in order, the UTXRs are trivially sorted by `UtxrId`.


### Unspent Transaction Record by Request ID
There is another store that contains UTXR IDs by request ID.
This store is used to help fast look up of UTXRs by request ID.

The request ID is generated by the tenant and is used to identify the UTXR when the tenant wants to query or cancel the UTXR.

- `UTXRByRequestId`: `RequestId -> UTXRId`
```go
struct {
  UTXRId bytes
}
```

#### Example Scenarios
**Insertion**
- A new UTXR is created with the following values.
    - `Recipients`: `List of recipients with weight`
    - `Amount`: `Amount`
    - `CreatedAt`: `CurrentBlockHeight`
- Add the new UTXR to the store with the following key: `((TenantID)-(UTXRID))`.

**Settlement**
- Iterate from the lowest key to the highest key in the `unspent_records` store.
- For each UTXR, do the following:
  - Check if the current UTXR's `CreatedAt + PayoutPeriod` is less than or equal to the current block height.  
    - If no, since every UTXR is sorted by `CreatedAt + PayoutPeriod`, we can stop checking the rest of the UTXRs.
    - If yes, check if the tenant has enough funds to settle the UTXR.
      - If the tenant has enough funds, remove the UTXR from the `unspent_records` state and transfer the amount from the tenant's treasury to the recipient's wallet.
      - If the tenant does not have enough funds, stop the iteration and emit a `NotEnoughTreasuryBalance` event.

**Cancel**
In the case of a cancel, the UTXR is simply removed from the `unspent_records` state.
- Get the UTXR ID by Request ID.
- Delete the UTXR from the `unspent_records` state by the UTXR ID.

## End Block
At each `EndBlock`, the `x/settlement` module iterates through the `unspent_records` state and checks if the UTXR's payout period has passed.

## Messages

### MsgRecord
The `MsgRecord` message allows tenant admins to record revenue.
`x/settlement` module will query the owner of the NFT with `x/nftownership` module and record the revenue to the owner's wallet.
```go
type MsgRecord struct {
	Sender string
    TenantId string
	RequestId string
    Amount sdk.Coin
    ChainId string
	ContractAddress string
	TokenId string
	metadata string
}
```
`metadata` is an optional field that can be used to send additional information about the UTXR.
`metadata` is not stored in the `unspent_records` state, but it is emitted in the `EventUTXRCreated` event.
An indexer can parse the emitted `metadata` and store it in a separate database.

### MsgCancel
The `MsgCancel` message allows tenant admins to cancel a UTXR.
```go
type MsgCancel struct {
	Sender string
    TenantId string
	RequestId string
}
```

### MsgDepositToTreasury
The `MsgDepositToTreasury` message allows tenant admins deposit funds to the treasury.
Anyone can deposit funds to the treasury.
```go
type MsgDepositToTreasury struct {
	Sender string
    Id string
    TenantId uint64
    Amount sdk.Coin
}
```

## Events
The `x/settlement` module emits the following events:


### EventRecord
```go
type EventRecord struct {
    TenantId uint64
	UtxrId uint64
	RequestId string
    Recipients []*Receipinet
	Nft Nft
	Amount sdk.Coins
	Metdata string
}
```

### EventSettled
```go
type EventSettled struct {
	Tenant uint64
	UtxrId uint64
}
```

### EventCancel
```go
type EventCancel struct {
    TenantId string
	RequestId string
}
```

### EventSetRecipients
```go
type EventSetRecipients struct {
    Tenant uint64
	UtxrId uint64 
	Recipients []*Recipient
}
```

## Parameters
The `x/settlement` module contains the following parameters:

| Key                        | Type      | Example                      |
|----------------------------|-----------|------------------------------|
| gas_prices                 | []DecCoin | [{denom: setl, amount: 0.1}] |
| oracle_fee_percentage      | dec       | "0.500000000000000000"       |
| supported_chains           | []Chain   | [{chain_id:1, chain_name: ethereum, chain_url: https://ethereum.org}]|


## Client

### CLI
A user can query and interact with the `x/settlement` module using the CLI.

#### Query
The query commands allow users to query `x/settlement` module's state.
```shell
settlusd query settlement --help
```

##### UTXRs
The `utxrs` command allows users to query the UTXRs of a tenant.
```shell
settlusd query settlement utxrs [tenant-id] [flags]
```

Example:
```shell
settlusd query settlement utxrs 1
```

Example Output:
```shell
{
  "id": "1",
  "tenant_id": "1",
  "created_at": "550",
  "recipients": [
    {
        "addr": "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5", 
        "weight": 1
    }
  ]
  "amount": [
    {
        "denom": "uusdc",
        "amount": "1000000"
    }
  ],
}
```


#### Transactions
The `tx` commands allow users to interact with the `x/settlement` module.
```shell
settlusd tx settlement --help
```

##### Record
The command `record` allows tenant admins to record a UTXR.

Usage:
```shell
settlusd tx settlement record [tenant-id] [request-id] [amount] [chain-id] [contract-address] [token-id] [metadata] [flags] 
```

Example:
```shell
settlusd tx settlement record \
    1 # tenant id \
    request-1 # request id \
    1000000usdc # amount
    1 # chain id \
    0x0000000000000000000000000000000000000001 # contract address \
    0x1 # token id \
    "metadata" # metadata
```

##### Cancel
The command `cancel` allows tenant admins to cancel a UTXR.

Usage:
```shell
settlusd tx settlement cancel [tenant-id] [request-id] [flags]
```

Example:
```shell
settlusd tx settlement cancel \
    1 # tenant id \
    request-1 # request id
```

##### Deposit to Treasury
The command `deposit-to-treasury` allows tenant admins to deposit funds to the treasury.

Usage:
```shell
settlusd tx settlement deposit-to-treasury [tenant-id] [amount] [flags]
```

Example:
```shell
settlusd tx settlement deposit-to-treasury \
    1 # tenant id \
    1000000usdc # amount
```

### gRPC
A user can query and interact with the `x/settlement` module using gRPC.


#### REST
A user can query and interact with the `x/settlement` module using REST.
