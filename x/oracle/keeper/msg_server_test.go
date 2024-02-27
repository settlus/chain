package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/oracle/keeper"
	"github.com/settlus/chain/x/oracle/types"
)

func (suite *OracleTestSuite) TestMsgServer_Prevote() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	blockHash := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"
	blockDataStr := types.BlockDataToString(&types.BlockData{
		ChainId: "1", BlockNumber: 100, BlockHash: blockHash,
	})
	salt := "TestMsgServer_Prevote"
	hash, _ := types.GetAggregateVoteHash(blockDataStr, salt)

	response, err := msgSvr.Prevote(s.ctx.WithBlockHeight(1), &types.MsgPrevote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		Hash:      hash,
	})
	suite.NoError(err)
	suite.Equal(response, &types.MsgPrevoteResponse{})

	prevote := s.app.OracleKeeper.GetAggregatePrevote(s.ctx, s.validators[0].GetOperator().String())
	suite.Equal(prevote.Hash, hash)
}

func (suite *OracleTestSuite) TestMsgServer_Prevote_should_be_failed_with_different_round_id() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	blockHash := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"
	blockDataStr := types.BlockDataToString(&types.BlockData{
		ChainId: "1", BlockNumber: 100, BlockHash: blockHash,
	})
	salt := "TestMsgServer_Prevote"
	hash, _ := types.GetAggregateVoteHash(blockDataStr, salt)

	_, err := msgSvr.Prevote(s.ctx.WithBlockHeight(20), &types.MsgPrevote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		Hash:      hash,
		RoundId:   0,
	})
	suite.Error(err)
}

func (suite *OracleTestSuite) TestMsgServer_Prevote_should_be_failed_if_exceed_prevote_period() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	blockHash := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"
	blockDataStr := types.BlockDataToString(&types.BlockData{
		ChainId: "1", BlockNumber: 100, BlockHash: blockHash,
	})
	salt := "TestMsgServer_Prevote"
	hash, _ := types.GetAggregateVoteHash(blockDataStr, salt)

	_, err := msgSvr.Prevote(s.ctx.WithBlockHeight(10), &types.MsgPrevote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		Hash:      hash,
		RoundId:   0,
	})
	suite.Error(err)
}

func (suite *OracleTestSuite) TestMsgServer_Vote() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	salt := "TestMsgServer_Vote"
	blockData := &types.BlockData{
		ChainId: "1", BlockNumber: 100, BlockHash: "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
	}
	blockDataStr := suite.doPrevote(msgSvr, blockData, salt, 1)

	_, err := msgSvr.Vote(s.ctx.WithBlockHeight(10), &types.MsgVote{
		Feeder:          sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator:       s.validators[0].GetOperator().String(),
		BlockDataString: blockDataStr,
		Salt:            salt,
	})
	suite.NoError(err)

	vote := s.app.OracleKeeper.GetAggregateVote(s.ctx, s.validators[0].GetOperator().String())
	suite.Equal(vote.BlockData[0], blockData)
}

func (suite *OracleTestSuite) TestMsgServer_Vote_should_be_failed_with_different_round_id() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	salt := "TestMsgServer_Vote"
	blockData := &types.BlockData{
		ChainId: "1", BlockNumber: 100, BlockHash: "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
	}
	blockDataStr := suite.doPrevote(msgSvr, blockData, salt, 1)

	_, err := msgSvr.Vote(s.ctx.WithBlockHeight(39), &types.MsgVote{
		Feeder:          sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator:       s.validators[0].GetOperator().String(),
		BlockDataString: blockDataStr,
		Salt:            salt,
		RoundId:         1,
	})
	suite.Error(err)
}

func (suite *OracleTestSuite) TestMsgServer_Vote_should_be_failed_if_exceed_vote_period() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	salt := "TestMsgServer_Vote"
	blockData := &types.BlockData{
		ChainId: "1", BlockNumber: 100, BlockHash: "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
	}
	blockDataStr := suite.doPrevote(msgSvr, blockData, salt, 1)

	_, err := msgSvr.Vote(s.ctx.WithBlockHeight(20), &types.MsgVote{
		Feeder:          sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator:       s.validators[0].GetOperator().String(),
		BlockDataString: blockDataStr,
		Salt:            salt,
		RoundId:         0,
	})
	suite.Error(err)
}

func (suite *OracleTestSuite) TestMsgServer_FeederDelegationConsent() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)

	_, err := msgSvr.FeederDelegationConsent(s.ctx, &types.MsgFeederDelegationConsent{
		Validator:     s.validators[0].GetOperator().String(),
		FeederAddress: s.address.String(),
	})
	suite.NoError(err)

	feeder := s.app.OracleKeeper.GetFeederDelegation(s.ctx, s.validators[0].GetOperator().String())
	suite.Equal(feeder, s.address)
}

func (suite *OracleTestSuite) doPrevote(msgSvr types.MsgServer, blockData *types.BlockData, salt string, height int64) string {
	blockDataStr := types.BlockDataToString(blockData)
	hash, _ := types.GetAggregateVoteHash(blockDataStr, salt)

	_, _ = msgSvr.Prevote(s.ctx.WithBlockHeight(height), &types.MsgPrevote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		Hash:      hash,
	})

	return blockDataStr
}
