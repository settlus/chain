package types_test

import (
	"context"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/settlus/chain/x/oracle/types"
)

func TestParseBlockDataString(t *testing.T) {
	tests := []struct {
		name         string
		blockData    []*types.BlockData
		blockDatastr string
		wantErr      bool
	}{
		{
			name: "valid block data string",
			blockData: []*types.BlockData{
				{
					ChainId:     "1",
					BlockNumber: 100,
					BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
				},
				{
					ChainId:     "2",
					BlockNumber: 200,
					BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
				},
			},
			blockDatastr: "1:100:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3,2:200:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
			wantErr:      false,
		}, {
			name:         "invalid block data string: no comma",
			blockData:    nil,
			blockDatastr: "1:100:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd32:200:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
			wantErr:      true,
		}, {
			name:         "invalid block data string: non-integer block number",
			blockData:    nil,
			blockDatastr: "1:100.5:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3,2:200.1:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd, err := types.ParseBlockDataString(tt.blockDatastr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.blockData, bd)
			}
		})
	}
}

func TestGetAggregateVoteHash(t *testing.T) {
	tests := []struct {
		name         string
		blockDataStr string
		salt         string
		wantHash     string
	}{
		{
			name:         "valid block data string",
			salt:         "1",
			blockDataStr: "1:100:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3,2:200:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
			wantHash:     "f54d660a42b63a4eb17bbea8337b108509ae70348b08e1dbd7e64f2d02bb13b4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := types.GetAggregateVoteHash(tt.blockDataStr, tt.salt)
			require.NoError(t, err)
			require.Equal(t, hash, strings.ToUpper(tt.wantHash))
		})
	}
}

func TestBlockDataToString(t *testing.T) {
	tests := []struct {
		name string
		bd   *types.BlockData
		want string
	}{
		{
			name: "valid block data",
			bd: &types.BlockData{
				ChainId:     "1",
				BlockNumber: 100,
				BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
			},
			want: "1:100:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, types.BlockDataToString(tt.bd))
		})
	}
}

func TestIsLastBlockOfVotePeriod(t *testing.T) {
	tests := []struct {
		name        string
		blockNumber uint64
		period      uint64
		want        bool
	}{
		{
			name:        "block number is last block of period",
			blockNumber: 100,
			period:      5,
			want:        true,
		}, {
			name:        "block number is not last block of period",
			blockNumber: 101,
			period:      5,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := sdk.Context{}.WithContext(context.Background()).WithBlockHeight(int64(tt.blockNumber))
			require.Equal(t, tt.want, types.IsLastBlockOfVotePeriod(ctx, tt.period))
		})
	}
}

func TestIsLastBlockOfSlashWindow(t *testing.T) {
	tests := []struct {
		name        string
		blockNumber uint64
		slashWindow uint64
		want        bool
	}{
		{
			name:        "block number is last block of slashWindow",
			blockNumber: 1000,
			slashWindow: 100,
			want:        true,
		}, {
			name:        "block number is not last block of slashWindow",
			blockNumber: 1001,
			slashWindow: 100,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := sdk.Context{}.WithContext(context.Background()).WithBlockHeight(int64(tt.blockNumber))
			require.Equal(t, tt.want, types.IsLastBlockOfSlashWindow(ctx, tt.slashWindow))
		})
	}
}

func TestPickMostVotedBlockData(t *testing.T) {
	tests := []struct {
		name           string
		votes          []*types.BlockDataAndWeight
		thresholdVotes sdkmath.Int
		want           *types.BlockData
	}{
		{
			name: "default case",
			votes: []*types.BlockDataAndWeight{
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 100,
						BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
					},
					Weight: 100,
				},
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 200,
						BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894eddd",
					},
					Weight: 200,
				},
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 200,
						BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894eddd",
					},
					Weight: 100,
				},
			},
			thresholdVotes: sdkmath.NewInt(200),
			want: &types.BlockData{
				ChainId:     "1",
				BlockNumber: 200,
				BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894eddd",
			},
		}, {
			name: "tie should return nil",
			votes: []*types.BlockDataAndWeight{
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 100,
						BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
					},
					Weight: 100,
				},
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 200,
						BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894eddd",
					},
					Weight: 100,
				},
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 10,
						BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edde",
					},
					Weight: 100,
				},
			},
			thresholdVotes: sdkmath.NewInt(100),
			want:           nil,
		}, {
			name: "check block number as well as block hash",
			votes: []*types.BlockDataAndWeight{
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 100,
						BlockHash:   "1111111111111111",
					},
					Weight: 100,
				},
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 100,
						BlockHash:   "1111111111111111",
					},
					Weight: 100,
				},
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 100,
						BlockHash:   "2222222222222222",
					},
					Weight: 100,
				},
			},
			thresholdVotes: sdkmath.NewInt(150),
			want: &types.BlockData{
				ChainId:     "1",
				BlockNumber: 100,
				BlockHash:   "1111111111111111",
			},
		}, {
			name: "same block number, different hash",
			votes: []*types.BlockDataAndWeight{
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 100,
						BlockHash:   "1111111111111111",
					},
					Weight: 100,
				},
				{
					BlockData: &types.BlockData{
						ChainId:     "1",
						BlockNumber: 100,
						BlockHash:   "2222222222222222",
					},
					Weight: 100,
				},
			},
			thresholdVotes: sdkmath.NewInt(100),
			want:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := types.PicMostVotedBlockData(tt.votes, tt.thresholdVotes)
			require.Equal(t, tt.want, bd)
		})
	}
}

func TestTallyVotes(t *testing.T) {
	validator1 := sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	validator2 := sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	validatorAbstain := sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	tests := []struct {
		name              string
		votesByChainId    map[string][]*types.BlockDataAndVoter
		validatorClaimMap map[string]types.Claim
		thresholdVotes    sdkmath.Int
		want              map[string]*types.BlockData
	}{
		{
			name: "valid block data string",
			votesByChainId: map[string][]*types.BlockDataAndVoter{
				"1": {
					{
						BlockData: &types.BlockData{
							ChainId:     "1",
							BlockNumber: 100,
							BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
						},
						Voter: validator1,
					},
					{
						BlockData: &types.BlockData{
							ChainId:     "1",
							BlockNumber: 100,
							BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
						},
						Voter: validator2,
					},
					{
						BlockData: &types.BlockData{
							ChainId:     "1",
							BlockNumber: -1,
							BlockHash:   "",
						},
						Voter: validatorAbstain,
					},
				},
				"2": {
					{
						BlockData: &types.BlockData{
							ChainId:     "2",
							BlockNumber: 200,
							BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
						},
						Voter: validator1,
					},
					{
						BlockData: &types.BlockData{
							ChainId:     "2",
							BlockNumber: 200,
							BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
						},
						Voter: validator2,
					},
					{
						BlockData: &types.BlockData{
							ChainId:     "2",
							BlockNumber: -1,
							BlockHash:   "",
						},
						Voter: validatorAbstain,
					},
				},
			},
			validatorClaimMap: map[string]types.Claim{
				validator1.String(): {
					Weight:    100,
					MissCount: 0,
					Recipient: validator1,
					Abstain:   false,
				},
				validator2.String(): {
					Weight:    200,
					MissCount: 0,
					Recipient: validator2,
					Abstain:   false,
				},
				validatorAbstain.String(): {
					Weight:    100,
					MissCount: 0,
					Recipient: validatorAbstain,
					Abstain:   true,
				},
			},
			thresholdVotes: sdkmath.NewInt(300),
			want: map[string]*types.BlockData{
				"1": {
					ChainId:     "1",
					BlockNumber: 100,
					BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
				},
				"2": {
					ChainId:     "2",
					BlockNumber: 200,
					BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, types.TallyVotes(tt.votesByChainId, tt.validatorClaimMap, tt.thresholdVotes))
		})
	}
}
