package testutil

import (
	"fmt"
	"math/big"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	evm "github.com/evmos/evmos/v19/x/evm/types"
	"github.com/gogo/protobuf/proto"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/testutil/tx"
)

// DeployContract deploys a contract with the provided private key,
// compiled contract data and constructor arguments
func DeployContract(
	ctx sdk.Context,
	app *app.SettlusApp,
	priv cryptotypes.PrivKey,
	queryClientEvm evm.QueryClient,
	contract evm.CompiledContract,
	constructorArgs ...interface{},
) (common.Address, error) {
	chainID := app.EvmKeeper.ChainID()
	from := common.BytesToAddress(priv.PubKey().Address().Bytes())
	nonce := app.EvmKeeper.GetNonce(ctx, from)

	ctorArgs, err := contract.ABI.Pack("", constructorArgs...)
	if err != nil {
		return common.Address{}, err
	}

	data := append(contract.Bin, ctorArgs...) //nolint:gocritic
	gas, err := tx.GasLimit(ctx, from, data, queryClientEvm)
	if err != nil {
		return common.Address{}, err
	}

	msgEthereumTx := evm.NewTx(&evm.EvmTxArgs{
		ChainID:   chainID,
		Nonce:     nonce,
		GasLimit:  gas,
		GasFeeCap: app.FeeMarketKeeper.GetBaseFee(ctx),
		GasTipCap: big.NewInt(1),
		Input:     data,
		Accesses:  &ethtypes.AccessList{},
	})
	msgEthereumTx.From = from.String()

	res, err := DeliverEthTx(app, priv, msgEthereumTx)
	if err != nil {
		return common.Address{}, err
	}

	if _, err := CheckEthTxResponse(res, app.AppCodec()); err != nil {
		return common.Address{}, err
	}

	return crypto.CreateAddress(from, nonce), nil
}

func MintNFT(
	ctx sdk.Context,
	appSettlus *app.SettlusApp,
	priv cryptotypes.PrivKey,
	contract evm.CompiledContract,
	contractAddress common.Address,
	constructorArgs ...interface{},
) error {
	from := common.BytesToAddress(priv.PubKey().Address().Bytes())

	fmt.Printf("constructorArgs %v\n", constructorArgs)

	_, err := appSettlus.EvmKeeper.CallEVM(ctx, contract.ABI, from, contractAddress, true, "safeMint", constructorArgs...)
	if err != nil {
		return fmt.Errorf("call evm failed: %w", err)
	}

	return nil
}

func CheckNFTExists(
	ctx sdk.Context,
	appSettlus *app.SettlusApp,
	priv cryptotypes.PrivKey,
	contract evm.CompiledContract,
	contractAddress common.Address,
	constructorArgs ...interface{},
) (bool, error) {
	from := common.BytesToAddress(priv.PubKey().Address().Bytes())

	res, err := appSettlus.EvmKeeper.CallEVM(ctx, contract.ABI, from, contractAddress, true, "exists", constructorArgs...)
	if err != nil {
		return false, fmt.Errorf("call evm failed: %w", err)
	}

	exists := new(big.Int).SetBytes(res.Ret).String()

	return exists == "1", nil
}

// CheckEthTxResponse checks that the transaction was executed successfully
func CheckEthTxResponse(r abci.ResponseDeliverTx, cdc codec.Codec) (*evm.MsgEthereumTxResponse, error) {
	if !r.IsOK() {
		return nil, fmt.Errorf("tx failed. Code: %d, Logs: %s", r.Code, r.Log)
	}
	var txData sdk.TxMsgData
	if err := cdc.Unmarshal(r.Data, &txData); err != nil {
		return nil, err
	}

	var res evm.MsgEthereumTxResponse
	if err := proto.Unmarshal(txData.MsgResponses[0].Value, &res); err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, fmt.Errorf("tx failed. VmError: %s", res.VmError)
	}

	return &res, nil
}
