package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/evmos/crypto/ethsecp256k1"

	"github.com/settlus/chain/x/oracle/types"
)

func (suite *OracleTestSuite) TestKeeper_Params() {
	paramsResponse, err := s.app.OracleKeeper.Params(s.ctx, &types.QueryParamsRequest{})
	suite.NoError(err)
	suite.Equal(types.DefaultParams(), paramsResponse.Params)
}

func (suite *OracleTestSuite) TestKeeper_BlockData() {
	expectedBlockData := types.BlockData{
		ChainId:     "1",
		BlockNumber: 100,
		BlockHash:   "foobar",
	}
	s.app.OracleKeeper.SetBlockData(s.ctx, expectedBlockData)

	actualBlockData, err := s.app.OracleKeeper.BlockData(s.ctx, &types.QueryBlockDataRequest{
		ChainId: "1",
	})
	suite.NoError(err)
	suite.Equal(expectedBlockData, *actualBlockData.BlockData)
}

func (suite *OracleTestSuite) TestKeeper_AllBlockData() {
	expectedBlockData := []types.BlockData{
		{
			ChainId:     "1",
			BlockNumber: 100,
			BlockHash:   "foobar",
		},
		{
			ChainId:     "2",
			BlockNumber: 200,
			BlockHash:   "barbaz",
		},
	}
	for _, blockData := range expectedBlockData {
		s.app.OracleKeeper.SetBlockData(s.ctx, blockData)
	}

	actualBlockData, err := s.app.OracleKeeper.AllBlockData(s.ctx, &types.QueryAllBlockDataRequest{})
	suite.NoError(err)
	suite.Equal(expectedBlockData, actualBlockData.BlockData)
}

func (suite *OracleTestSuite) TestKeeper_AggregatePrevote() {
	expectedAggregatePrevote := types.AggregatePrevote{
		Hash:        "foobar",
		Voter:       s.validators[0].GetOperator().String(),
		SubmitBlock: 100,
	}
	s.app.OracleKeeper.SetAggregatePrevote(s.ctx, expectedAggregatePrevote)

	actualAggregatePrevote, err := s.app.OracleKeeper.AggregatePrevote(s.ctx, &types.QueryAggregatePrevoteRequest{
		ValidatorAddress: s.validators[0].GetOperator().String(),
	})
	suite.NoError(err)
	suite.Equal(expectedAggregatePrevote, *actualAggregatePrevote.AggregatePrevote)
}

func (suite *OracleTestSuite) TestKeeper_AggregatePrevotes() {
	expectedAggregatePrevotes := []types.AggregatePrevote{
		{
			Hash:        "foobar",
			Voter:       s.validators[0].GetOperator().String(),
			SubmitBlock: 100,
		},
		{
			Hash:        "barbaz",
			Voter:       s.validators[1].GetOperator().String(),
			SubmitBlock: 200,
		},
	}
	for _, aggregatePrevote := range expectedAggregatePrevotes {
		s.app.OracleKeeper.SetAggregatePrevote(s.ctx, aggregatePrevote)
	}

	actualAggregatePrevotes, err := s.app.OracleKeeper.AggregatePrevotes(s.ctx, &types.QueryAggregatePrevotesRequest{})
	suite.NoError(err)
	suite.Equal(len(expectedAggregatePrevotes), len(actualAggregatePrevotes.AggregatePrevotes))
}

func (suite *OracleTestSuite) TestKeeper_AggregateVote() {
	expectedAggregateVote := types.AggregateVote{
		BlockData: []*types.BlockData{
			{
				ChainId:     "1",
				BlockNumber: 100,
				BlockHash:   "foobar",
			},
		},
		Voter: s.validators[0].GetOperator().String(),
	}
	s.app.OracleKeeper.SetAggregateVote(s.ctx, expectedAggregateVote)
	actualAggregateVote, err := s.app.OracleKeeper.AggregateVote(s.ctx, &types.QueryAggregateVoteRequest{
		ValidatorAddress: s.validators[0].GetOperator().String(),
	})
	suite.NoError(err)
	suite.Equal(expectedAggregateVote, *actualAggregateVote.AggregateVote)
}

func (suite *OracleTestSuite) TestKeeper_AggregateVotes() {
	expectedAggregateVotes := []types.AggregateVote{
		{
			BlockData: []*types.BlockData{
				{
					ChainId:     "1",
					BlockNumber: 100,
					BlockHash:   "foobar",
				},
			},
			Voter: s.validators[0].GetOperator().String(),
		},
		{
			BlockData: []*types.BlockData{
				{
					ChainId:     "2",
					BlockNumber: 200,
					BlockHash:   "barbaz",
				},
			},
			Voter: s.validators[1].GetOperator().String(),
		},
	}
	for _, aggregateVote := range expectedAggregateVotes {
		s.app.OracleKeeper.SetAggregateVote(s.ctx, aggregateVote)
	}

	actualAggregateVotes, err := s.app.OracleKeeper.AggregateVotes(s.ctx, &types.QueryAggregateVotesRequest{})
	suite.NoError(err)
	suite.Equal(len(expectedAggregateVotes), len(actualAggregateVotes.AggregateVotes))
}

func (suite *OracleTestSuite) TestKeeper_FeederDelegation() {
	priv, _ := ethsecp256k1.GenerateKey()
	feederAddress := sdk.AccAddress(priv.PubKey().Address().Bytes())
	err := s.app.OracleKeeper.SetFeederDelegation(s.ctx, s.validators[0].GetOperator().String(), feederAddress.String())
	suite.NoError(err)

	response, err := s.app.OracleKeeper.FeederDelegation(s.ctx, &types.QueryFeederDelegationRequest{
		ValidatorAddress: s.validators[0].GetOperator().String(),
	})
	suite.NoError(err)
	suite.Equal(feederAddress.String(), response.FeederDelegation.FeederAddress)
}

func (suite *OracleTestSuite) TestKeeper_MissCount() {
	s.app.OracleKeeper.SetMissCount(s.ctx, s.validators[0].GetOperator().String(), 100)

	response, err := s.app.OracleKeeper.MissCount(s.ctx, &types.QueryMissCountRequest{
		ValidatorAddress: s.validators[0].GetOperator().String(),
	})
	suite.NoError(err)
	suite.Equal(uint64(100), response.MissCount)
}

func (suite *OracleTestSuite) TestKeeper_RewardPool() {
	amt := sdk.NewCoins(sdk.NewInt64Coin("asetl", 1000000))
	err := s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, s.address, types.ModuleName, amt)
	suite.NoError(err)

	response, err := s.app.OracleKeeper.RewardPool(s.ctx, &types.QueryRewardPoolRequest{})
	suite.NoError(err)
	suite.Equal(amt, response.Balance)
}
