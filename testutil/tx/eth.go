package tx

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	evmosserverconfig "github.com/evmos/evmos/v19/server/config"
	evmtypes "github.com/evmos/evmos/v19/x/evm/types"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/cmd/settlusd/config"
)

// PrepareEthTx creates an ethereum tx and signs it with the provided messages and private key.
// It returns the signed transaction and an error
func PrepareEthTx(
	txCfg client.TxConfig,
	appSettlus *app.SettlusApp,
	priv cryptotypes.PrivKey,
	msgs ...sdk.Msg,
) (authsigning.Tx, error) {
	txBuilder := txCfg.NewTxBuilder()

	signer := ethtypes.LatestSignerForChainID(appSettlus.EvmKeeper.ChainID())
	txFee := sdk.Coins{}
	txGasLimit := uint64(0)

	// Sign messages and compute gas/fees.
	for _, m := range msgs {
		msg, ok := m.(*evmtypes.MsgEthereumTx)
		if !ok {
			return nil, errorsmod.Wrapf(errorsmod.Error{}, "cannot mix Ethereum and Cosmos messages in one Tx")
		}

		if priv != nil {
			err := msg.Sign(signer, NewSigner(priv))
			if err != nil {
				return nil, err
			}
		}

		msg.From = ""

		txGasLimit += msg.GetGas()
		txFee = txFee.Add(sdk.Coin{Denom: config.BaseDenom, Amount: sdkmath.NewIntFromBigInt(msg.GetFee())})
	}

	if err := txBuilder.SetMsgs(msgs...); err != nil {
		return nil, err
	}

	// Set the extension
	var option *codectypes.Any
	option, err := codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
	if err != nil {
		return nil, err
	}

	builder, ok := txBuilder.(authtx.ExtensionOptionsTxBuilder)
	if !ok {
		return nil, errorsmod.Wrapf(errorsmod.Error{}, "could not set extensions for Ethereum tx")
	}

	builder.SetExtensionOptions(option)

	txBuilder.SetGasLimit(txGasLimit)
	txBuilder.SetFeeAmount(txFee)

	return txBuilder.GetTx(), nil
}

// GasLimit estimates the gas limit for the provided parameters. To achieve
// this, need to provide the corresponding QueryClient to call the
// `eth_estimateGas` rpc method. If not provided, returns a default value
func GasLimit(ctx sdk.Context, from common.Address, data evmtypes.HexString, queryClientEvm evmtypes.QueryClient) (uint64, error) {
	// default gas limit (used if no queryClientEvm is provided)
	gas := uint64(100000000000)

	if queryClientEvm != nil {
		args, err := json.Marshal(&evmtypes.TransactionArgs{
			From: &from,
			Data: (*hexutil.Bytes)(&data),
		})
		if err != nil {
			return gas, err
		}

		goCtx := sdk.WrapSDKContext(ctx)
		res, err := queryClientEvm.EstimateGas(goCtx, &evmtypes.EthCallRequest{
			Args:   args,
			GasCap: evmosserverconfig.DefaultGasCap,
		})
		if err != nil {
			return gas, err
		}
		gas = res.Gas
	}
	return gas, nil
}
