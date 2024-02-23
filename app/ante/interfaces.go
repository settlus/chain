package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	settlementtypes "github.com/settlus/chain/x/settlement/types"
)

type SettlementKeeper interface {
	GetParams(ctx sdk.Context) (params settlementtypes.Params)
}

type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
}
