package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	// GetValidator get a single validator
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, bool)
	// MaxValidators returns the maximum amount of bonded validators
	MaxValidators(sdk.Context) uint32
	// ValidatorsPowerStoreIterator  returns an iterator for the current validator power store
	ValidatorsPowerStoreIterator(ctx sdk.Context) sdk.Iterator
	// PowerReduction - is the amount of staking tokens required for 1 unit of consensus-engine power.
	// Currently, this returns a global variable that the app developer can tweak.
	PowerReduction(ctx sdk.Context) math.Int
	// TotalBondedTokens total staking tokens supply which is bonded.
	TotalBondedTokens(ctx sdk.Context) math.Int
	// Slash a validator for an infraction committed at a known height
	// Find the contributing stake at that height and burn the specified slashFactor
	// of it, updating unbonding delegations & redelegations appropriately
	//
	// CONTRACT:
	//
	//	slashFactor is non-negative
	//
	// CONTRACT:
	//
	//	Infraction was committed equal to or less than an unbonding period in the past,
	//	so all unbonding delegations and redelegations from that height are stored
	//
	// CONTRACT:
	//
	//	Slash will not slash unbonded validators (for the above reason)
	//
	// CONTRACT:
	//
	//	Infraction was committed at the current height or at a past height,
	//	not at a height in the future
	Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) math.Int
	// Jail a validator
	Jail(ctx sdk.Context, consAddr sdk.ConsAddress)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// DistributionKeeper defines the expected interface needed to manage validator reward distribution.
type DistributionKeeper interface {
	// AllocateTokensToValidator allocate tokens to a particular validator,
	// splitting according to commission.
	AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins)
}
