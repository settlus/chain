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

func (k SettlementKeeper) GetSupportedChainIds(ctx sdk.Context) []string {
	params := k.GetParams(ctx)
	chainIds := make([]string, 0)
	for _, chain := range params.GetSupportedChains() {
		chainIds = append(chainIds, chain.ChainId)
	}
	return chainIds
}

func (k SettlementKeeper) IsSupportedChain(ctx sdk.Context, chainId string) bool {
	params := k.GetParams(ctx)
	for _, chain := range params.GetSupportedChains() {
		if chain.ChainId == chainId {
			return true
		}
	}
	return false
}
