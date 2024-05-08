package keeper

import (
	"github.com/settlus/chain/x/settlement/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k SettlementKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	if params.GasPrices == nil {
		return types.DefaultParams()
	}

	return params
}

// SetParams set the params
func (k SettlementKeeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
