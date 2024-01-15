# InterOp Node
InterOp Node is an application that facilitates communication between Settlus and external chains, supporting multi-chain NFT interoperability.

## Concepts

### Subscriber
Each subscriber retrieves block data and NFT ownership information from external blockchain networks and stores it locally in a database.

#### DB Schema
For each block, the Subscriber stores three types of key-value pairs in a single transaction:

| Type          | Prefix |                                                                Key |        Value |
|---------------|-------|-------------------------------------------------------------------|-------------|
| Block Hash    |   BH   |                                               BLOCK_HASH (32bytes) | BLOCK_NUMBER |
| Block Number  |   BN   |                                             BLOCK_NUMBER (32bytes) |   BLOCK_HASH |
| NFT Ownership |   NO   | NFT_ADDR (20bytes) \+ TOKEN_ID (32bytes) \+ BLOCK_NUMBER (32bytes) |  OWNER_ADDR  |

To find the owner of an NFT at block number A, a reverse iterator is utilized with the key (NO + NFT_ADDR + TOKEN_ID + A). To identify the owner at a specific block hash X, the process begins by searching with the key (BH + BLOCK_HASH) to retrieve the block number, then proceeds in the same manner.

#### Fallback
If the database lacks the ownership data, the system can directly query the blockchain network. This situation might arise when the node is new and hasn't yet processed all the block data.

### Feeder

Feeder sends vote to Settlus. There can be multiple topics for voting, but currently we only support topic to decide block number.
- The feeder will check if Settlus is in the `VOTING` period.
   - If Settlus is in the `VOTING` period, feeder submits the `Vote` TX to Settlus.
   - If Settlus is not in the `VOTING` period, feeder gathers block data from the chains and submits the `Prevote` TX to Settlus.
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
First run local Settlus.
```shell
ignite chain serve
```

Create feeder account.
```shell
chaind keys add foo
- address: settlus1nyw0ruj3t5wdh9ycgcsxles6mpfz9xmk93m9cy
  name: foo
  pubkey: '{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A8KPbBSg2xj/OFCkgcN0doTTGD4MikLFmfjQy4CjQ/lw"}'
  type: local
```

Fund the feeder account.
```shell
chaind tx bank send YOUR_KEY_WITH_FUND settlus1nyw0ruj3t5wdh9ycgcsxles6mpfz9xmk93m9cy 10setl --fees 0.001setl --keyring-backend test
```

Get private key of the feeder account
```shell
chaind keys export foo --unarmored-hex --unsafe --keyring-backend test

WARNING: The private key will be exported as an unarmored hexadecimal string. USE AT YOUR OWN RISK. Continue? [y/N]: y
1b8dac2949968eae623859330172283d27dda8b496f1dae2cbed9b1bcce51cb1
```

Get the validator address.
```shell
chaind q staking validators

...
  min_self_delegation: "1"
  operator_address: settlusvaloper12g8w5dr5jyncct8jwdxwsy2g9ktdrjjluy76df
  status: BOND_STATUS_BONDED
...
```

Set oracle feeder delegation.
```shell
chaind tx oracle feeder-delegation-consent settlus1nyw0ruj3t5wdh9ycgcsxles6mpfz9xmk93m9cy --from YOUR_VALIDATOR --fees 0.01setl --keyring-backend test
```

Check if the oracle feeder delegation is set.
```shell
chaind q oracle feeder-delegation settlusvaloper12g8w5dr5jyncct8jwdxwsy2g9ktdrjjluy76df

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
  private_key: 1b8dac2949968eae623859330172283d27dda8b496f1dae2cbed9b1bcce51cb1
  validator_address: settlusvaloper12g8w5dr5jyncct8jwdxwsy2g9ktdrjjluy76df
chains:
- chain_id: "1"
  chain_name: ethereum
  chain_type: ethereum
  rpc_url: https://mainnet.infura.io/v3/YOUR_INFURA_KEY
db_home: ...
log_level: info
```

Run the oracle feeder.
```shell
./interop-node start
```