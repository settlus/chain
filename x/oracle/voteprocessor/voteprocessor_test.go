package voteprocessor

import (
	"testing"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/settlus/chain/types"
	"github.com/settlus/chain/x/oracle/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func genValAddrs(num int) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, num)
	for i := 0; i < num; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		valAddrs[i] = sdk.ValAddress(pk.Address())
	}
	return valAddrs
}

func TestBlockVoteProcessor(t *testing.T) {
	thresholdVotes := math.NewInt(2)
	addrs := genValAddrs(3)

	aggregateVotes := []types.AggregateVote{
		{
			Voter: addrs[0].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OralceTopic_BLOCK,
					Data:  []string{"1:100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"},
				},
			},
		},
		{
			Voter: addrs[1].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OralceTopic_BLOCK,
					Data:  []string{"1:100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"},
				},
			},
		},
		{
			Voter: addrs[2].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OralceTopic_BLOCK,
					Data:  []string{"1:101/415f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"},
				},
			},
		},
	}

	validatorClaimMap := map[string]types.Claim{
		addrs[0].String(): {
			Weight:  1,
			Miss:    false,
			Abstain: false,
		},
		addrs[1].String(): {
			Weight:  1,
			Miss:    false,
			Abstain: false,
		},
		addrs[2].String(): {
			Weight:  1,
			Miss:    false,
			Abstain: false,
		},
	}

	result := make(map[string]types.BlockData)
	onConsensus := func(ctx sdk.Context, voteData map[string]types.BlockData) {
		for chainId, block := range voteData {
			result[chainId] = block
		}
	}
	vp := NewBlockVoteProcessor(onConsensus, aggregateVotes, thresholdVotes)
	ctx := sdk.NewContext(nil, tmproto.Header{Height: 1}, false, log.NewNopLogger())

	vp.TallyVotes(ctx, validatorClaimMap)

	require.Equal(t, "1", result["1"].ChainId)
	require.Equal(t, int64(100), result["1"].BlockNumber)
	require.Equal(t, "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3", result["1"].BlockHash)

	require.False(t, validatorClaimMap[addrs[0].String()].Miss)
	require.False(t, validatorClaimMap[addrs[1].String()].Miss)
	require.True(t, validatorClaimMap[addrs[2].String()].Miss)
}

func TestOwnershipVoteProcessor(t *testing.T) {
	thresholdVotes := math.NewInt(2)
	addrs := genValAddrs(3)

	aggregateVotes := []types.AggregateVote{
		{
			Voter: addrs[0].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OralceTopic_OWNERSHIP,
					Data:  []string{"1/0x123/0x0:0x777"},
				},
			},
		},
		{
			Voter: addrs[1].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OralceTopic_OWNERSHIP,
					Data:  []string{"1/0x123/0x0:0x777"},
				},
			},
		},
		{
			Voter: addrs[2].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OralceTopic_OWNERSHIP,
					Data:  []string{"1/0x123/0x0:0x666"},
				},
			},
		},
	}

	validatorClaimMap := map[string]types.Claim{
		addrs[0].String(): {
			Weight:  1,
			Miss:    false,
			Abstain: false,
		},
		addrs[1].String(): {
			Weight:  1,
			Miss:    false,
			Abstain: false,
		},
		addrs[2].String(): {
			Weight:  1,
			Miss:    false,
			Abstain: false,
		},
	}

	result := make(map[string]ctypes.HexAddressString)
	onConsensus := func(ctx sdk.Context, voteData map[ctypes.Nft]ctypes.HexAddressString) {
		for nft, owner := range voteData {
			result[nft.FormatString()] = owner
		}
	}
	vp := NewOwnershipVoteProcessor(onConsensus, aggregateVotes, thresholdVotes)
	ctx := sdk.NewContext(nil, tmproto.Header{Height: 1}, false, log.NewNopLogger())

	vp.TallyVotes(ctx, validatorClaimMap)

	require.Equal(t, ctypes.HexAddressString("0x777"), result["1/0x123/0x0"])

	require.False(t, validatorClaimMap[addrs[0].String()].Miss)
	require.False(t, validatorClaimMap[addrs[1].String()].Miss)
	require.True(t, validatorClaimMap[addrs[2].String()].Miss)
}
