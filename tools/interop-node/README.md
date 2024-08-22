# InterOp Node
InterOp Node is an application that facilitates communication between Settlus and external chains, supporting multi-chain NFT interoperability.

## Concepts
There are mainly two components in the InterOp Node: Subscriber and Feeder.
The Subscriber is responsible for fetching information from external chains, and the Feeder is responsible for sending votes to Settlus.

### Subscriber
Each subscriber retrieves block data and NFT ownership information from external blockchain networks.
External blockchain's block information (hash and number) is stored in the cache.
NFT owner information is always fetched from the external chain.

Supported chains:
- [Ethereum](./subscriber/ethereum_subscriber.go)

#### Cache
For each block, the subscriber stores the block hash and number in a memory cache.
Each subscriber runs a separate goroutine to periodically check the latest block number and hash from the external chain and update the cache.

When initializing the cache, users can specify the cache size to limit memory usage.

Cache implementation: [cache.go](./subscriber/cache.go)

### Feeder

The Feeder sends votes to Settlus. There can be multiple topics for voting, but currently, there is only one topic: `NFT Ownership`.
1. The feeder fetches `RoundInfo` from the Settlus blockchain.
2. From the `RoundInfo`, the feeder checks if Settlus is in the `VOTING` period.
  - If Settlus is in the `VOTING` period, the feeder submits the `Vote` TX to Settlus.
  - If Settlus is not in the `VOTING` period, the feeder gathers block data from the chains and submits the `Prevote` TX to Settlus.
  - If the feeder cannot send the vote appropriately, it submits an Abstain vote to Settlus.

## How to Run
Initialize the config file.
The config file will be located at ~/.interop/
```shell
./interop-node config init
```

Edit the config file with appropriate values and run the oracle feeder.
```shell
./interop-node start
```

## Local Development
First, run local Settlus.
```shell
ignite chain serve
```

Create a feeder account.
```shell
settlusd keys add foo
- address: settlus1nyw0ruj3t5wdh9ycgcsxles6mpfz9xmk93m9cy
  name: foo
  pubkey: '{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A8KPbBSg2xj/OFCkgcN0doTTGD4MikLFmfjQy4CjQ/lw"}'
  type: local
```

Fund the feeder account.
```shell
settlusd tx bank send YOUR_KEY_WITH_FUND settlus1nyw0ruj3t5wdh9ycgcsxles6mpfz9xmk93m9cy 10setl --fees 0.001setl --keyring-backend test
```

Get the private key of the feeder account.
```shell
settlusd keys export foo --unarmored-hex --unsafe --keyring-backend test

WARNING: The private key will be exported as an unarmored hexadecimal string. USE AT YOUR OWN RISK. Continue? [y/N]: y
1b8dac2949968eae623859330172283d27dda8b496f1dae2cbed9b1bcce51cb1
```

Get the validator address.
```shell
settlusd q staking validators

...
  min_self_delegation: "1"
  operator_address: settlusvaloper12g8w5dr5jyncct8jwdxwsy2g9ktdrjjluy76df
  status: BOND_STATUS_BONDED
...
```

Set oracle feeder delegation.
```shell
settlusd tx oracle feeder-delegation-consent settlus1nyw0ruj3t5wdh9ycgcsxles6mpfz9xmk93m9cy --from YOUR_VALIDATOR --fees 0.01setl --keyring-backend test
```

Check if the oracle feeder delegation is set.
```shell
settlusd q oracle feeder-delegation settlusvaloper12g8w5dr5jyncct8jwdxwsy2g9ktdrjjluy76df

feeder_delegation:
  feeder_address: settlus1nyw0ruj3t5wdh9ycgcsxles6mpfz9xmk93m9cy
  validator_address: settlusvaloper12g8w5dr5jyncct8jwdxwsy2g9ktdrjjluy76df
```

Update config.yaml
```yaml
settlus:
  chain_id: "settlus_5371-1"
  rpc_url: http://localhost:26657
  grpc_url: http://localhost:9090
  insecure: true
  avg_block_time: 5s
  voting_period: 10
  gas_limit: 200000
  fees:
    denom: asetl
    amount: "10000000000000"
feeder:
  address: settlus1nyw0ruj3t5wdh9ycgcsxles6mpfz9xmk93m9cy
  private_key: 1b8dac2949968eae623859330172283d27dda8b496f1dae2cbed9b1bcce51cb1 # <- example private key
  validator_address: settlusvaloper12g8w5dr5jyncct8jwdxwsy2g9ktdrjjluy76df
chains:
- chain_id: "1"
  chain_name: ethereum
  chain_type: ethereum
  rpc_url: https://mainnet.infura.io/v3/YOUR_INFURA_KEY
log_level: info
```

Run the oracle feeder.
```shell
./interop-node start
```
