package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// StakingKeeper expected staking keeper (noalias)
type StakingKeeper interface {
	MinCommissionRate(ctx sdk.Context) sdk.Dec
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) (validator stakingtypes.Validator, found bool)
	BondDenom(ctx sdk.Context) string
	SetValidator(ctx sdk.Context, validator stakingtypes.Validator)
	SetValidatorByConsAddr(ctx sdk.Context, validator stakingtypes.Validator) error
	SetNewValidatorByPowerIndex(ctx sdk.Context, validator stakingtypes.Validator)
	AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error)

	GetParams(ctx sdk.Context) stakingtypes.Params
}
