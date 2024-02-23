package feeder

import (
	"context"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/settlus/chain/tools/interop-node/client"
	"github.com/settlus/chain/tools/interop-node/config"
	"github.com/settlus/chain/tools/interop-node/subscriber"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

const BlockTopic = "block"

type BlockVoteInfo struct {
	RoundId         uint64
	Salt            string
	BlockDataString string
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
		RoundId:         roundId,
		Salt:            feeder.lastPreVote.Salt,
		BlockDataString: feeder.lastPreVote.BlockDataString,
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

	blockDataStr, err := feeder.gatherBlockDataString(round)
	if err != nil {
		feeder.abstainRoundId = round.Id
		return err
	}

	salt, err := GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	if err := feeder.sendPrevote(ctx, blockDataStr, salt, round.Id); err != nil {
		return fmt.Errorf("failed to send prevote: %w", err)
	}

	feeder.lastPreVote = &BlockVoteInfo{
		RoundId:         round.Id,
		Salt:            salt,
		BlockDataString: blockDataStr,
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
	abstainString := GenerateAbstainString(round.ChainIds)
	if err := feeder.sendPrevote(ctx, abstainString, salt, round.Id); err != nil {
		return fmt.Errorf("failed to send abstain prevote: %w", err)
	}

	feeder.lastPreVote = &BlockVoteInfo{
		RoundId:         round.Id,
		Salt:            salt,
		BlockDataString: abstainString,
	}

	return nil
}

func (feeder *BlockFeeder) gatherBlockDataString(round oracletypes.RoundInfo) (string, error) {
	blockDataList := make([]oracletypes.BlockData, len(round.ChainIds))
	for idx, cid := range round.ChainIds {
		cc, ok := feeder.subscribers[cid]
		if !ok {
			return "", fmt.Errorf("chain client not found for chain id: %s", cc)
		}
		bd, err := cc.GetRecentBlockData(uint64(round.Timestamp))
		if err != nil {
			return "", err
		}

		blockDataList[idx] = bd
	}

	return BlockDataListToBlockDataString(blockDataList), nil
}

// sendVote sends a vote to the Settlus node
func (feeder *BlockFeeder) sendVote(ctx context.Context) error {
	if feeder.lastPreVote.Salt == "" || feeder.lastPreVote.BlockDataString == "" {
		// we skip if salt is empty, which means no previous prevote was sent.
		feeder.logger.Info("salt or blockDataString is empty, skipping this vote period...")
		return nil
	}

	msg := oracletypes.NewMsgVote(
		feeder.address,
		feeder.validatorAddress,
		feeder.lastPreVote.BlockDataString,
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
func (feeder *BlockFeeder) sendPrevote(ctx context.Context, bd, salt string, roundId uint64) error {
	hash := GeneratePrevoteHash(bd, salt)

	msg := oracletypes.NewMsgPrevote(
		feeder.address,
		feeder.validatorAddress,
		hash,
		roundId,
	)

	feeder.logger.Debug("try to send prevote", "msg", msg.String(), "blockDataString", bd)
	if err := feeder.sc.BuildAndSendTxWithRetry(ctx, msg); err != nil {
		return fmt.Errorf("failed to send prevote tx: %w", err)
	}
	feeder.logger.Info("prevote sent successfully", "msg", msg.String(), "blockDataString", bd)

	return nil
}

// GenerateAbstainString generates an abstain string from a list of chains
// To abstain from voting, we send a prevote with a negative blocknumber -1
func GenerateAbstainString(chainIds []string) string {
	var builder strings.Builder
	for _, chainId := range chainIds {
		builder.WriteString(fmt.Sprintf("%s:%d:%s", chainId, -1, ""))
	}
	return builder.String()
}

// BlockDataListToBlockDataString converts a list of BlockData to a string
func BlockDataListToBlockDataString(bdList []oracletypes.BlockData) string {
	var builder strings.Builder
	for i, bd := range bdList {
		builder.WriteString(fmt.Sprintf("%s:%d:%s", bd.ChainId, bd.BlockNumber, bd.BlockHash))
		if i != len(bdList)-1 {
			builder.WriteString(",")
		}
	}
	return builder.String()
}
