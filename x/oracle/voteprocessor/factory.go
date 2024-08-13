package voteprocessor

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/settlus/chain/types"
	"github.com/settlus/chain/x/oracle/keeper"
	"github.com/settlus/chain/x/oracle/types"
)

func NewSettlusVoteProcessors(keeper keeper.Keeper, aggregateVotes []types.AggregateVote, thresholdVotes math.Int) []IVoteProcessor {
	ownershipConsensus := func(ctx sdk.Context, voteData map[ctypes.Nft]ctypes.HexAddressString) {
		keeper.FillSettlementRecipients(ctx, voteData)
	}

	return []IVoteProcessor{
		NewOwnershipVoteProcessor(ownershipConsensus, aggregateVotes, thresholdVotes),
	}
}

func NewOwnershipVoteProcessor(
	onConsensus func(ctx sdk.Context, voteData map[ctypes.Nft]ctypes.HexAddressString),
	aggregateVotes []types.AggregateVote,
	thresholdVotes math.Int) *VoteProcessor[ctypes.Nft, ctypes.HexAddressString] {

	return &VoteProcessor[ctypes.Nft, ctypes.HexAddressString]{
		topic:          types.OracleTopic_OWNERSHIP,
		aggregateVotes: aggregateVotes,
		thresholdVotes: thresholdVotes,
		dataConverter:  types.StringToOwnershipData,
		onConsensus:    onConsensus,
	}
}
