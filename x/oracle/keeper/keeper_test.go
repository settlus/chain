package keeper_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/evmos/evmos/v19/crypto/ethsecp256k1"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/cmd/settlusd/config"
	"github.com/settlus/chain/testutil"
	utiltx "github.com/settlus/chain/testutil/tx"
	"github.com/settlus/chain/utils"
	"github.com/settlus/chain/x/oracle"
	"github.com/settlus/chain/x/oracle/types"
)

type OracleTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.SettlusApp
	queryClient types.QueryClient

	address     sdk.AccAddress
	consAddress sdk.ConsAddress
	clientCtx   client.Context //nolint:unused
	priv        cryptotypes.PrivKey
	validators  []stakingtypes.Validator
	signer      keyring.Signer
}

var s *OracleTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(OracleTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *OracleTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func (suite *OracleTestSuite) DoSetupTest(t *testing.T) {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.priv = priv
	suite.address = sdk.AccAddress(priv.PubKey().Address().Bytes())
	suite.signer = utiltx.NewSigner(priv)

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())
	suite.consAddress = consAddress

	checkTx := false

	// init app
	suite.app = app.Setup(checkTx, nil, utils.MainnetChainID)

	// setup context
	header := testutil.NewHeader(
		1, time.Now().UTC(), "settlus_5371-1", suite.consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, header)

	// query clients
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.OracleKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	// bond denom
	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = config.BaseDenom
	if err = suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams); err != nil {
		return
	}

	oracleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)
	suite.app.AccountKeeper.SetModuleAccount(suite.ctx, oracleAcc)

	// fund signer acc to pay for tx fees
	amt := sdk.NewInt(int64(math.Pow10(18) * 2))
	err = testutil.FundAccount(
		suite.ctx,
		suite.app.BankKeeper,
		suite.priv.PubKey().Address().Bytes(),
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amt)),
	)
	suite.Require().NoError(err)

	staking.EndBlocker(suite.ctx, suite.app.StakingKeeper)
	oracle.EndBlocker(suite.ctx, *suite.app.OracleKeeper)

	validators := s.app.StakingKeeper.GetValidators(s.ctx, 99)
	var bondedValidators []stakingtypes.Validator
	for i, validator := range validators {
		// four validators are enough for test
		if len(bondedValidators) > 4 {
			break
		}

		if validator.Status != stakingtypes.Bonded {
			panic("validator " + fmt.Sprintf("%d", i) + " is not bonded")
		}

		validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
		err = suite.app.StakingKeeper.Hooks().AfterValidatorCreated(suite.ctx, validator.GetOperator())
		require.NoError(t, err)
		err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
		require.NoError(t, err)

		bondedValidators = append(bondedValidators, validator)
	}
	suite.validators = bondedValidators
}

// NewValidator creates a new validator with a given amount of bonded tokens.
func (suite *OracleTestSuite) NewValidator(amt sdkmath.Int) *stakingtypes.MsgCreateValidator {
	privEd := ed25519.GenPrivKey()
	multiplier := sdk.NewInt(2)
	err := testutil.FundAccount(suite.ctx, suite.app.BankKeeper, privEd.PubKey().Address().Bytes(), sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amt.Mul(multiplier))))
	suite.Require().NoError(err)
	valAddr := sdk.ValAddress(privEd.PubKey().Address().Bytes())
	msgCreate, err := stakingtypes.NewMsgCreateValidator(
		valAddr,
		privEd.PubKey(),
		sdk.NewCoin(config.BaseDenom, amt),
		stakingtypes.NewDescription("moniker", "indentity", "website", "security_contract", "details"),
		stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		sdkmath.OneInt(),
		sdkmath.ZeroInt(),
		false,
	)
	suite.Require().NoError(err)
	return msgCreate
}

func (suite *OracleTestSuite) TestKeeper_GetBlockData() {
	bd := types.BlockData{
		ChainId:     "1",
		BlockNumber: 100,
		BlockHash:   "foobar",
	}
	s.app.OracleKeeper.SetBlockData(s.ctx, bd)
	actual, err := s.app.OracleKeeper.GetBlockData(s.ctx, "1")
	s.Require().NoError(err)
	s.Equal(bd, *actual)
}

func (suite *OracleTestSuite) TestKeeper_GetAllBlockData() {
	var l []types.BlockData
	a := types.BlockData{
		ChainId:     "1",
		BlockNumber: 100,
		BlockHash:   "foobar",
	}
	s.app.OracleKeeper.SetBlockData(s.ctx, a)
	l = append(l, a)

	b := types.BlockData{
		ChainId:     "2",
		BlockNumber: 200,
		BlockHash:   "foobar",
	}
	s.app.OracleKeeper.SetBlockData(s.ctx, b)
	l = append(l, b)

	actual := s.app.OracleKeeper.GetAllBlockData(s.ctx)
	s.Equal(l, actual)
}

func (suite *OracleTestSuite) TestKeeper_GetFeederDelegation() {
	validator := s.validators[0]
	err := s.app.OracleKeeper.SetFeederDelegation(s.ctx, validator.GetOperator().String(), s.address.String())
	s.NoError(err)
	actual := s.app.OracleKeeper.GetFeederDelegation(s.ctx, validator.GetOperator().String())
	s.Equal(s.address, actual)
}

func (suite *OracleTestSuite) TestKeeper_GetSetDeleteMissCount() {
	validator := s.validators[0]
	s.app.OracleKeeper.SetMissCount(suite.ctx, validator.GetOperator().String(), 10)
	actual := s.app.OracleKeeper.GetMissCount(suite.ctx, validator.GetOperator().String())
	s.Equal(uint64(10), actual)

	s.app.OracleKeeper.DeleteMissCount(suite.ctx, validator.GetOperator().String())
	actual = s.app.OracleKeeper.GetMissCount(suite.ctx, validator.GetOperator().String())
	s.Equal(uint64(0), actual)
}

func (suite *OracleTestSuite) TestKeeper_GetRewardPool() {
	amt := sdk.NewCoins(sdk.NewInt64Coin("asetl", 1000000))
	err := s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, s.address, types.ModuleName, amt)
	s.NoError(err)

	s.Equal(amt, s.app.OracleKeeper.GetRewardPool(s.ctx))
}

func (suite *OracleTestSuite) TestKeeper_GetSetDeletePrevote() {
	validator := s.validators[0]
	prevote := types.AggregatePrevote{
		Hash:  "foobar",
		Voter: validator.GetOperator().String(),
	}
	s.app.OracleKeeper.SetAggregatePrevote(s.ctx, prevote)

	actual := s.app.OracleKeeper.GetAggregatePrevote(s.ctx, validator.GetOperator().String())
	s.Equal(prevote, *actual)

	s.app.OracleKeeper.DeleteAggregatePrevote(s.ctx, validator.GetOperator().String())
	actual = s.app.OracleKeeper.GetAggregatePrevote(s.ctx, validator.GetOperator().String())
	s.Nil(actual)

	all, err := s.app.OracleKeeper.AggregatePrevotes(s.ctx, &types.QueryAggregatePrevotesRequest{})
	s.NoError(err)
	s.Equal([]*types.AggregatePrevote(nil), all.AggregatePrevotes)
}

func (suite *OracleTestSuite) TestKeeper_GetSetDeleteVote() {
	validator := s.validators[0]
	vote := types.AggregateVote{
		VoteData: []*types.VoteData{
			{
				Topic: types.OralceTopic_BLOCK,
				Data:  []string{"1:100:foobar"},
			},
		},
		Voter: validator.GetOperator().String(),
	}
	s.app.OracleKeeper.SetAggregateVote(s.ctx, vote)

	actual := s.app.OracleKeeper.GetAggregateVote(s.ctx, validator.GetOperator().String())
	s.Equal(vote, *actual)

	s.app.OracleKeeper.DeleteAggregateVote(s.ctx, validator.GetOperator().String())
	actual = s.app.OracleKeeper.GetAggregateVote(s.ctx, validator.GetOperator().String())
	s.Nil(actual)

	all, err := s.app.OracleKeeper.AggregateVotes(s.ctx, &types.QueryAggregateVotesRequest{})
	s.NoError(err)
	s.Equal([]*types.AggregateVote(nil), all.AggregateVotes)
}

func (suite *OracleTestSuite) TestKeeper_ValidateFeeder() {
	validator := s.validators[0]
	err := s.app.OracleKeeper.SetFeederDelegation(s.ctx, validator.GetOperator().String(), s.address.String())
	s.NoError(err)
	ok, err := s.app.OracleKeeper.ValidateFeeder(s.ctx, s.address.String(), validator.GetOperator().String())
	s.NoError(err)
	s.True(ok)
}

func (suite *OracleTestSuite) TestKeeper_RewardBallotWinners() {
	tests := []struct {
		name      string
		vcm       map[string]types.Claim
		totalCoin sdk.Coins
		rewardMap map[string]sdk.DecCoins
	}{
		{
			name: "all validators with same weight get reward",
			vcm: map[string]types.Claim{
				s.validators[0].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[1].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[2].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[3].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
			},
			totalCoin: sdk.NewCoins(sdk.NewInt64Coin("asetl", 4000000)),
			rewardMap: map[string]sdk.DecCoins{
				s.validators[0].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[1].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[2].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[3].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
			},
		},
		{
			name: "all validators with different weights get reward",
			vcm: map[string]types.Claim{
				s.validators[0].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[1].GetOperator().String(): {
					Weight:  200,
					Miss:    false,
					Abstain: false,
				},
				s.validators[2].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[3].GetOperator().String(): {
					Weight:  200,
					Miss:    false,
					Abstain: false,
				},
			},
			totalCoin: sdk.NewCoins(sdk.NewInt64Coin("asetl", 4000000)),
			rewardMap: map[string]sdk.DecCoins{
				s.validators[0].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 666666)),
				s.validators[1].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1333333)),
				s.validators[2].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 666666)),
				s.validators[3].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1333333)),
			},
		},
		{
			name: "2 validators get reward, 2 validators do not",
			vcm: map[string]types.Claim{
				s.validators[0].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[1].GetOperator().String(): {
					Weight:  200,
					Miss:    false,
					Abstain: false,
				},
				s.validators[2].GetOperator().String(): {
					Weight:  100,
					Miss:    true,
					Abstain: false,
				},
				s.validators[3].GetOperator().String(): {
					Weight:  200,
					Miss:    false,
					Abstain: true,
				},
			},
			totalCoin: sdk.NewCoins(sdk.NewInt64Coin("asetl", 3000000)),
			rewardMap: map[string]sdk.DecCoins{
				s.validators[0].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[1].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 2000001)),
				s.validators[2].GetOperator().String(): sdk.NewDecCoins(),
				s.validators[3].GetOperator().String(): sdk.NewDecCoins(),
			},
		},
		{
			name: "no validators get reward",
			vcm: map[string]types.Claim{
				s.validators[0].GetOperator().String(): {
					Weight:  100,
					Miss:    true,
					Abstain: false,
				},
				s.validators[1].GetOperator().String(): {
					Weight:  200,
					Miss:    true,
					Abstain: false,
				},
				s.validators[2].GetOperator().String(): {
					Weight:  100,
					Miss:    true,
					Abstain: false,
				},
				s.validators[3].GetOperator().String(): {
					Weight:  200,
					Miss:    true,
					Abstain: false,
				},
			},
			totalCoin: sdk.NewCoins(sdk.NewInt64Coin("asetl", 1000000)),
			rewardMap: map[string]sdk.DecCoins{
				s.validators[0].GetOperator().String(): sdk.NewDecCoins(),
				s.validators[1].GetOperator().String(): sdk.NewDecCoins(),
				s.validators[2].GetOperator().String(): sdk.NewDecCoins(),
				s.validators[3].GetOperator().String(): sdk.NewDecCoins(),
			},
		},
	}

	err := testutil.FundAccount(s.ctx, s.app.BankKeeper, s.address, sdk.NewCoins(sdk.NewInt64Coin("asetl", 10000000000)))
	s.NoError(err)

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err = s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, s.address, types.ModuleName, tt.totalCoin)
			s.NoError(err)

			err = s.app.OracleKeeper.RewardBallotWinners(s.ctx, tt.vcm)
			s.NoError(err)

			for _, validator := range s.validators {
				rewards := s.app.DistrKeeper.GetValidatorCurrentRewards(s.ctx, validator.GetOperator())
				s.Equal(tt.rewardMap[validator.GetOperator().String()].AmountOf("asetl"), rewards.Rewards.AmountOf("asetl"))
				s.app.DistrKeeper.DeleteValidatorCurrentRewards(s.ctx, validator.GetOperator())
			}
		})
	}
}

func (suite *OracleTestSuite) TestKeeper_RewardBallotWinners_WithProbono() {
	tests := []struct {
		name         string
		vcm          map[string]types.Claim
		totalCoin    sdk.Coins
		rewardMap    map[string]sdk.DecCoins
		probonoIndex []int
	}{
		{
			name: "Probono validators send rewards to community pool, normal validators get rewards as usual",
			vcm: map[string]types.Claim{
				s.validators[0].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[1].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[2].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[3].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
			},
			totalCoin: sdk.NewCoins(sdk.NewInt64Coin("asetl", 4000000)),
			rewardMap: map[string]sdk.DecCoins{
				s.validators[0].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[1].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[2].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[3].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
			},
			probonoIndex: []int{0, 1},
		},
		{
			name: "All probono validator, there's no normal validator",
			vcm: map[string]types.Claim{
				s.validators[0].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[1].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[2].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
				s.validators[3].GetOperator().String(): {
					Weight:  100,
					Miss:    false,
					Abstain: false,
				},
			},
			totalCoin: sdk.NewCoins(sdk.NewInt64Coin("asetl", 4000000)),
			rewardMap: map[string]sdk.DecCoins{
				s.validators[0].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[1].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[2].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
				s.validators[3].GetOperator().String(): sdk.NewDecCoins(sdk.NewInt64DecCoin("asetl", 1000000)),
			},
			probonoIndex: []int{0, 1, 2, 3},
		},
	}

	err := testutil.FundAccount(s.ctx, s.app.BankKeeper, s.address, sdk.NewCoins(sdk.NewInt64Coin("asetl", 10000000000)))
	s.NoError(err)

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err = s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, s.address, types.ModuleName, tt.totalCoin)
			s.NoError(err)

			s.app.DistrKeeper.SetFeePool(s.ctx, disttypes.InitialFeePool())
			s.Equal(s.app.DistrKeeper.GetFeePoolCommunityCoins(s.ctx).AmountOf("asetl"), sdk.ZeroDec())

			for _, idx := range tt.probonoIndex {
				s.validators[idx].Probono = true
				s.validators[idx] = stakingkeeper.TestingUpdateValidator(s.app.StakingKeeper, s.ctx, s.validators[idx], true)
			}

			probonoRewards := sdk.ZeroDec()

			err = s.app.OracleKeeper.RewardBallotWinners(s.ctx, tt.vcm)
			s.NoError(err)

			for _, validator := range s.validators {
				if validator.IsProbono() {
					probonoRewards = probonoRewards.Add(tt.rewardMap[validator.GetOperator().String()].AmountOf("asetl"))
					validator.Probono = false
					continue
				}
				rewards := s.app.DistrKeeper.GetValidatorCurrentRewards(s.ctx, validator.GetOperator())
				s.Equal(tt.rewardMap[validator.GetOperator().String()].AmountOf("asetl"), rewards.Rewards.AmountOf("asetl"))
				s.app.DistrKeeper.DeleteValidatorCurrentRewards(s.ctx, validator.GetOperator())
			}

			// check community pool
			actualCommunityAmount := s.app.DistrKeeper.GetFeePoolCommunityCoins(s.ctx).AmountOf("asetl")
			s.Equal(probonoRewards, actualCommunityAmount)
			s.app.DistrKeeper.SetFeePool(s.ctx, disttypes.InitialFeePool())
		})
	}
}

func (suite *OracleTestSuite) TestKeeper_ClearBallots() {
	tests := []struct {
		name                  string
		setupPrevotesAndVotes func()
		wantLength            int
	}{
		{
			name:                  "no prevotes and votes",
			setupPrevotesAndVotes: func() {},
			wantLength:            0,
		}, {
			name: "all prevotes and votes need to be cleared",
			setupPrevotesAndVotes: func() {
				s.app.OracleKeeper.SetAggregatePrevote(s.ctx, types.AggregatePrevote{
					Hash:  "foobar",
					Voter: s.validators[0].GetOperator().String(),
				})
				s.app.OracleKeeper.SetAggregatePrevote(s.ctx, types.AggregatePrevote{
					Hash:  "foobar",
					Voter: s.validators[1].GetOperator().String(),
				})
				s.app.OracleKeeper.SetAggregateVote(s.ctx, types.AggregateVote{
					VoteData: []*types.VoteData{
						{
							Topic: types.OralceTopic_BLOCK,
							Data:  []string{"1:100:foobar"},
						},
					},
				})
			},
			wantLength: 0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupPrevotesAndVotes()
			s.app.OracleKeeper.ClearBallots(s.ctx.WithBlockHeight(3))

			prevotes, err := s.app.OracleKeeper.AggregatePrevotes(s.ctx.WithBlockHeight(4), &types.QueryAggregatePrevotesRequest{Pagination: &query.PageRequest{}})
			s.NoError(err)
			s.Equal(tt.wantLength, len(prevotes.AggregatePrevotes))

			votes, err := s.app.OracleKeeper.AggregateVotes(s.ctx.WithBlockHeight(4), &types.QueryAggregateVotesRequest{Pagination: &query.PageRequest{}})
			s.NoError(err)
			s.Equal(0, len(votes.AggregateVotes))

			// delete all prevotes for next tests
			for _, validator := range s.validators {
				s.app.OracleKeeper.DeleteAggregatePrevote(s.ctx, validator.GetOperator().String())
			}
		})
	}
}

func (suite *OracleTestSuite) TestKeeper_SlashValidatorsAndResetMissCount() {
	amt := sdk.DefaultPowerReduction

	tests := []struct {
		name                 string
		setupValidators      func()
		validatorsMissCount  map[string]uint64
		wantValidatorSlashed map[string]bool
		wantValidatorJailed  map[string]bool
	}{
		{
			name:            "no validators need to be slashed",
			setupValidators: func() {},
			validatorsMissCount: map[string]uint64{
				s.validators[0].GetOperator().String(): 0,
				s.validators[1].GetOperator().String(): 0,
				s.validators[2].GetOperator().String(): 0,
				s.validators[3].GetOperator().String(): 0,
			},
			wantValidatorSlashed: map[string]bool{
				s.validators[0].GetOperator().String(): false,
				s.validators[1].GetOperator().String(): false,
				s.validators[2].GetOperator().String(): false,
				s.validators[3].GetOperator().String(): false,
			},
			wantValidatorJailed: map[string]bool{
				s.validators[0].GetOperator().String(): false,
				s.validators[1].GetOperator().String(): false,
				s.validators[2].GetOperator().String(): false,
				s.validators[3].GetOperator().String(): false,
			},
		}, {
			name:            "all validators need to be slashed, default max miss count is 60",
			setupValidators: func() {},
			validatorsMissCount: map[string]uint64{
				s.validators[0].GetOperator().String(): 100,
				s.validators[1].GetOperator().String(): 100,
				s.validators[2].GetOperator().String(): 100,
				s.validators[3].GetOperator().String(): 100,
			},
			wantValidatorSlashed: map[string]bool{
				s.validators[0].GetOperator().String(): true,
				s.validators[1].GetOperator().String(): true,
				s.validators[2].GetOperator().String(): true,
				s.validators[3].GetOperator().String(): true,
			},
			wantValidatorJailed: map[string]bool{
				s.validators[0].GetOperator().String(): true,
				s.validators[1].GetOperator().String(): true,
				s.validators[2].GetOperator().String(): true,
				s.validators[3].GetOperator().String(): true,
			},
		}, {
			name:            "some validators need to be slashed, default max miss count is 60",
			setupValidators: func() {},
			validatorsMissCount: map[string]uint64{
				s.validators[0].GetOperator().String(): 100,
				s.validators[1].GetOperator().String(): 61,
				s.validators[2].GetOperator().String(): 5,
				s.validators[3].GetOperator().String(): 10,
			},
			wantValidatorSlashed: map[string]bool{
				s.validators[0].GetOperator().String(): true,
				s.validators[1].GetOperator().String(): true,
				s.validators[2].GetOperator().String(): false,
				s.validators[3].GetOperator().String(): false,
			},
			wantValidatorJailed: map[string]bool{
				s.validators[0].GetOperator().String(): true,
				s.validators[1].GetOperator().String(): true,
				s.validators[2].GetOperator().String(): false,
				s.validators[3].GetOperator().String(): false,
			},
		}, {
			name: "do not slash unbonded validators",
			setupValidators: func() {
				// unbond validator 0
				validator, _ := s.app.StakingKeeper.GetValidator(s.ctx, s.validators[0].GetOperator())
				validator.Status = stakingtypes.Unbonded
				validator.Jailed = false
				validator.Tokens = amt
				s.app.StakingKeeper.SetValidator(s.ctx, validator)
			},
			validatorsMissCount: map[string]uint64{
				s.validators[0].GetOperator().String(): 100,
				s.validators[1].GetOperator().String(): 100,
			},
			wantValidatorSlashed: map[string]bool{
				s.validators[0].GetOperator().String(): false,
				s.validators[1].GetOperator().String(): true,
			},
			wantValidatorJailed: map[string]bool{
				s.validators[0].GetOperator().String(): false,
				s.validators[1].GetOperator().String(): true,
			},
		}, {
			name: "do not slash jailed validators",
			setupValidators: func() {
				// jail validator 0
				validator, _ := s.app.StakingKeeper.GetValidator(s.ctx, s.validators[0].GetOperator())
				validator.Status = stakingtypes.Bonded
				validator.Jailed = true
				validator.Tokens = amt
				s.app.StakingKeeper.SetValidator(s.ctx, validator)
			},
			validatorsMissCount: map[string]uint64{
				s.validators[0].GetOperator().String(): 100,
				s.validators[1].GetOperator().String(): 100,
			},
			wantValidatorSlashed: map[string]bool{
				s.validators[0].GetOperator().String(): false,
				s.validators[1].GetOperator().String(): true,
			},
			wantValidatorJailed: map[string]bool{
				s.validators[0].GetOperator().String(): true,
				s.validators[1].GetOperator().String(): true,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupValidators()
			for _, validator := range s.validators {
				s.app.OracleKeeper.SetMissCount(s.ctx, validator.GetOperator().String(), tt.validatorsMissCount[validator.GetOperator().String()])
			}

			s.app.OracleKeeper.SlashValidatorsAndResetMissCount(s.ctx)

			// check if validators are jailed
			for _, validator := range s.validators {
				validator, found := s.app.StakingKeeper.GetValidator(s.ctx, validator.GetOperator())
				s.True(found)

				if tt.wantValidatorSlashed[validator.GetOperator().String()] {
					s.Less(validator.GetTokens().Uint64(), amt.Uint64())
				} else {
					s.Equal(validator.GetTokens().Uint64(), amt.Uint64())
				}

				s.Equal(tt.wantValidatorJailed[validator.GetOperator().String()], validator.Jailed)
			}

			// check if the miss count is reset
			for _, validator := range s.validators {
				missCount := s.app.OracleKeeper.GetMissCount(s.ctx, validator.GetOperator().String())
				s.Equal(uint64(0), missCount)
			}

			// reset validator status
			// this shows why ginkgo and gomega are better than testify. we can use BeforeEach and AfterEach in ginkgo
			for _, validator := range s.validators {
				validator, _ := s.app.StakingKeeper.GetValidator(s.ctx, validator.GetOperator())
				validator.Status = stakingtypes.Bonded
				validator.Jailed = false
				validator.Tokens = amt
				s.app.StakingKeeper.SetValidator(s.ctx, validator)
			}
		})
	}
}
