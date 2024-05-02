package voteprocessor

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/settlus/chain/types"
	"github.com/settlus/chain/x/oracle/keeper"
	"github.com/settlus/chain/x/oracle/types"
)

func NewSettlusVoteProcessors(keeper keeper.Keeper, aggregateVotes []types.AggregateVote, thresholdVotes math.Int) []IVoteProcessor {
	blockConsensus := func(ctx sdk.Context, voteData map[string]types.BlockData) {
		for chainId, block := range voteData {
			if block.ChainId == "" { // for debug
				ctx.Logger().Error("block data is empty", "chainId", chainId, "height", ctx.BlockHeight())
				continue
			}
			keeper.SetBlockData(ctx, block)
		}
	}

	ownershipConsensus := func(ctx sdk.Context, voteData map[types.Nft]ctypes.HexAddressString) {
		keeper.FillSettlementRecipients(ctx, voteData)
	}

	return []IVoteProcessor{
		NewBlockVoteProcessor(blockConsensus, aggregateVotes, thresholdVotes),
		NewOwnershipVoteProcessor(ownershipConsensus, aggregateVotes, thresholdVotes),
	}
}

func NewBlockVoteProcessor(
	onConsensus func(ctx sdk.Context, voteData map[string]types.BlockData),
	aggregateVotes []types.AggregateVote,
	thresholdVotes math.Int) *VoteProcessor[string, types.BlockData] {

	return &VoteProcessor[string, types.BlockData]{
		topic:          types.OralceTopic_BLOCK,
		aggregateVotes: aggregateVotes,
		thresholdVotes: thresholdVotes,
		dataConverter:  types.StringToBlockData,
		onConsensus:    onConsensus,
	}
}

func NewOwnershipVoteProcessor(
	onConsensus func(ctx sdk.Context, voteData map[types.Nft]ctypes.HexAddressString),
	aggregateVotes []types.AggregateVote,
	thresholdVotes math.Int) *VoteProcessor[types.Nft, ctypes.HexAddressString] {

	return &VoteProcessor[types.Nft, ctypes.HexAddressString]{
		topic:          types.OralceTopic_OWNERSHIP,
		aggregateVotes: aggregateVotes,
		thresholdVotes: thresholdVotes,
		dataConverter:  types.StringToOwnershipData,
		onConsensus:    onConsensus,
	}
}
