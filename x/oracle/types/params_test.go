package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/settlus/chain/x/oracle/types"
)

func TestParams(t *testing.T) {
	p1 := types.DefaultParams()
	err := p1.Validate()
	require.NoError(t, err)

	// minus vote period
	p1.VotePeriod = 0
	err = p1.Validate()
	require.Error(t, err)

	// small vote threshold
	p2 := types.DefaultParams()
	p2.VoteThreshold = sdk.ZeroDec()
	err = p2.Validate()
	require.Error(t, err)

	// negative slash fraction
	p4 := types.DefaultParams()
	p4.SlashFraction = sdk.NewDec(-1)
	err = p4.Validate()
	require.Error(t, err)

	// max miss count per slash window is greater than slash window
	p5 := types.DefaultParams()
	p5.MaxMissCountPerSlashWindow = 100
	p5.SlashWindow = 10
	err = p5.Validate()
	require.Error(t, err)

	// max miss count per slash window is less than or equal to 0
	p6 := types.DefaultParams()
	p6.MaxMissCountPerSlashWindow = 0
	err = p6.Validate()
	require.Error(t, err)

	// slash window is less than or equal to 0
	p7 := types.DefaultParams()
	p7.SlashWindow = 0
	err = p7.Validate()
	require.Error(t, err)

	// slash window is less than vote period
	p8 := types.DefaultParams()
	p8.SlashWindow = 10
	p8.VotePeriod = 20
	err = p8.Validate()
	require.Error(t, err)

	// slash window is not divisible by vote period
	p9 := types.DefaultParams()
	p9.SlashWindow = 10
	p9.VotePeriod = 3
	err = p9.Validate()
	require.Error(t, err)

	// default params
	p10 := types.DefaultParams()
	require.NotNil(t, p10.ParamSetPairs())
	require.NotNil(t, p10.String())
}

func TestCalculateVotePeriod(t *testing.T) {
	tests := []struct {
		votePeriod  uint64
		blockHeight int64
		prevoteEnd  int64
		voteEnd     int64
	}{
		{votePeriod: 10, blockHeight: 0, prevoteEnd: 9, voteEnd: 19},
		{votePeriod: 10, blockHeight: 1, prevoteEnd: 9, voteEnd: 19},
		{votePeriod: 10, blockHeight: 9, prevoteEnd: 9, voteEnd: 19},
		{votePeriod: 10, blockHeight: 10, prevoteEnd: 9, voteEnd: 19},
		{votePeriod: 10, blockHeight: 19, prevoteEnd: 9, voteEnd: 19},
		{votePeriod: 10, blockHeight: 20, prevoteEnd: 29, voteEnd: 39},
		{votePeriod: 1, blockHeight: 0, prevoteEnd: 0, voteEnd: 1},
		{votePeriod: 1, blockHeight: 1, prevoteEnd: 0, voteEnd: 1},
		{votePeriod: 1, blockHeight: 2, prevoteEnd: 2, voteEnd: 3},
		{votePeriod: 1, blockHeight: 3, prevoteEnd: 2, voteEnd: 3},
		{votePeriod: 3, blockHeight: 0, prevoteEnd: 2, voteEnd: 5},
		{votePeriod: 3, blockHeight: 2, prevoteEnd: 2, voteEnd: 5},
		{votePeriod: 3, blockHeight: 4, prevoteEnd: 2, voteEnd: 5},
		{votePeriod: 3, blockHeight: 5, prevoteEnd: 2, voteEnd: 5},
	}

	for _, tt := range tests {
		prevoteEnd, voteEnd := types.CalculateVotePeriod(tt.blockHeight, tt.votePeriod)
		require.Equal(t, tt.prevoteEnd, prevoteEnd)
		require.Equal(t, tt.voteEnd, voteEnd)
	}
}
