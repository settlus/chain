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

type Feeder struct {
	logger           log.Logger
	validatorAddress string
	address          string
	currentRound     oracletypes.RoundInfo
	abstainRoundId   uint64
	sc               *client.SettlusClient

	lastPreVote *BlockVoteInfo
	lastVote    *BlockVoteInfo
	subscribers map[string]subscriber.Subscriber
}

func (feeder *Feeder) FetchNewRoundInfo(ctx context.Context) {
	roundInfo, err := feeder.sc.FetchNewRoundInfo(ctx)
	if err != nil {
		feeder.logger.Error(fmt.Sprintf("failed to fetch round info: %v", err))
		return
	}

	feeder.currentRound = *roundInfo
}

// IsVotingPeriod returns true if the current height is a voting period
func (feeder *Feeder) IsVotingPeriod(height int64) bool {
	return height > int64(feeder.currentRound.PrevoteEnd) && height <= int64(feeder.currentRound.VoteEnd)
}

// IsPreVotingPeriod returns true if the current height is a prevoting period
func (feeder *Feeder) IsPreVotingPeriod(height int64) bool {
	return height <= int64(feeder.currentRound.PrevoteEnd)
}

// WantAbstain returns true if the feeder wants to abstain from voting
func (feeder *Feeder) WantAbstain(height int64) bool {
	return feeder.IsPreVotingPeriod(height) && feeder.currentRound.Id == feeder.abstainRoundId
}

type BlockVoteInfo struct {
	RoundId  uint64
	Salt     string
	VoteData types.VoteDataArr
}

func NewFeeder(
	config *config.Config,
	sc *client.SettlusClient,
	subscribers []subscriber.Subscriber,
	logger log.Logger,
) (*Feeder, error) {
	subscribersMap := make(map[string]subscriber.Subscriber)
	for _, cc := range subscribers {
		subscribersMap[cc.Id()] = cc
	}

	return &Feeder{
		logger:           logger.With("topic", "block"),
		validatorAddress: config.Feeder.ValidatorAddress,
		address:          config.Feeder.Address,
		sc:               sc,
		subscribers:      subscribersMap,
	}, nil
}

// HandleVote Handles a vote period
func (feeder *Feeder) HandleVote(ctx context.Context) error {
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
func (feeder *Feeder) HandlePrevote(ctx context.Context) error {
	round := feeder.currentRound
	if feeder.lastPreVote != nil && round.Id == feeder.lastPreVote.RoundId {
		feeder.logger.Debug("already prevote, skipping this period...")
		return nil
	}

	voteData := types.VoteDataArr{}
	for _, od := range round.OracleData {
		switch od.Topic {
		case oracletypes.OralceTopic_BLOCK:
			blockDataStr, err := feeder.gatherBlockDataString(od.Sources, uint64(round.Timestamp))
			if err != nil {
				feeder.abstainRoundId = round.Id
				return err
			}

			voteData = append(voteData, &oracletypes.VoteData{
				Topic: oracletypes.OralceTopic_BLOCK,
				Data:  blockDataStr,
			})

		case oracletypes.OralceTopic_OWNERSHIP:
			nftDataStr, err := feeder.gatherNftOwnerDataString(od.Sources, uint64(round.Timestamp))
			if err != nil {
				feeder.abstainRoundId = round.Id
				return err
			}
			voteData = append(voteData, &oracletypes.VoteData{
				Topic: oracletypes.OralceTopic_OWNERSHIP,
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
func (feeder *Feeder) HandleAbstain(ctx context.Context) error {
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

func (feeder *Feeder) gatherBlockDataString(chainIds []string, timestamp uint64) ([]string, error) {
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

func (feeder *Feeder) getOldestBlock(chainId string, timestamp uint64) (bd oracletypes.BlockData, err error) {
	cc, ok := feeder.subscribers[chainId]
	if !ok {
		return bd, fmt.Errorf("chain client not found for chain id: %s", chainId)
	}

	return cc.GetOldestBlock(timestamp)
}

func (feeder *Feeder) gatherNftOwnerDataString(nftIds []string, timestamp uint64) ([]string, error) {
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
func (feeder *Feeder) sendVote(ctx context.Context) error {
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
func (feeder *Feeder) sendPrevote(ctx context.Context, vda types.VoteDataArr, salt string, roundId uint64) error {
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
