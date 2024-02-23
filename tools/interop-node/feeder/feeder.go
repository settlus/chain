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

type Feeder interface {
	IsPreVotingPeriod(height int64) bool
	IsVotingPeriod(height int64) bool
	WantAbstain(height int64) bool

	HandlePrevote(ctx context.Context) error
	HandleVote(ctx context.Context) error
	HandleAbstain(ctx context.Context) error
	FetchNewRoundInfo(ctx context.Context)
}

func InitFeeders(config *config.Config, sc *client.SettlusClient, chainClients []subscriber.Subscriber, logger log.Logger) ([]Feeder, error) {
	var feeders []Feeder

	feederTypes := strings.Split(config.Feeder.Topics, ",")
	for _, feederType := range feederTypes {
		feeder, err := NewFeeder(feederType, config, sc, chainClients, logger)
		if err != nil {
			return nil, err
		}

		if feeder == nil {
			return nil, fmt.Errorf("failed to create feeder without error: %s", feederType)
		}

		feeders = append(feeders, feeder)
	}

	return feeders, nil
}

func NewFeeder(feederType string, config *config.Config, sc *client.SettlusClient, chainClients []subscriber.Subscriber, logger log.Logger) (Feeder, error) {
	switch feederType {
	case BlockTopic:
		return NewBlockFeeder(config, sc, chainClients, logger)
	default:
		return nil, fmt.Errorf("unsupported feeder type: %s", feederType)
	}
}

type BaseFeeder struct {
	logger           log.Logger
	validatorAddress string
	address          string
	currentRound     oracletypes.RoundInfo
	abstainRoundId   uint64
	sc               *client.SettlusClient
}

func (feeder *BaseFeeder) FetchNewRoundInfo(ctx context.Context) {
	roundInfo, err := feeder.sc.FetchNewRoundInfo(ctx)
	if err != nil {
		feeder.logger.Error(fmt.Sprintf("failed to fetch round info: %v", err))
		return
	}

	feeder.currentRound = *roundInfo
}

// IsVotingPeriod returns true if the current height is a voting period
func (feeder *BaseFeeder) IsVotingPeriod(height int64) bool {
	return height > int64(feeder.currentRound.PrevoteEnd) && height <= int64(feeder.currentRound.VoteEnd)
}

// IsPreVotingPeriod returns true if the current height is a prevoting period
func (feeder *BaseFeeder) IsPreVotingPeriod(height int64) bool {
	return height <= int64(feeder.currentRound.PrevoteEnd)
}

// WantAbstain returns true if the feeder wants to abstain from voting
func (feeder *BaseFeeder) WantAbstain(height int64) bool {
	return feeder.IsPreVotingPeriod(height) && feeder.currentRound.Id == feeder.abstainRoundId
}
