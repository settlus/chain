#!/bin/bash

ETH_RPC_URL="$1"

if [ -z "$ETH_RPC_URL" ]; then
    echo "Error: ETH_RPC_URL is required as an argument."
    echo "Usage: $0 <ETH_RPC_URL>"
    exit 1
fi

CHAINID="${CHAIN_ID:-settlus_5371-1}"
MONIKER="test"
KEYRING="test"
KEYALGO="eth_secp256k1" #gitleaks:allow
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""
PRUNING="default"
#PRUNING="custom"

CHAINDIR="$HOME/.settlus"
GENESIS="$CHAINDIR/config/genesis.json"
TMP_GENESIS="$CHAINDIR/config/tmp_genesis.json"
APP_TOML="$CHAINDIR/config/app.toml"
CONFIG_TOML="$CHAINDIR/config/config.toml"

# feemarket params basefee
BASEFEE=1000000000

# treasury address 0x7cb61d4117ae31a12e393a1cfa3bac666481d02e
VAL_KEY="treasury"
VAL_MNEMONIC="equal broken goose strong twenty upgrade cool pen run opinion gain brick husband repeat magnet foam creek purse alcohol this margin lunch hip birth"

# bob from config.yml
BOB_KEY="bob"
BOB_MNEMONIC="police tube stay federal expire veteran roof fossil simple purse ridge knee wheel topple omit review spider public tone prosper side imitate auto inhale"

# faucet from config.yml
FAUCET_KEY="faucet"
FAUCET_MNEMONIC="island club point history solution tonight festival maid zebra business nasty clap spirit science excess win caution hand embrace heavy snow derive nuclear head"

# validate dependencies are installed
command -v jq >/dev/null 2>&1 || {
	echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"
	exit 1
}

# used to exit on first error (any non-zero exit code)
set -e

# Check settlusd version to decide how to set the client configuration
sdk_version=$(settlusd version --long | grep 'cosmos_sdk_version' | awk '{print $2}')
if [[ $sdk_version == *v0.4* ]]; then
	settlusd config chain-id "$CHAINID"
	settlusd config keyring-backend "$KEYRING"
else
	settlusd config set client chain-id "$CHAINID"
	settlusd config set client keyring-backend "$KEYRING"
fi

# Import keys from mnemonics
echo "$VAL_MNEMONIC" | settlusd keys add "$VAL_KEY" --recover --keyring-backend "$KEYRING" --algo "$KEYALGO"
echo "$BOB_MNEMONIC" | settlusd keys add "$BOB_KEY" --recover --keyring-backend "$KEYRING" --algo "$KEYALGO"
echo "$FAUCET_MNEMONIC" | settlusd keys add "$FAUCET_KEY" --recover --keyring-backend "$KEYRING" --algo "$KEYALGO"

# Set moniker and chain-id for Evmos (Moniker can be anything, chain-id must be an integer)
settlusd init "$MONIKER" --chain-id "$CHAINID" --overwrite --home "$CHAINDIR"

# Change parameter token denominations to asetl
jq '.app_state["staking"]["params"]["bond_denom"]="asetl"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="asetl"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
# When upgrade to cosmos-sdk v0.47, use gov.params to edit the deposit params
jq '.app_state["gov"]["params"]["min_deposit"][0]["denom"]="asetl"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["evm"]["params"]["evm_denom"]="asetl"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["inflation"]["params"]["mint_denom"]="asetl"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# set gov proposing && voting period
jq '.app_state.gov.deposit_params.max_deposit_period="10s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state.gov.voting_params.voting_period="10s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
sed -i.bak 's/"expedited_voting_period": "86400s"/"expedited_voting_period": "5s"/g' "$GENESIS"

jq '.app_state.gov.params.min_deposit[0].denom="asetl"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state.gov.params.max_deposit_period="10s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state.gov.params.voting_period="10s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# Set gas limit in genesis
jq '.consensus_params.block.max_gas="10000000000000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# Set base fee in genesis
jq '.app_state["feemarket"]["params"]["base_fee"]="'${BASEFEE}'"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# Set settlement params
jq '.app_state.settlement.params.gas_prices = [{"amount": "1", "denom": "uusdc"}, {"amount": "0.0001", "denom": "setl"}] | 
    .app_state.settlement.params.oracle_fee_percentage = "1" | 
    .app_state.settlement.params.supported_chains = [{"chain_id": "1", "chain_name": "Ethereum", "chain_url": "https://ethereum.org"}] | 
    .app_state.settlement.tenants = [{"admins": ["settlus12g8w5dr5jyncct8jwdxwsy2g9ktdrjjlcs5f0a"], "denom": "eBLUC", "id": "0", "payout_method": "mintable_contract", "payout_period": 10}] | 
    .app_state.settlement.utxrs = []' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# Set oracle params
jq '.app_state.oracle.params = {
    "vote_period": 2,
    "vote_threshold": "0.5",
    "slash_fraction": "0.01",
    "slash_window": 604800,
    "max_miss_count_per_slash_window": 302400
}' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# disable produce empty block
sed -i.bak 's/create_empty_blocks = true/create_empty_blocks = false/g' "$CONFIG_TOML"

# Allocate genesis accounts (cosmos formatted addresses)
settlusd add-genesis-account "$(settlusd keys show "$VAL_KEY" -a --keyring-backend "$KEYRING")" 10000000000000000000000000asetl,10000000000000000000000000uusdc --keyring-backend "$KEYRING"
settlusd add-genesis-account "$(settlusd keys show "$BOB_KEY" -a --keyring-backend "$KEYRING")" 100000000000000000000000asetl,10000000000000000000000000uusdc --keyring-backend "$KEYRING"
settlusd add-genesis-account "$(settlusd keys show "$FAUCET_KEY" -a --keyring-backend "$KEYRING")" 100000000000000000000000asetl,10000000000000000000000000uusdc --keyring-backend "$KEYRING"

# Update total supply with claim values
# Bc is required to add this big numbers
# total_supply=$(bc <<< "$validators_supply")
total_supply=10200000000000000000000000
jq -r --arg total_supply "$total_supply" '.app_state.bank.supply[0].amount=$total_supply' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# set custom pruning settings
if [ "$PRUNING" = "custom" ]; then
	sed -i.bak 's/pruning = "default"/pruning = "custom"/g' "$APP_TOML"
	sed -i.bak 's/pruning-keep-recent = "0"/pruning-keep-recent = "2"/g' "$APP_TOML"
	sed -i.bak 's/pruning-interval = "0"/pruning-interval = "10"/g' "$APP_TOML"
fi

# make sure the localhost IP is 0.0.0.0
sed -i.bak 's/localhost/0.0.0.0/g' "$CONFIG_TOML"
sed -i.bak 's/127.0.0.1/0.0.0.0/g' "$APP_TOML"

# use timeout_commit 1s to make test faster
sed -i.bak 's/timeout_commit = "3s"/timeout_commit = "1s"/g' "$CONFIG_TOML"

# Sign genesis transaction
settlusd gentx "$VAL_KEY" 1000000000000000000000asetl --gas-prices ${BASEFEE}asetl --keyring-backend "$KEYRING" --chain-id "$CHAINID"
## In case you want to create multiple validators at genesis
## 1. Back to `settlusd keys add` step, init more keys
## 2. Back to `settlusd add-genesis-account` step, add balance for those
## 3. Clone this ~/.settlusd home directory into some others, let's say `~/.clonedsettlusd`
## 4. Run `gentx` in each of those folders
## 5. Copy the `gentx-*` folders under `~/.clonedsettlusd/config/gentx/` folders into the original `~/.settlusd/config/gentx`

# Enable the APIs for the tests to be successful
sed -i.bak 's/enable = false/enable = true/g' "$APP_TOML"
# Don't enable Rosetta API by default
grep -q -F '[rosetta]' "$APP_TOML" && sed -i.bak '/\[rosetta\]/,/^\[/ s/enable = true/enable = false/' "$APP_TOML"
# Don't enable memiavl by default
grep -q -F '[memiavl]' "$APP_TOML" && sed -i.bak '/\[memiavl\]/,/^\[/ s/enable = true/enable = false/' "$APP_TOML"
# Don't enable versionDB by default
grep -q -F '[versiondb]' "$APP_TOML" && sed -i.bak '/\[versiondb\]/,/^\[/ s/enable = true/enable = false/' "$APP_TOML"

echo "collect-gentxs"

# Collect genesis tx
settlusd collect-gentxs

# echo "validate-genesis"

# Run this to ensure everything worked and that the genesis file is setup correctly
# settlusd validate-genesis

echo "Genesis file validated"

# Start the node in the background
settlusd start "$TRACE" \
	--log_level $LOGLEVEL \
	--minimum-gas-prices=0.0001asetl \
	--json-rpc.api eth,txpool,personal,net,debug,web3 \
	--chain-id "$CHAINID" > /dev/null 2>&1 &

# Wait for the node to start
sleep 10

# Optional: You can add a message to indicate that the node has started
echo "Node started in the background"

# Create directory for interop config
mkdir -p /root/.interop

# Save the config file
cat << EOF > /root/.interop/config.yaml
settlus:
  chain_id: settlus_5371-1
  rpc_url: http://localhost:26657
  grpc_url: http://localhost:9090
  insecure: true
  gas_limit: 200000
  fees:
    denom: asetl
    amount: "210000000000000"
feeder:
  topics: block
  address: settlus12g8w5dr5jyncct8jwdxwsy2g9ktdrjjlcs5f0a
  signer_mode: local
  key: 8be29a465f945630ca905af7a6977a5b2bfa735fb7996d44b630420af8fc9ed4
  validator_address: settlusvaloper12g8w5dr5jyncct8jwdxwsy2g9ktdrjjluy76df
chains:
- chain_id: "1"
  chain_name: Ethereum
  chain_type: ethereum
  rpc_url: $ETH_RPC_URL
log_level: debug
EOF

echo "Interop config file created at /root/.interop/config.yaml"

# Start the interop node
interop-node start > /dev/null 2>&1 &

# Wait for the interop node to start
sleep 10

# Optional: You can add a message to indicate that the interop node has started
echo "Interop node started in the background"

# Run the test
go test ./tests/e2e -v || (echo "Tests failed"; exit 1)
