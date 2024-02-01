package subscriber

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/settlus/chain/tools/interop-node/repository"
	"github.com/settlus/chain/tools/interop-node/types"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

type EthereumSubscriber struct {
	chainId   *big.Int
	repo      repository.Repository
	client    *ethclient.Client
	lastBlock oracletypes.BlockData
	logger    log.Logger
	dbCh      chan *types.BlockEventData
}

var _ Subscriber = (*EthereumSubscriber)(nil)

func NewEthereumSubscriber(chainId string, rpcUrl string, logger log.Logger, repo repository.Repository) (*EthereumSubscriber, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	chainIdInt, ok := math.ParseBig256(chainId)
	if !ok {
		return nil, fmt.Errorf("failed to parse chainId: %v", err)
	}

	dbCh := make(chan *types.BlockEventData)

	return &EthereumSubscriber{
		chainId: chainIdInt,
		repo:    repo,
		client:  client,
		logger:  logger,
		dbCh:    dbCh,
	}, nil
}

func (sub *EthereumSubscriber) Id() uint64 {
	return sub.chainId.Uint64()
}

func (sub *EthereumSubscriber) Start(ctx context.Context) {
	go sub.dbLoop(ctx)
	go sub.fetchLoop(ctx)
	go sub.fillDB()
}

func (sub *EthereumSubscriber) Stop() {
	sub.client.Close()
}

func (sub *EthereumSubscriber) OwnerOf(ctx context.Context, nftAddressHex string, tokenIdHex string, blockHash string) (string, error) {
	ownerHex, err := sub.findOwnerFromDb(nftAddressHex, tokenIdHex, blockHash)

	switch err.(type) {
	case *types.NotFoundError:
		sub.logger.Info("Not found in db, try to find from network", "error", err)
		ownerHex, err := sub.findOwnerFromNetwork(ctx, nftAddressHex, tokenIdHex, blockHash)
		if err == nil {
			// lazily register ownership data to db
			go func() {
				block, err := sub.client.BlockByHash(context.Background(), common.HexToHash(blockHash))
				if err != nil {
					sub.logger.Error("Failed to get block", "blockHash", blockHash, "error", err)
					return
				}
				events := []*types.OwnershipTransferEvent{{
					ContractAddr: common.HexToAddress(nftAddressHex).Bytes(),
					To:           common.HexToAddress(ownerHex).Bytes(),
					TokenId:      common.HexToHash(tokenIdHex).Big(),
				}}

				if err := sub.repo.PutBlockData(block.Hash().Bytes(), block.Number().Bytes(), events); err != nil {
					sub.logger.Error("Failed to putBlockData", "error", err)
				} else {
					sub.logger.Info("put new block", "number", block.Number(), "hash", blockHash)
				}
			}()
			return ownerHex, nil
		}
	default:
		return ownerHex, err
	}

	return ownerHex, err
}

// GetBlockData returns the latest block data
func (sub *EthereumSubscriber) GetBlockData() (oracletypes.BlockData, error) {
	if sub.lastBlock.BlockHash == "" {
		return sub.lastBlock, errors.New("no block data available")
	}

	return sub.lastBlock, nil
}

func (sub *EthereumSubscriber) findOwnerFromNetwork(ctx context.Context, nftAddressHex string, tokenIdHex string, blockHash string) (string, error) {
	addr := common.HexToAddress(nftAddressHex)
	data := common.Hex2Bytes(fmt.Sprintf("6352211e%064x", common.HexToHash(tokenIdHex))) // ownerOf = 0x6352211e
	msg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}

	res, err := sub.client.CallContractAtHash(ctx, msg, common.HexToHash(blockHash))

	return common.BytesToAddress(res).Hex(), err
}

func (sub *EthereumSubscriber) findOwnerFromDb(nftAddressHex string, tokenIdHex string, blockHash string) (string, error) {
	blockNumberHex, err := sub.repo.GetBlockNumber(blockHash)
	if err != nil {
		return "", types.NewNotFoundError(fmt.Sprintf("block(%s) data not found", blockHash))
	}

	return sub.repo.GetNftOwnership(nftAddressHex, tokenIdHex, blockNumberHex)
}

func (sub *EthereumSubscriber) dbLoop(ctx context.Context) {
	for {
		select {
		case event := <-sub.dbCh:
			// if channel closed
			if event == nil {
				return
			}

			// check whether block data already exists
			if hash, err := sub.repo.GetBlockHash(common.Bytes2Hex(event.BlockNumber.Bytes())); err == nil {
				if hash != "" && hash != common.Bytes2Hex(event.BlockHash) {
					sub.logger.Error("BlockHash mismatch", "db", hash, "new", common.Bytes2Hex(event.BlockHash))
				}
				continue
			}

			if err := sub.repo.PutBlockData(event.BlockHash, event.BlockNumber.Bytes(), event.NftTransferred); err != nil {
				sub.logger.Error("Failed to putBlockData", "error", err)
				continue
			}

			sub.lastBlock = oracletypes.BlockData{
				ChainId:     sub.chainId.String(),
				BlockNumber: event.BlockNumber.Int64(),
				BlockHash:   common.Bytes2Hex(event.BlockHash),
			}

			sub.logger.Info("put new block", "number", event.BlockNumber.Int64(), "hash", common.Bytes2Hex(event.BlockHash))

		case <-ctx.Done():
			close(sub.dbCh)
		}
	}
}

const fetchInterval = 3 * time.Second

func (sub *EthereumSubscriber) fetchLoop(ctx context.Context) {
	ticker := time.NewTicker(fetchInterval)

	for {
		select {
		case <-ticker.C:
			blockNumber, err := sub.client.BlockNumber(ctx)
			if err != nil {
				sub.logger.Error(err.Error(), "from", "latest block")
				continue
			}
			block, err := sub.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
			if err != nil {
				sub.logger.Error(err.Error(), "from", "block fetch", "number", blockNumber)
				continue
			}
			event, err := sub.parseBlock(block)
			if err != nil {
				sub.logger.Error(err.Error(), "from", "parse block", "number", blockNumber)
				continue
			}

			sub.dbCh <- event
		case <-ctx.Done():
			sub.logger.Info("fetchLoop stopped")
			close(sub.dbCh)
			return
		}
	}

}

func (sub *EthereumSubscriber) parseBlock(block *ethtypes.Block) (*types.BlockEventData, error) {
	erc721Transferred, err := sub.getTransferEventsFromBlock(block)
	if err != nil {
		return nil, err
	}

	return &types.BlockEventData{
		BlockNumber:    block.Number(),
		BlockHash:      block.Hash().Bytes(),
		NftTransferred: erc721Transferred,
	}, nil
}

func (sub *EthereumSubscriber) getTransferEventsFromBlock(block *ethtypes.Block) ([]*types.OwnershipTransferEvent, error) {
	logs, err := sub.client.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: block.Number(),
		ToBlock:   block.Number(),
		Topics:    [][]common.Hash{{common.HexToHash(types.EventTransferSignature)}},
	})
	if err != nil {
		return nil, err
	}

	erc721events := make([]*types.OwnershipTransferEvent, 0)
	for _, vLog := range logs {
		if len(vLog.Topics) != 4 || vLog.Topics[0].Hex() != types.EventTransferSignature {
			continue
		}

		erc721events = append(erc721events, &types.OwnershipTransferEvent{
			TxId:         vLog.TxHash.Hex(),
			ContractAddr: vLog.Address.Bytes(),
			To:           vLog.Topics[2].Bytes(),
			TokenId:      vLog.Topics[3].Big(),
			BlockNumber:  vLog.BlockNumber,
		})
	}

	return erc721events, nil
}
