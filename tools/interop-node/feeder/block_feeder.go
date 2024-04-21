package feeder

import (
	"context"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/settlus/chain/tools/interop-node/client"
	"github.com/settlus/chain/tools/interop-node/config"
	"github.com/settlus/chain/tools/interop-node/subscriber"
	"github.com/settlus/chain/tools/interop-node/types"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

const BlockTopic = "block"

type BlockVoteInfo struct {
	RoundId  uint64
	Salt     string
	VoteData types.VoteDataArr
}

type BlockFeeder struct {
	BaseFeeder

	lastPreVote *BlockVoteInfo
	lastVote    *BlockVoteInfo
	subscribers map[string]subscriber.Subscriber
}

var _ Feeder = (*BlockFeeder)(nil)

func NewBlockFeeder(
	config *config.Config,
	sc *client.SettlusClient,
	subscribers []subscriber.Subscriber,
	logger log.Logger,
) (*BlockFeeder, error) {
	BaseFeeder := BaseFeeder{
		logger:           logger.With("topic", "block"),
		validatorAddress: config.Feeder.ValidatorAddress,
		address:          config.Feeder.Address,
		sc:               sc,
	}

	subscribersMap := make(map[string]subscriber.Subscriber)
	for _, cc := range subscribers {
		subscribersMap[cc.Id()] = cc
	}

	return &BlockFeeder{
		BaseFeeder:  BaseFeeder,
		subscribers: subscribersMap,
	}, nil
}

// HandleVote Handles a vote period
func (feeder *BlockFeeder) HandleVote(ctx context.Context) error {
	roundId := feeder.currentRound.Id
	if feeder.lastPreVote == nil {
		return nil
	}

	if feeder.lastVote != nil && roundId == feeder.lastVote.RoundId {
		feeder.logger.Debug("already sent a vote for this height, skipping this vote period...")
		return nil
	}

	if err := feeder.sendVote(ctx); err != nil {
		feeder.logger.Error(fmt.Sprintf("failed to send vote: %v", err))
		return fmt.Errorf("failed to send vote: %w", err)
	}

	feeder.lastVote = &BlockVoteInfo{
		RoundId:  roundId,
		Salt:     feeder.lastPreVote.Salt,
		VoteData: feeder.lastPreVote.VoteData,
	}

	return nil
}

// HandlePrevote Handles a prevote period
func (feeder *BlockFeeder) HandlePrevote(ctx context.Context) error {
	round := feeder.currentRound
	if feeder.lastPreVote != nil && round.Id == feeder.lastPreVote.RoundId {
		feeder.logger.Debug("already prevote, skipping this period...")
		return nil
	}

	voteData := types.VoteDataArr{}
	for _, od := range round.OracleData {
		switch od.Topic {
		case oracletypes.OralceTopic_Block:
			blockDataStr, err := feeder.gatherBlockDataString(od.Sources, uint64(round.Timestamp))
			if err != nil {
				feeder.abstainRoundId = round.Id
				return err
			}

			voteData = append(voteData, &oracletypes.VoteData{
				Topic: oracletypes.OralceTopic_Block,
				Data:  blockDataStr,
			})

		case oracletypes.OralceTopic_Ownership:
			nftDataStr, err := feeder.gatherNftOwnerDataString(od.Sources, uint64(round.Timestamp))
			if err != nil {
				feeder.abstainRoundId = round.Id
				return err
			}
			voteData = append(voteData, &oracletypes.VoteData{
				Topic: oracletypes.OralceTopic_Ownership,
				Data:  nftDataStr,
			})
		}
	}

	salt, err := GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	if err := feeder.sendPrevote(ctx, voteData, salt, round.Id); err != nil {
		return fmt.Errorf("failed to send prevote: %w", err)
	}

	feeder.lastPreVote = &BlockVoteInfo{
		RoundId:  round.Id,
		Salt:     salt,
		VoteData: voteData,
	}

	return nil
}

// HandleAbstain Handles a prevote period when block data string cannot be gathered
func (feeder *BlockFeeder) HandleAbstain(ctx context.Context) error {
	round := feeder.currentRound
	if feeder.lastPreVote != nil && round.Id == feeder.lastPreVote.RoundId {
		feeder.logger.Debug("already prevote, skipping this period...")
		return nil
	}

	feeder.logger.Info("failed to gather block data string, abstaining from this vote period")

	salt, err := GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// abstain from voting
	abstainData := []*oracletypes.VoteData{}
	if err := feeder.sendPrevote(ctx, abstainData, salt, round.Id); err != nil {
		return fmt.Errorf("failed to send abstain prevote: %w", err)
	}

	feeder.lastPreVote = &BlockVoteInfo{
		RoundId:  round.Id,
		Salt:     salt,
		VoteData: abstainData,
	}

	return nil
}

func (feeder *BlockFeeder) gatherBlockDataString(chainIds []string, timestamp uint64) ([]string, error) {
	blockDataList := make([]oracletypes.BlockData, len(chainIds))
	for idx, cid := range chainIds {
		bd, err := feeder.getOldestBlock(cid, timestamp)
		if err != nil {
			return nil, err
		}

		blockDataList[idx] = bd
	}

	return BlockDataListToBlockDataString(blockDataList), nil
}

func (feeder *BlockFeeder) getOldestBlock(chainId string, timestamp uint64) (bd oracletypes.BlockData, err error) {
	cc, ok := feeder.subscribers[chainId]
	if !ok {
		return bd, fmt.Errorf("chain client not found for chain id: %s", chainId)
	}

	return cc.GetOldestBlock(timestamp)
}

func (feeder *BlockFeeder) gatherNftOwnerDataString(nftIds []string, timestamp uint64) ([]string, error) {
	s := make([]string, len(nftIds))
	for i, nftId := range nftIds {
		nft, err := oracletypes.ParseNftId(nftId)
		if err != nil {
			return nil, err
		}

		bd, err := feeder.getOldestBlock(nft.ChainId, timestamp)
		if err != nil {
			return nil, err
		}

		cc, ok := feeder.subscribers[nft.ChainId]
		if !ok {
			return nil, fmt.Errorf("chain client not found for chain id: %s", nft.ChainId)
		}

		owner, err := cc.OwnerOf(context.TODO(), nft.ContractAddr.String(), nft.TokenId.String(), bd.BlockHash)
		if err != nil {
			return nil, err
		}

		s[i] = fmt.Sprintf("%s:%s", nftId, owner)
	}

	return s, nil
}

// sendVote sends a vote to the Settlus node
func (feeder *BlockFeeder) sendVote(ctx context.Context) error {
	if feeder.lastPreVote.Salt == "" {
		// we skip if salt is empty, which means no previous prevote was sent.
		feeder.logger.Info("salt or blockDataString is empty, skipping this vote period...")
		return nil
	}

	msg := oracletypes.NewMsgVote(
		feeder.address,
		feeder.validatorAddress,
		feeder.lastPreVote.VoteData,
		feeder.lastPreVote.Salt,
		feeder.lastPreVote.RoundId,
	)

	feeder.logger.Debug("try to send vote", "msg", msg.String())
	if err := feeder.sc.BuildAndSendTxWithRetry(ctx, msg); err != nil {
		return fmt.Errorf("failed to send vote tx: %w", err)
	}
	feeder.logger.Info("vote sent successfully", "msg", msg.String())

	return nil
}

// sendPrevote sends a prevote to the Settlus node
func (feeder *BlockFeeder) sendPrevote(ctx context.Context, vda types.VoteDataArr, salt string, roundId uint64) error {
	hash := GeneratePrevoteHash(vda, salt)

	msg := oracletypes.NewMsgPrevote(
		feeder.address,
		feeder.validatorAddress,
		hash,
		roundId,
	)

	feeder.logger.Debug("try to send prevote", "msg", msg.String(), "vote data", vda)
	if err := feeder.sc.BuildAndSendTxWithRetry(ctx, msg); err != nil {
		return fmt.Errorf("failed to send prevote tx: %w", err)
	}
	feeder.logger.Info("prevote sent successfully", "msg", msg.String(), "vote data", vda)

	return nil
}

// BlockDataListToBlockDataString converts a list of BlockData to a string
func BlockDataListToBlockDataString(bdList []oracletypes.BlockData) []string {
	s := make([]string, len(bdList))
	for i, bd := range bdList {
		s[i] = fmt.Sprintf("%s:%d/%s", bd.ChainId, bd.BlockNumber, bd.BlockHash)
	}
	return s
}
