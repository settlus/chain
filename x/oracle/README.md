# Abstract

The `x/oracle` module equips the Settlus blockchain with real-time, accurate block information and NFT ownerships from external blockchains.
This enables the [`x/settlement`](../settlement/README.md) module to utilize NFTs from external chains without using a bridge.
The module has a generic structure designed to support various kinds of information from outside blockchains in the future.

Since this information is external to the Settlus blockchain, it relies on validators to periodically submit votes on the most recent block data.
The protocol tallies these votes at the end of each `VotePeriod` to update the on-chain block data based on the weighted majority.


## Contents

1. **[Concepts](#Concepts)**
   - [Voting Procedure](#voting-procedure)
   - [Topics](#topics)
   - [Reward](#reward)
   - [Slashing](#slashing)
   - [Abstaining from Voting](#abstaining-from-voting)
2. **[State](#State)**
   - [BlockData](#blockdata)
   - [AggregatePrevote](#aggregateprevote)
   - [AggregateVote](#aggregatevote)
   - [FeederDelegation](#feederdelegation)
   - [MissCounter](#misscounter)
3. **[EndBlock](#end-block)**
   - [Tally Exchange Rate Votes](#tally-exchange-rate-votes)
4. **[Messages](#messages)**
   - [MsgPrevote](#msgprevote)
   - [MsgVote](#msgvote)
   - [MsgFeederDelegationConsent](#msgfeederdelegationconsent)
5. **[Parameters](#parameters)**


# Concepts

## Voting Procedure

Settlus provides `RoundInfo`, which includes details about the voting period, voting topics, criteria, and more. Currently, there are two types of topics: Block and Ownership.

During each `VotePeriod`, validators are required to submit two types of messages: `Prevote` and `Vote`. Validators must initially pre-commit to `Data` for every topic with a `Prevote` message. Before the current `VotePeriod` ends, they must reveal their pre-committed `Data` alongside proof of the pre-commitment with a `Vote` message.

This process ensures that validators commit to their choices before seeing other votes, thereby reducing centralization and free-rider risks.

### Prevote and Vote

- A `MsgPrevote` contains the SHA-256 hash of the combined `Data` for each topic.
- `MsgPrevote` can only be submitted before `PrevoteEnd`.

The `Prevote` contains the SHA-256 hash of the block data in the following string format:

```
sha256(<salt><data for topic1><data for topic2>...)
```

- A `MsgVote` includes the salt used to generate the hash for the `Prevote` submitted in the prevote stage.
- `MsgVote` can only be submitted between `PrevoteEnd` and `VoteEnd`.


The `Vote` contains an array of data for each topic:

```
[
    {
        "Topic": "Block",
        "Data": ["data1", "data2", "data3", ...]
    },
    {
        "Topic": "Ownership",
        "Data": ["data1", "data2", "data3", ...]
    },
    ...
]
```


### Vote Tally

At each round, the protocol tallies the votes to calculate the weighted majority at the `EndBlock` stage.

The submitted salt of each vote is used to verify consistency with the `Prevote` submitted by the validator during the prevote stage. If the validator has not submitted a `Prevote`, or if the SHA-256 hash resulting from the salt does not match the hash from the `Prevote`, the vote is dropped.

For each topic, the most voted `Data` that has more than the `VoteThreshold` total voting power is accepted for the following voting period.

### Ballot Rewards

After tallying, the `tally()` function identifies the winning ballots.

### Block Reorganization Handling

Block reorganization, or reorg, occurs when an existing block at a specific height is replaced. This can happen when an external chain forks.

In the event of a reorg, the list of information is updated to reflect the new longest chain following standard voting procedures.

## Topics

### Block

The Block topic is for determining the current block height of an external chain corresponding to the Settlus block. The list of external chains is from the `SupportedChains` of the Settlement Module. Feeders should submit the number and hash of the block which has the smallest timestamp exceeding the `Timestamp` given in `RoundInfo`. The data format of the Block topic is:

```
<chain-id>:<block-number>/<block-hash>
```

Note that the block data itself doesn't currently have any role in the Settlus chain, but it represents the point-of-view criteria of NFT ownership.

### Ownership

The Ownership topic is for determining the current owner of NFTs from external chains. The list of NFTs that need to be verified comes from `UTXR`s of the Settlement Module. Feeders should submit the NFT owner at the block that will be submitted with the Block topic in the same round (the block having the smallest timestamp exceeding the `timestamp` given in `RoundInfo`). The data format of the Ownership topic is:

```
<chain-id>/<contract-addr>/<token-id>:<owner-addr>/<weight>,<owner-addr>/<weight>,...
```

## Reward

Rewards for oracles are pooled from the fees collected from the [x/settlement](../settlement/README.md) module.

## Slashing

The following events are considered a "miss":

- The validator fails to submit a vote for the current `Data` for each topic in every round.
- The validator fails to vote within the weighted majority for any topic.

During every `SlashWindow`, participating validators must maintain a valid vote rate of at least `MinValidPercentPerSlashWindow` (currently set to 50%). If a validator fails to maintain a valid vote rate, the validator is slashed by `SlashFraction` (currently set to 1%). The slashed validator is automatically temporarily "jailed" by the protocol (to protect the funds of delegators), and the operator is expected to fix the discrepancy promptly to resume validator participation.

## Abstaining from Voting

A validator may abstain from voting by submitting empty data in `MsgVote`. Doing so will absolve them of any penalties for missing `VotePeriod`s, but also disqualify them from receiving Oracle rewards for faithful reporting.

# State

## AggregatePrevote

`AggregatePrevote` contains validator's aggregated prevotes.

- `AggregatePrevote`: `valAddress -> AggregatePrevote`

```go
type AggregatePrevote struct {
    Hash        string
    Voter       string    // Voter validator address 
}
```

## AggregateVote

`AggregateVote` contains validator's aggregated votes.

- `AggregateVote`: `valAddress -> AggregateVote`

```go
type AggregateVote struct {
    repeated VoteData VoteData
    Voter    string // voter val address of validator
}
```

## FeederDelegation

An `sdk.AccAddress` (`setl-` account) address of `validator`'s delegated price feeder.

- FeederDelegation: `valAddress -> sdk.AccAddress`

```go
type FeederDelegation struct {
    FeederAddress string
    ValidatorAddress string
}
```

## MissCount

An `int64` representing the number of `VotePeriods` that validator `operator` missed during the current `SlashWindow`.

- MissCount: `valAddress -> uint64`

```go
type MissCount struct {
    ValidatorAddress string
    MissCount uint64
}
```

# End Block

## Tally Block Data Votes

At the end of every block, the `x/oracle` module checks whether it's the last block of the `VotePeriod`.
If it is, it runs the [Voting Procedure](#voting-procedure):

1. Received votes are organized into ballots. Abstained votes, as well as votes by inactive or jailed validators are ignored.

2. Run `VoteProcessor` for each Topic. A `VoteProcessor`
    - Tally up votes and find the weighted majority Data and winners with `TallyVotes()`.
    - Iterate through winners of the ballot and add their weight to their running total.
    - Store the specific actions for each topic. For example, store the updated `BlockData` on the blockchain for the `Block` topic

3. Count up the validators who [missed](#misscount) the Oracle vote and increase the appropriate miss counters.

4. If at the end of a `SlashWindow`, penalize validators who have missed more than the penalty threshold (submitted
   fewer valid votes than `MinValidPercentPerSlashWindow`)

5. Distribute rewards to ballot winners with `k.RewardBallotWinners()`

6. Clear all prevotes (except ones for the next `VotePeriod`) and votes from the store

# Messages

## MsgPrevote

`Hash` is a hex string generated the SHA256 hash (hex string) of a string with the following format:


Note that in the subsequent `MsgVote`, the salt will have to be revealed.
The salt used must be regenerated for each prevote submission.

```go
type MsgPrevote struct {
    Feeder    string
    Validator string
    Hash string
}
```

## MsgVote

The `MsgVote` contains the actual exchange rates vote.
The `Salt` parameter must match the salt used to create the prevote, otherwise the voter cannot be rewarded.
The `BlockData` field contains the block number and block hash of the block for each chain ID in the whitelist.

```go
type MsgVote struct {
    Feeder    string
    Validator string
	BlockData *BlockData
    Salt      string
}
```

## MsgFeederDelegationConsent

Validators may elect to delegate voting rights to another key to prevent the block signing key from being kept online.
To do so, they must submit a `MsgFeederDelegationConsent`, delegating their oracle voting rights to a `Delegate` that
sign `MsgPrevote` and `MsgVote` on behalf of the validator.

The `Validator` field contains the operator address of the validator (prefixed `settlusvaloper1`).
The `FeederAddress` field is the account address (prefixed `settlus1-`) of the delegate account that will be submitting
votes and prevotes on behalf of the `Validator`.

```go
type MsgDelegateFeedConsent struct {
	Validator string
	FeederAddress string
}
```

# Parameters

The oracle module contains the following parameters:

| Key                        | Type    | Example                                                                          |
|----------------------------|---------|----------------------------------------------------------------------------------|
| votePeriod                 | int     | 3                                                                                |
| voteThreshold              | dec     | "0.500000000000000000"                                                           |
| slashFraction              | dec     | "0.001000000000000000"                                                           |
| slashWindow                | int     | 100                                                                              |
| maxMissCountPerSlashWindow | int     | 10                                                                               |

## Validations
- `votePeriod` must be larger than 0.
- `voteThreshold` must be larger than 0.5.
- `slashFraction` must be larger than 0 and less than 0.1.
- `slashWindow` must be larger than 0 and divisible by the `votePeriod`.
- `maxMissCountPerSlashWindow` must be less than the `slashWindow`.
