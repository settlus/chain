package keeper_test

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/mock"

	"github.com/settlus/chain/evmos/x/evm/statedb"
	evm "github.com/settlus/chain/evmos/x/evm/types"
	"github.com/settlus/chain/x/nftownership/types"
)

var _ types.EVMKeeper = &MockEVMKeeper{}

type MockEVMKeeper struct {
	mock.Mock
}

func (m *MockEVMKeeper) GetParams(_ sdk.Context) evm.Params {
	args := m.Called(mock.Anything)
	return args.Get(0).(evm.Params)
}

func (m *MockEVMKeeper) GetAccountWithoutBalance(_ sdk.Context, _ common.Address) *statedb.Account {
	args := m.Called(mock.Anything, mock.Anything)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*statedb.Account)
}

func (m *MockEVMKeeper) EstimateGas(_ context.Context, _ *evm.EthCallRequest) (*evm.EstimateGasResponse, error) {
	args := m.Called(mock.Anything, mock.Anything)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*evm.EstimateGasResponse), args.Error(1)
}

func (m *MockEVMKeeper) ApplyMessage(_ sdk.Context, _ core.Message, _ vm.EVMLogger, _ bool) (*evm.MsgEthereumTxResponse, error) {
	args := m.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*evm.MsgEthereumTxResponse), args.Error(1)
}

// CallEVM implements types.EVMKeeper.
func (m *MockEVMKeeper) CallEVM(ctx sdk.Context, abi abi.ABI, from common.Address, contract common.Address, commit bool, method string, arg ...interface{}) (*evm.MsgEthereumTxResponse, error) {
	args := m.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*evm.MsgEthereumTxResponse), args.Error(1)
}

// CallEVMWithData implements types.EVMKeeper.
func (m *MockEVMKeeper) CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool) (*evm.MsgEthereumTxResponse, error) {
	args := m.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*evm.MsgEthereumTxResponse), args.Error(1)
}
