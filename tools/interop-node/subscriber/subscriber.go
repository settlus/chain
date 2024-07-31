package subscriber

import (
	"context"
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/settlus/chain/tools/interop-node/config"
	"github.com/settlus/chain/tools/interop-node/types"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

type Subscriber interface {
	Id() string
	Start(ctx context.Context)
	Stop()
	GetOldestBlock(timestamp uint64) (oracletypes.BlockData, error)
	OwnerOf(ctx context.Context, nftAddressHex string, tokenIdHex string, blockHash string) (string, error)
}

func InitSubscribers(config *config.Config, logger log.Logger) (subscribers []Subscriber, err error) {
	subscribers = make([]Subscriber, len(config.Chains))
	for i, chain := range config.Chains {
		subscriber, err := NewSubscriber(chain, logger)
		if err != nil {
			return nil, err
		}
		subscribers[i] = subscriber
	}

	return subscribers, err
}

func NewSubscriber(config config.ChainConfig, logger log.Logger) (Subscriber, error) {
	if config.ChainType == types.CHAINTYPE_ETHEREUM {
		return NewEthereumSubscriber(config.ChainID, config.RpcUrl, logger)
	}

	return nil, fmt.Errorf("unsupported chain type: %s", config.ChainType)
}
