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
    * [Settlement](#settlement)
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
Each UTXR contains details such as the recipient's address, the amount of the transaction, and the end of the payout period
Stored in the `unspent_records` state, these records are the backbone of the settlement process.

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

> `payout_period` is set to 100800.
> Assume average block period is 6 seconds.
> Then a `payout_period` of 100800 (60 * 60 * 24 * 7 / 6) which is about a week, means that the UTXR can be canceled for 100800 blocks after the UTXR is created.
> After 100800 blocks, the UTXR is considered settled and cannot be canceled.

### Settlement
Once the payout period of a UTXR has concluded, the UTXR is eligible for settlement.
This process is the final step in the lifecycle of a transaction and is essential for the actual transfer of funds from tenants to recipients.

At each `BeginBlock`, the `x/settlement` module iterates through the `unspent_records` state and checks if the UTXR's payout period has passed.
The UTXRs that have passed the payout period are considered eligible for settlement.
If the tenant has enough funds in the treasury to settle the UTXR, the UTXR is removed from the `unspent_records` state and the amount is transferred from the tenant's treasury to the recipient's wallet.
If the tenant does not have enough funds in the treasury to settle the UTXR, the settlement will be deferred until the tenant has enough funds.
The UTXR will remain in the `unspent_records` state until the tenant has enough funds to settle the UTXR.

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
  RequestId string
  Recipient sdk.AccAddress
  Amount sdk.Coins
  PayoutBlock uint64
}
```
The UTXR ID is incremented by 1 for each UTXR. Because the UTXRs are created in order, the UTXRs are trivially sorted by `PayoutBlock`.

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
- A tenant creates a UTXR with payout period of `100800` blocks.
- A new UTXR is created with the following values.
    - `Recipient`: `RecipientAddress`
    - `Amount`: `Amount`
    - `PayoutBlock`: `CurrentBlockHeight + 100800 (PayoutPeriod)`
- Add the new UTXR to the store with the following key: `((TenantID)-(UTXRID))`.

**Settlement**
- Iterate from the lowest key to the highest key in the `unspent_records` store.
- For each UTXR, do the following:
  - Check if the current UTXR's `PayoutBlock` is less than or equal to the current block height.  
    - If no, since every UTXR is sorted by `PayoutBlock`, we can stop checking the rest of the UTXRs.
    - If yes, check if the tenant has enough funds to settle the UTXR.
      - If the tenant has enough funds, remove the UTXR from the `unspent_records` state and transfer the amount from the tenant's treasury to the recipient's wallet.
      - If the tenant does not have enough funds, stop the iteration and emit a `NotEnoughTreasuryBalance` event.

**Cancel**
In the case of a cancel, the UTXR is simply removed from the `unspent_records` state.
- Get the UTXR ID by Request ID.
- Delete the UTXR from the `unspent_records` state by the UTXR ID.

## Begin Block
At each `BeginBlock`, the `x/settlement` module iterates through the `unspent_records` state and checks if the UTXR's payout period has passed.

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

### MstCancel
The `MstCancel` message allows tenant admins to cancel a UTXR.
```go
type MstCancel struct {
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
	ChainId string
    Recipient string
	NftAddress string
	NftTokenId string
	Ammount sdk.Coins
	Metdata string
}
```

### EventCancel
```go
type EventCancel struct {
    TenantId string
	RequestId string
}
```

## Parameters
The `x/settlement` module contains the following parameters:

| Key                        | Type     | Example                |
|----------------------------|----------|------------------------|
| fee                        | Coin     | {denom: usdc}          |
| oracle_fee_percentage      | dec      | "0.500000000000000000" |
| tenants                    | []Tenant | {}                     |

```go
type Tenant struct {
    Id uint64
	Name string
    Admins []sdk.AccAddress
    PayoutPeriod uint64
}
```

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
settlusd query settlement utxrs [tenant-id] [id] [flags]
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
  "payout_block": "100800",
  "recipient": "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
  "amount": [
      {
      "denom": "usdc",
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
settlusd tx settlement record-revenue \
    1 # tenant id \
    request-1 # request id \
    1000000usdc # amount
    ethereum # chain id \
    0x0000000000000000000000000000000000000001 # contract address \
    1 # token id \
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
