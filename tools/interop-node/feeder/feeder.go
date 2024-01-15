package feeder

import (
	"context"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/settlus/chain/tools/interop-node/client"
	"github.com/settlus/chain/tools/interop-node/config"
	"github.com/settlus/chain/tools/interop-node/subscriber"
)

type Feeder interface {
	IsPreVotingPeriod(height int64) bool
	IsVotingPeriod(height int64) bool
	WantAbstain(height int64) bool

	HandlePrevote(ctx context.Context, height int64) error
	HandleVote(ctx context.Context, height int64) error
	HandleAbstain(ctx context.Context, height int64) error
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
