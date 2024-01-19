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
	Height          int64
	Salt            string
	BlockDataString string
}

type BlockFeeder struct {
	logger log.Logger

	abstainHeight int64
	lastPreVote   BlockVoteInfo
	lastVote      BlockVoteInfo

	chainClients []subscriber.Subscriber

	validatorAddress string
	address          string

	sc *client.SettlusClient
}

var _ Feeder = (*BlockFeeder)(nil)

func NewBlockFeeder(
	config *config.Config,
	sc *client.SettlusClient,
	chainclients []subscriber.Subscriber,
	logger log.Logger,
) (*BlockFeeder, error) {
	return &BlockFeeder{
		logger:           logger.With("topic", "block"),
		chainClients:     chainclients,
		validatorAddress: config.Feeder.ValidatorAddress,
		address:          config.Feeder.Address,
		sc:               sc,

		lastPreVote: BlockVoteInfo{
			Height: -1,
		},
		lastVote: BlockVoteInfo{
			Height: -1,
		},
	}, nil
}

// IsVotingPeriod returns true if the current height is a voting period
func (feeder *BlockFeeder) IsVotingPeriod(height int64) bool {
	// TODO: Settlus chain should provide voting period information
	return height%2 == 0
}

// IsPreVotingPeriod returns true if the current height is a prevoting period
func (feeder *BlockFeeder) IsPreVotingPeriod(height int64) bool {
	return height%2 == 1
}

// WantAbstain returns true if the feeder wants to abstain from voting
func (feeder *BlockFeeder) WantAbstain(height int64) bool {
	return feeder.abstainHeight == height
}

// HandleVote Handles a vote period
func (feeder *BlockFeeder) HandleVote(ctx context.Context, height int64) error {
	if height == feeder.lastVote.Height {
		feeder.logger.Debug("already sent a vote for this height, skipping this vote period...")
		return nil
	}

	if err := feeder.sendVote(ctx); err != nil {
		feeder.logger.Error(fmt.Sprintf("failed to send vote: %v", err))
		return fmt.Errorf("failed to send vote: %w", err)
	}

	feeder.lastVote = BlockVoteInfo{
		Height:          height,
		Salt:            feeder.lastPreVote.Salt,
		BlockDataString: feeder.lastPreVote.BlockDataString,
	}

	return nil
}

// HandlePrevote Handles a prevote period
func (feeder *BlockFeeder) HandlePrevote(ctx context.Context, height int64) error {
	blockDataStr, err := feeder.gatherBlockDataString()
	if err != nil {
		feeder.abstainHeight = height
		return err
	}

	if blockDataStr == feeder.lastPreVote.BlockDataString && height == feeder.lastPreVote.Height {
		feeder.logger.Debug("blockDataString is the same as the last used one, skipping this vote period...")
		return nil
	}

	salt, err := GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	if err := feeder.sendPrevote(ctx, blockDataStr, salt); err != nil {
		return fmt.Errorf("failed to send prevote: %w", err)
	}

	feeder.lastPreVote = BlockVoteInfo{
		Height:          height,
		Salt:            salt,
		BlockDataString: blockDataStr,
	}

	return nil
}

// HandleAbstain Handles a prevote period when block data string cannot be gathered
func (feeder *BlockFeeder) HandleAbstain(ctx context.Context, height int64) error {
	if height == feeder.lastPreVote.Height {
		feeder.logger.Debug("blockDataString is the same as the last used one, skipping this vote period...")
		return nil
	}

	feeder.logger.Info("failed to gather block data string, abstaining from this vote period")

	salt, err := GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// abstain from voting
	abstainString := GenerateAbstainString(feeder.chainClients)
	if err := feeder.sendPrevote(ctx, abstainString, salt); err != nil {
		return fmt.Errorf("failed to send abstain prevote: %w", err)
	}

	feeder.lastPreVote = BlockVoteInfo{
		Height:          height,
		Salt:            salt,
		BlockDataString: abstainString,
	}

	return nil
}

func (feeder *BlockFeeder) gatherBlockDataString() (string, error) {
	blockDataList := make([]oracletypes.BlockData, len(feeder.chainClients))
	for i, cc := range feeder.chainClients {
		bd, err := cc.GetBlockData()
		if err != nil {
			return "", err
		}
		blockDataList[i] = bd
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
	)

	feeder.logger.Debug("try to send vote", "msg", msg.String())
	if err := feeder.sc.BuildAndSendTxWithRetry(ctx, msg); err != nil {
		return fmt.Errorf("failed to send vote tx: %w", err)
	}
	feeder.logger.Info("vote sent successfully", "msg", msg.String())

	return nil
}

// sendPrevote sends a prevote to the Settlus node
func (feeder *BlockFeeder) sendPrevote(ctx context.Context, bd, salt string) error {
	hash := GeneratePrevoteHash(bd, salt)

	msg := oracletypes.NewMsgPrevote(
		feeder.address,
		feeder.validatorAddress,
		hash,
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
func GenerateAbstainString(subs []subscriber.Subscriber) string {
	var builder strings.Builder
	for _, sub := range subs {
		builder.WriteString(fmt.Sprintf("%d:%d:%s", sub.Id(), -1, ""))
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
