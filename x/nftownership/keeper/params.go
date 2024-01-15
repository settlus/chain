package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/nftownership/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)

	if params.AllowedChainIds == nil {
		return types.DefaultParams()
	}

	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
