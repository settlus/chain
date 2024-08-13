package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/settlus/chain/x/oracle/keeper"
	"github.com/settlus/chain/x/oracle/types"
)

func (suite *OracleTestSuite) TestMsgServer_Prevote() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	suite.app.BeginBlocker(suite.ctx, abci.RequestBeginBlock{})
	hash := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"
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
	hash := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"

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
	hash := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"

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
	blockStr := []string{"1:100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"}
	ownershipStr := []string{"1/0x1234567890abcdef/0x1234567890abcdef:0x77791"}
	voteData := []*types.VoteData{
		{
			Topic: types.OracleTopic_BLOCK,
			Data:  blockStr,
		},
		{
			Topic: types.OracleTopic_OWNERSHIP,
			Data:  ownershipStr,
		},
	}

	suite.app.BeginBlocker(suite.ctx, abci.RequestBeginBlock{})
	suite.doPrevote(msgSvr, voteData, salt, 1)
	_, err := msgSvr.Vote(s.ctx.WithBlockHeight(10), &types.MsgVote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		VoteData:  voteData,
		Salt:      salt,
	})
	suite.NoError(err)

	vote := s.app.OracleKeeper.GetAggregateVote(s.ctx, s.validators[0].GetOperator().String())
	suite.Equal(vote.VoteData[0].Data, blockStr)
	suite.Equal(vote.VoteData[1].Data, ownershipStr)
}

func (suite *OracleTestSuite) TestMsgServer_Vote_should_be_failed_with_different_round_id() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	salt := "TestMsgServer_Vote"
	blockStr := []string{"1:100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"}
	ownershipStr := []string{"1/0x1234567890abcdef/0x1234567890abcdef:0x77791"}
	voteData := []*types.VoteData{
		{
			Topic: types.OracleTopic_BLOCK,
			Data:  blockStr,
		},
		{
			Topic: types.OracleTopic_OWNERSHIP,
			Data:  ownershipStr,
		},
	}

	suite.doPrevote(msgSvr, voteData, salt, 1)

	_, err := msgSvr.Vote(s.ctx.WithBlockHeight(39), &types.MsgVote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		VoteData:  voteData,
		Salt:      salt,
		RoundId:   1,
	})
	suite.Error(err)
}

func (suite *OracleTestSuite) TestMsgServer_Vote_should_be_failed_if_exceed_vote_period() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	salt := "TestMsgServer_Vote"
	blockStr := []string{"1:100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"}
	ownershipStr := []string{"1/0x1234567890abcdef/0x1234567890abcdef:0x77791"}
	voteData := []*types.VoteData{
		{
			Topic: types.OracleTopic_BLOCK,
			Data:  blockStr,
		},
		{
			Topic: types.OracleTopic_OWNERSHIP,
			Data:  ownershipStr,
		},
	}

	suite.doPrevote(msgSvr, voteData, salt, 1)

	_, err := msgSvr.Vote(s.ctx.WithBlockHeight(20), &types.MsgVote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		VoteData:  voteData,
		Salt:      salt,
		RoundId:   0,
	})
	suite.Error(err)
}

func (suite *OracleTestSuite) TestMsgServer_vote_should_fail_if_block_str_is_invalid() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	salt := "TestMsgServer_Vote"
	voteData := buildVoteData(
		"1100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
		"1/0x1234567890abcdef/0x1234567890abcdef:0x77791")
	suite.doPrevote(msgSvr, voteData, salt, 1)

	_, err := msgSvr.Vote(s.ctx.WithBlockHeight(10), &types.MsgVote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		VoteData:  voteData,
		Salt:      salt,
		RoundId:   0,
	})
	suite.Error(err)
}

func (suite *OracleTestSuite) TestMsgServer_vote_should_fail_if_nft_str_is_invalid() {
	msgSvr := keeper.NewMsgServerImpl(*suite.app.OracleKeeper)
	salt := "TestMsgServer_Vote"
	voteData := buildVoteData(
		"1:100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
		"10x1234567890abcdef/0x1234567890abcdef:0x77791")
	suite.doPrevote(msgSvr, voteData, salt, 1)

	_, err := msgSvr.Vote(s.ctx.WithBlockHeight(10), &types.MsgVote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		VoteData:  voteData,
		Salt:      salt,
		RoundId:   0,
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

func (suite *OracleTestSuite) doPrevote(msgSvr types.MsgServer, voteData []*types.VoteData, salt string, height int64) []string {
	hash, _ := types.GetAggregateVoteHash(voteData, salt)
	_, _ = msgSvr.Prevote(s.ctx.WithBlockHeight(height), &types.MsgPrevote{
		Feeder:    sdk.AccAddress(s.validators[0].GetOperator().Bytes()).String(),
		Validator: s.validators[0].GetOperator().String(),
		Hash:      hash,
	})

	return voteData[0].Data
}

func buildVoteData(blockStr, ownershipStr string) []*types.VoteData {
	return []*types.VoteData{
		{
			Topic: types.OracleTopic_BLOCK,
			Data:  []string{blockStr},
		},
		{
			Topic: types.OracleTopic_OWNERSHIP,
			Data:  []string{ownershipStr},
		},
	}
}
