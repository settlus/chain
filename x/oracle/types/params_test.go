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

	// empty ChainId
	p10 := types.DefaultParams()
	p10.Whitelist[0].ChainId = ""
	err = p10.Validate()
	require.Error(t, err)

	// empty ChainName
	p11 := types.DefaultParams()
	p11.Whitelist[0].ChainName = ""
	err = p11.Validate()
	require.Error(t, err)

	// empty ChainUrl
	p12 := types.DefaultParams()
	p12.Whitelist[0].ChainUrl = ""
	err = p12.Validate()
	require.Error(t, err)

	// duplicate ChainId
	p13 := types.DefaultParams()
	p13.Whitelist = append(p13.Whitelist, &types.Chain{
		ChainId:   "1",
		ChainName: "Foo",
		ChainUrl:  "http://localhost:8545",
	})
	err = p13.Validate()
	require.Error(t, err)

	// duplicate ChainName
	p14 := types.DefaultParams()
	p14.Whitelist = append(p13.Whitelist, &types.Chain{
		ChainId:   "2",
		ChainName: "Ethereum",
		ChainUrl:  "http://localhost:8545",
	})
	err = p14.Validate()
	require.Error(t, err)

	// max miss count per slash window is greater than slash window
	p15 := types.DefaultParams()
	p15.MaxMissCountPerSlashWindow = 100
	p15.SlashWindow = 10
	err = p15.Validate()
	require.Error(t, err)

	// max miss count per slash window is less than or equal to 0
	p16 := types.DefaultParams()
	p16.MaxMissCountPerSlashWindow = 0
	err = p16.Validate()
	require.Error(t, err)

	// slash window is less than or equal to 0
	p17 := types.DefaultParams()
	p17.SlashWindow = 0
	err = p17.Validate()
	require.Error(t, err)

	// slash window is less than vote period
	p18 := types.DefaultParams()
	p18.SlashWindow = 10
	p18.VotePeriod = 20
	err = p18.Validate()
	require.Error(t, err)

	// slash window is not divisible by vote period
	p19 := types.DefaultParams()
	p19.SlashWindow = 10
	p19.VotePeriod = 3
	err = p19.Validate()
	require.Error(t, err)

	// default params
	p20 := types.DefaultParams()
	require.NotNil(t, p20.ParamSetPairs())
	require.NotNil(t, p20.String())
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
