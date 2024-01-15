package keeper_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/settlus/chain/x/oracle"
	oraclekeeper "github.com/settlus/chain/x/oracle/keeper"
	"github.com/settlus/chain/x/oracle/types"
)

const (
	DefaultVotePeriod = 10
)

var _ = Describe("Oracle module integration tests", Ordered, func() {
	BeforeEach(func() {
		s.SetupTest()
	})

	blockHash := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"
	blockDataStr := types.BlockDataToString(&types.BlockData{ChainId: "1", BlockNumber: 100, BlockHash: blockHash})

	Context("Test oracle consensus threshold", func() {
		BeforeEach(func() {
		})
		It("Less than the threshold signs, block data consensus fails", func() {
			salt := "1"
			validator := s.validators[0]
			SendPrevoteAndVote(validator, blockDataStr, salt)

			oracle.EndBlocker(s.ctx, *s.app.OracleKeeper)

			bd, err := s.app.OracleKeeper.GetBlockData(s.ctx, "1")
			Expect(err).To(Equal(types.ErrBlockDataNotFound))
			Expect(bd).To(BeNil())
		})
		It("All validators sign same block data, block data consensus succeeds", func() {
			// All validators sign
			for i, validator := range s.validators {
				salt := fmt.Sprintf("%d", i)
				SendPrevoteAndVote(validator, blockDataStr, salt)
			}

			oracle.EndBlocker(s.ctx.WithBlockHeight(DefaultVotePeriod), *s.app.OracleKeeper)

			bd, err := s.app.OracleKeeper.GetBlockData(s.ctx, "1")
			Expect(err).To(BeNil())
			Expect(*bd).To(Equal(types.BlockData{ChainId: "1", BlockNumber: 100, BlockHash: blockHash}))
		})
		It("Tie between two block data, block data consensus succeeds", func() {
			// Two validators sign block number 100
			for i, validator := range s.validators[:2] {
				salt := fmt.Sprintf("%d", i)
				SendPrevoteAndVote(validator, blockDataStr, salt)
			}

			// Two validators sign block number 99
			blockDataStr2 := types.BlockDataToString(&types.BlockData{ChainId: "1", BlockNumber: 99, BlockHash: blockHash})
			for i, validator := range s.validators[2:] {
				salt := fmt.Sprintf("%d", i+2)
				SendPrevoteAndVote(validator, blockDataStr2, salt)
			}

			oracle.EndBlocker(s.ctx.WithBlockHeight(DefaultVotePeriod), *s.app.OracleKeeper)

			bd, err := s.app.OracleKeeper.GetBlockData(s.ctx, "1")
			Expect(err).To(Equal(types.ErrBlockDataNotFound))
			Expect(bd).To(BeNil())
		})
		It("Abstain validator's power is majority, block data consensus fails", func() {
			// One validator abstains
			abstainValidator := s.validators[0]
			blockDataStr := types.BlockDataToString(&types.BlockData{ChainId: "1", BlockNumber: -1, BlockHash: blockHash})
			salt := "abstain"
			SendPrevoteAndVote(abstainValidator, blockDataStr, salt)

			// One validator signs block number 100
			validator := s.validators[1]
			salt = "1"
			SendPrevoteAndVote(validator, blockDataStr, salt)

			oracle.EndBlocker(s.ctx.WithBlockHeight(2), *s.app.OracleKeeper)

			bd, err := s.app.OracleKeeper.GetBlockData(s.ctx, "1")
			Expect(err).To(Equal(types.ErrBlockDataNotFound))
			Expect(bd).To(BeNil())
		})
	})

	Context("Test oracle vote period", func() {
		BeforeEach(func() {
			for i, validator := range s.validators {
				salt := fmt.Sprintf("%d", i)
				SendPrevote(s.ctx.WithBlockHeight(1), validator, blockDataStr, salt)
			}
			oracle.EndBlocker(s.ctx.WithBlockHeight(1), *s.app.OracleKeeper)
		})
		It("Vote period is 10, prevote is submitted at block 1, vote is submitted at block 10, block data is updated at block 11", func() {
			for i, validator := range s.validators {
				salt := fmt.Sprintf("%d", i)
				SendVote(s.ctx.WithBlockHeight(DefaultVotePeriod), validator, blockDataStr, salt)
			}

			oracle.EndBlocker(s.ctx.WithBlockHeight(DefaultVotePeriod), *s.app.OracleKeeper)

			bd, err := s.app.OracleKeeper.GetBlockData(s.ctx.WithBlockHeight(DefaultVotePeriod+1), "1")
			Expect(err).To(BeNil())
			Expect(*bd).To(Equal(types.BlockData{ChainId: "1", BlockNumber: 100, BlockHash: blockHash}))
		})
		It("Vote period is 10, prevote is submitted at block 1, vote is submitted at block 11, block data is not updated at block 12", func() {
			for i, validator := range s.validators {
				salt := fmt.Sprintf("%d", i)
				SendVote(s.ctx.WithBlockHeight(DefaultVotePeriod+1), validator, blockDataStr, salt)
			}

			oracle.EndBlocker(s.ctx.WithBlockHeight(DefaultVotePeriod+1), *s.app.OracleKeeper)

			bd, err := s.app.OracleKeeper.GetBlockData(s.ctx.WithBlockHeight(DefaultVotePeriod+2), "1")
			Expect(err).To(Equal(types.ErrBlockDataNotFound))
			Expect(bd).To(BeNil())

			// Check if old prevotes are deleted
			prevotes, err := s.app.OracleKeeper.AggregatePrevotes(s.ctx, &types.QueryAggregatePrevotesRequest{
				Pagination: &query.PageRequest{},
			})
			Expect(err).To(BeNil())
			Expect(prevotes.AggregatePrevotes).To(BeNil())
		})
		It("Vote period is 10, send prevote at block 1 and another at block 2. Send vote at block 10. Block data is updated at block 11", func() {
			blockDataStr := types.BlockDataToString(&types.BlockData{ChainId: "1", BlockNumber: 101, BlockHash: blockHash})

			for i, validator := range s.validators {
				salt := fmt.Sprintf("%d", i)
				SendPrevote(s.ctx.WithBlockHeight(2), validator, blockDataStr, salt)
			}
			oracle.EndBlocker(s.ctx.WithBlockHeight(2), *s.app.OracleKeeper)

			for i, validator := range s.validators {
				salt := fmt.Sprintf("%d", i)
				SendVote(s.ctx.WithBlockHeight(DefaultVotePeriod), validator, blockDataStr, salt)
			}
			oracle.EndBlocker(s.ctx.WithBlockHeight(DefaultVotePeriod), *s.app.OracleKeeper)

			bd, err := s.app.OracleKeeper.GetBlockData(s.ctx.WithBlockHeight(DefaultVotePeriod+1), "1")
			Expect(err).To(BeNil())
			Expect(*bd).To(Equal(types.BlockData{ChainId: "1", BlockNumber: 101, BlockHash: blockHash}))
		})
	})

})

func SendPrevote(ctx sdk.Context, validator stakingtypes.Validator, blockDataStr, salt string) {
	hash, _ := types.GetAggregateVoteHash(blockDataStr, salt)
	prevoteMsg := types.NewMsgPrevote(validator.OperatorAddress, validator.OperatorAddress, hash)
	oracleMsgSvr := oraclekeeper.NewMsgServerImpl(*s.app.OracleKeeper)
	_, err := oracleMsgSvr.Prevote(ctx, prevoteMsg)
	Expect(err).To(BeNil())
}

func SendVote(ctx sdk.Context, validator stakingtypes.Validator, blockDataStr, salt string) {
	voteMsg := types.NewMsgVote(validator.OperatorAddress, validator.OperatorAddress, blockDataStr, salt)
	oracleMsgSvr := oraclekeeper.NewMsgServerImpl(*s.app.OracleKeeper)
	_, err := oracleMsgSvr.Vote(ctx, voteMsg)
	Expect(err).To(BeNil())
}

// SendPrevoteAndVote sends prevote and vote messages to the oracle module
func SendPrevoteAndVote(validator stakingtypes.Validator, blockDataStr, salt string) {
	hash, _ := types.GetAggregateVoteHash(blockDataStr, salt)
	prevoteMsg := types.NewMsgPrevote(validator.OperatorAddress, validator.OperatorAddress, hash)
	voteMsg := types.NewMsgVote(validator.OperatorAddress, validator.OperatorAddress, blockDataStr, salt)

	oracleMsgSvr := oraclekeeper.NewMsgServerImpl(*s.app.OracleKeeper)
	_, err := oracleMsgSvr.Prevote(s.ctx.WithBlockHeight(1), prevoteMsg)
	Expect(err).To(BeNil())
	_, err = oracleMsgSvr.Vote(s.ctx.WithBlockHeight(DefaultVotePeriod), voteMsg)
	Expect(err).To(BeNil())
}
