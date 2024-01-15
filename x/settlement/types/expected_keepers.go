package types

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	erc20types "github.com/settlus/chain/evmos/x/erc20/types"
	evmtypes "github.com/settlus/chain/evmos/x/evm/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	NewAccount(ctx sdk.Context, acc authtypes.AccountI) authtypes.AccountI
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)

	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type Erc20Keeper interface {
	ConvertCoin(goCtx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
	IsDenomRegistered(ctx sdk.Context, denom string) bool
}

type EvmKeeper interface {
	CallEVM(ctx sdk.Context, abi abi.ABI, from, contract common.Address, commit bool, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
	CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
}

// NftOwnershipKeeper defines the expected interface needed to retrieve owner of NFT.
type NftOwnershipKeeper interface {
	CheckValidChainId(ctx sdk.Context, chainId string) bool
	OwnerOf(ctx sdk.Context, chainId string, contractAddr string, tokenIdHex string) (*common.Address, error)
}
