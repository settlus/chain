package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/oracle/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)

	if params.VotePeriod == 0 {
		return types.DefaultParams()
	}

	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	k.paramstore.SetParamSet(ctx, &params)

	return nil
}

// GetChains returns the list of whitelist chains in params
func (k Keeper) GetChains(ctx sdk.Context) []*types.Chain {
	params := k.GetParams(ctx)
	wl := params.GetWhitelist()
	if wl == nil {
		return nil
	}
	return wl
}

// GetChain returns the chain with the given chain id
func (k Keeper) GetChain(ctx sdk.Context, chainId string) (*types.Chain, error) {
	params := k.GetParams(ctx)
	wl := params.GetWhitelist()
	if len(wl) == 0 {
		return nil, fmt.Errorf("whitelist is not set")
	}
	for _, chain := range wl {
		if chain.ChainId == chainId {
			return chain, nil
		}
	}

	return nil, types.ErrChainNotFound
}
