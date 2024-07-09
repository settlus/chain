package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/settlement/types"
)

type SettlementKeeper interface {
	GetParams(ctx sdk.Context) (params types.Params)
}

type OracleKeeper interface {
	ValidateFeeder(ctx sdk.Context, feederAddr string, validatorAddr string) (bool, error) 
}