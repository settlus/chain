package subscriber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/settlus/chain/tools/interop-node/types"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

type EthereumSubscriber struct {
	chainId string
	cache   *BlockCache
	client  *http.Client
	rpcUrl  string
	logger  log.Logger
}

var _ Subscriber = (*EthereumSubscriber)(nil)

func NewEthereumSubscriber(chainId string, rpcUrl string, logger log.Logger) (*EthereumSubscriber, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &EthereumSubscriber{
		chainId: chainId,
		cache:   NewBlockCache(100),
		client:  client,
		rpcUrl:  rpcUrl,
		logger:  logger,
	}, nil
}

func (sub *EthereumSubscriber) Id() string {
	return sub.chainId
}

func (sub *EthereumSubscriber) Start(ctx context.Context) {
	go sub.fetchLoop(ctx)
}

func (sub *EthereumSubscriber) Stop() {
	sub.client.CloseIdleConnections()
}

// OwnerOf returns the owner of the given NFT
func (sub *EthereumSubscriber) OwnerOf(_ context.Context, nftAddressHex string, tokenIdHex string, blockHash string) (string, error) {
	return sub.findOwnerFromNetwork(nftAddressHex, tokenIdHex, blockHash)
}

func (sub *EthereumSubscriber) findOwnerFromNetwork(nftAddressHex string, tokenIdHex string, blockHash string) (string, error) {
	if !strings.HasPrefix(nftAddressHex, "0x") {
		nftAddressHex = "0x" + nftAddressHex
	}

	if !strings.HasPrefix(blockHash, "0x") {
		blockHash = "0x" + blockHash
	}

	data := fmt.Sprintf("0x6352211e%064x", common.HexToHash(tokenIdHex))

	type CallMsg struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}

	type BlockHash struct {
		BlockHash string `json:"blockHash"`
	}

	body, err := json.Marshal(types.EthRpcInput{
		JsonRpc: "2.0",
		Method:  "eth_call",
		Params: []interface{}{CallMsg{
			To:   nftAddressHex,
			Data: data,
		}, BlockHash{
			BlockHash: blockHash,
		}},
		Id: "1",
	})
	if err != nil {
		return "failed to marshal json", err
	}

	res, err := sub.client.Post(sub.rpcUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "failed to get ownerOf http response", err
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return "failed to read response body", err
	}

	var owner types.OwnerOfOutput
	if err := json.Unmarshal(resp, &owner); err != nil {
		return "failed to unmarshal json", err
	}

	if owner.Error != nil {
		if owner.Error.Code == 3 { //Tranasction Reverted
			return "0x00", nil
		}

		return "", fmt.Errorf("failed to get ownerOf: %s code: %d", owner.Error.Message, owner.Error.Code)
	}

	return owner.Result, nil
}

const fetchInterval = 5 * time.Second

func (sub *EthereumSubscriber) fetchLoop(ctx context.Context) {
	ticker := time.NewTicker(fetchInterval)

	for {
		select {
		case <-ticker.C:
			blockNumber, err := sub.blockNumber()
			if err != nil {
				sub.logger.Error(err.Error(), "from", "latest block")
				continue
			}
			blockHash, blockTime, err := sub.getBlockByNumber(blockNumber)
			if err != nil {
				sub.logger.Error(err.Error(), "from", "block fetch", "hash", blockHash)
				continue
			}

			sub.cache.PutBlockData(blockHash, blockNumber.Int64(), blockTime)
		case <-ctx.Done():
			sub.logger.Info("fetchLoop stopped")
			return
		}
	}
}

// GetOldestBlock returns the latest block data
func (sub *EthereumSubscriber) GetOldestBlock(timestamp uint64) (oracletypes.BlockData, error) {
	hash, number := sub.cache.GetOldestBlock(timestamp)
	if hash == "" || number == 0 {
		return oracletypes.BlockData{}, fmt.Errorf("failed to get oldest block")
	}

	return oracletypes.BlockData{
		ChainId:     sub.chainId,
		BlockNumber: number,
		BlockHash:   hash,
	}, nil
}

// blockNumber returns the latest block number
func (sub *EthereumSubscriber) blockNumber() (*big.Int, error) {
	res, err := sub.client.Post(sub.rpcUrl, "application/json", bytes.NewBuffer([]byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)))
	if err != nil {
		return nil, fmt.Errorf("failed to get block number: %w", err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var blockNumber types.BlockNumberOutput
	if err := json.Unmarshal(respBody, &blockNumber); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return common.HexToHash(blockNumber.Result).Big(), nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	finalized := big.NewInt(int64(rpc.FinalizedBlockNumber))
	if number.Cmp(finalized) == 0 {
		return "finalized"
	}
	safe := big.NewInt(int64(rpc.SafeBlockNumber))
	if number.Cmp(safe) == 0 {
		return "safe"
	}
	return hexutil.EncodeBig(number)
}

// getBlockByNumber returns the block by number
func (sub *EthereumSubscriber) getBlockByNumber(blockNumber *big.Int) (string, uint64, error) {
	body, err := json.Marshal(types.EthRpcInput{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{toBlockNumArg(blockNumber), false},
		Id:      "1",
	})
	if err != nil {
		return "", 0, fmt.Errorf("failed to marshal json: %w", err)
	}

	res, err := sub.client.Post(sub.rpcUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", 0, fmt.Errorf("failed to get block: %w", err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var block types.BlockByNumberOutput
	if err := json.Unmarshal(respBody, &block); err != nil {
		return "", 0, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	if block.Error != nil {
		return "", 0, fmt.Errorf("failed to get block: %s", block.Error.Message)
	}

	timestamp, err := hexutil.DecodeUint64(block.Result.Timestamp)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decode timestamp: %w", err)
	}

	return block.Result.Hash, timestamp * 1000, nil
}
