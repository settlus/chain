package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/settlus/chain/evmos/x/evm/statedb"
	evmtypes "github.com/settlus/chain/evmos/x/evm/types"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// EVMKeeper defines the expected EVM keeper interface used on erc20
type EVMKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
	GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account
	EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
	CallEVM(ctx sdk.Context, abi abi.ABI, from, contract common.Address, commit bool, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
	CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
}

type OracleKeeper interface {
	GetBlockData(ctx sdk.Context, chainId string) (*oracletypes.BlockData, error)
}
