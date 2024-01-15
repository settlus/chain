package subscriber

import (
	"context"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/settlus/chain/tools/interop-node/config"
	"github.com/settlus/chain/tools/interop-node/repository"
	"github.com/settlus/chain/tools/interop-node/types"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

type Subscriber interface {
	Id() uint64
	Start(ctx context.Context)
	Stop()
	GetBlockData() (oracletypes.BlockData, error)
	OwnerOf(ctx context.Context, nftAddressHex string, tokenIdHex string, blockHash string) (string, error)
}

func InitSubscribers(config *config.Config, logger log.Logger) (subscribers []Subscriber, err error) {
	subscribers = make([]Subscriber, len(config.Chains))
	for i, chain := range config.Chains {
		repo := repository.NewLevelDbRepostiory(config.DBHome, chain.ChainID)
		subscriber, err := NewSubscriber(chain, logger, repo)
		if err != nil {
			return nil, err
		}
		subscribers[i] = subscriber
	}

	return subscribers, err
}

func NewSubscriber(config config.ChainConfig, logger log.Logger, repo repository.Repository) (Subscriber, error) {
	if config.ChainType == types.CHAINTYPE_ETHEREUM {
		return NewEthereumSubscriber(config.ChainID, config.RpcUrl, logger, repo)
	}

	return nil, fmt.Errorf("unsupported chain type: %s", config.ChainType)
}
