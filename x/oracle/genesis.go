package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/oracle/keeper"
	"github.com/settlus/chain/x/oracle/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(fmt.Errorf("failed to set params (%s)", err))
	}
	for _, aggregateVote := range genState.AggregateVotes {
		k.SetAggregateVote(ctx, aggregateVote)
	}
	for _, aggregatePrevote := range genState.AggregatePrevotes {
		k.SetAggregatePrevote(ctx, aggregatePrevote)
	}
	for _, missCount := range genState.MissCounts {
		k.SetMissCount(ctx, missCount.ValidatorAddress, missCount.MissCount)
	}
	for _, blockData := range genState.BlockData {
		k.SetBlockData(ctx, blockData)
	}
	for _, feederDelegation := range genState.FeederDelegation {
		if err := k.SetFeederDelegation(ctx, feederDelegation.ValidatorAddress, feederDelegation.FeederAddress); err != nil {
			panic(fmt.Errorf("failed to set feeder delegation (%s)", err))
		}
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export
	genesis.AggregateVotes = k.GetAggregateVotes(ctx)
	genesis.AggregatePrevotes = k.GetAggregatePrevotes(ctx)
	genesis.MissCounts = k.GetMissCounts(ctx)
	genesis.BlockData = k.GetAllBlockData(ctx)
	genesis.FeederDelegation = k.GetFeederDelegations(ctx)

	return genesis
}
