package voteprocessor

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	ctypes "github.com/settlus/chain/types"
	"github.com/settlus/chain/x/oracle/types"
)

func genValAddrs(num int) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, num)
	for i := 0; i < num; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		valAddrs[i] = sdk.ValAddress(pk.Address())
	}
	return valAddrs
}

func TestOwnershipVoteProcessor(t *testing.T) {
	thresholdVotes := math.NewInt(2)
	addrs := genValAddrs(3)

	aggregateVotes := []types.AggregateVote{
		{
			Voter: addrs[0].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OracleTopic_OWNERSHIP,
					Data:  []string{"1/0x123/0x0:0x777"},
				},
			},
		},
		{
			Voter: addrs[1].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OracleTopic_OWNERSHIP,
					Data:  []string{"1/0x123/0x0:0x777"},
				},
			},
		},
		{
			Voter: addrs[2].String(),
			VoteData: []*types.VoteData{
				{
					Topic: types.OracleTopic_OWNERSHIP,
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

	require.Equal(t, ctypes.HexAddressString("0x0000000000000000000000000000000000000777"),
		result["1/0x0000000000000000000000000000000000000123/0x0000000000000000000000000000000000000000"])

	require.False(t, validatorClaimMap[addrs[0].String()].Miss)
	require.False(t, validatorClaimMap[addrs[1].String()].Miss)
	require.True(t, validatorClaimMap[addrs[2].String()].Miss)
}
