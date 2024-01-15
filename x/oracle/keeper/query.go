package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/settlus/chain/x/oracle/types"
)

var _ types.QueryServer = Keeper{}

// Params queries params of oracle module
func (k Keeper) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// BlockData queries block data of oracle module
func (k Keeper) BlockData(goCtx context.Context, req *types.QueryBlockDataRequest) (*types.QueryBlockDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	blockData, err := k.GetBlockData(ctx, req.ChainId)
	if err != nil {
		return nil, fmt.Errorf("error while getting block data: %w", err)
	}
	return &types.QueryBlockDataResponse{BlockData: blockData}, nil
}

// AllBlockData queries all block data of oracle module
func (k Keeper) AllBlockData(goCtx context.Context, _ *types.QueryAllBlockDataRequest) (*types.QueryAllBlockDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	blockData := k.GetAllBlockData(ctx)
	return &types.QueryAllBlockDataResponse{BlockData: blockData}, nil
}

// AggregatePrevote queries aggregate prevote of a given validator
func (k Keeper) AggregatePrevote(goCtx context.Context, request *types.QueryAggregatePrevoteRequest) (*types.QueryAggregatePrevoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	aggregatePrevote := k.GetAggregatePrevote(ctx, request.ValidatorAddress)
	return &types.QueryAggregatePrevoteResponse{AggregatePrevote: aggregatePrevote}, nil
}

// AggregatePrevotes queries aggregate prevotes of all validators
func (k Keeper) AggregatePrevotes(goCtx context.Context, request *types.QueryAggregatePrevotesRequest) (*types.QueryAggregatePrevotesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	aggPrevoteStore := k.GetAggregatePrevoteStore(ctx)
	var paginatedAggPrevotes []*types.AggregatePrevote
	pageRes, err := query.Paginate(aggPrevoteStore, request.Pagination, func(_ []byte, value []byte) error {
		var aggregatePrevote types.AggregatePrevote
		k.cdc.MustUnmarshal(value, &aggregatePrevote)

		paginatedAggPrevotes = append(paginatedAggPrevotes, &aggregatePrevote)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while paginating aggregate prevotes: %w", err)
	}

	return &types.QueryAggregatePrevotesResponse{
		AggregatePrevotes: paginatedAggPrevotes,
		Pagination:        pageRes,
	}, nil
}

// AggregateVote queries aggregate vote of a given validator
func (k Keeper) AggregateVote(goCtx context.Context, request *types.QueryAggregateVoteRequest) (*types.QueryAggregateVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	aggregateVote := k.GetAggregateVote(ctx, request.ValidatorAddress)
	return &types.QueryAggregateVoteResponse{AggregateVote: aggregateVote}, nil
}

// AggregateVotes queries aggregate votes of all validators
func (k Keeper) AggregateVotes(goCtx context.Context, request *types.QueryAggregateVotesRequest) (*types.QueryAggregateVotesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	aggVoteStore := k.GetAggregateVoteStore(ctx)

	var paginatedAggVotes []*types.AggregateVote
	pageRes, err := query.Paginate(aggVoteStore, request.Pagination, func(_ []byte, value []byte) error {
		var aggregateVote types.AggregateVote
		err := k.cdc.Unmarshal(value, &aggregateVote)
		if err != nil {
			return fmt.Errorf("error while unmarshalling aggregate vote: %w", err)
		}

		paginatedAggVotes = append(paginatedAggVotes, &aggregateVote)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while paginating aggregate votes: %w", err)
	}

	return &types.QueryAggregateVotesResponse{
		AggregateVotes: paginatedAggVotes,
		Pagination:     pageRes,
	}, nil
}

// FeederDelegation queries feeder delegation of a given validator
func (k Keeper) FeederDelegation(goCtx context.Context, request *types.QueryFeederDelegationRequest) (*types.QueryFeederDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	feeder := k.GetFeederDelegation(ctx, request.ValidatorAddress)
	return &types.QueryFeederDelegationResponse{FeederDelegation: &types.FeederDelegation{FeederAddress: feeder.String(), ValidatorAddress: request.ValidatorAddress}}, nil
}

// MissCount queries miss count of a given validator
func (k Keeper) MissCount(goCtx context.Context, request *types.QueryMissCountRequest) (*types.QueryMissCountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	missCount := k.GetMissCount(ctx, request.ValidatorAddress)
	return &types.QueryMissCountResponse{MissCount: missCount}, nil
}

// RewardPool queries the reward pool balance
func (k Keeper) RewardPool(goCtx context.Context, request *types.QueryRewardPoolRequest) (*types.QueryRewardPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool := k.GetRewardPool(ctx)
	return &types.QueryRewardPoolResponse{Balance: pool}, nil
}
