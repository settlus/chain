package subscriber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/common"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	oracletypes "github.com/settlus/chain/x/oracle/types"

	"github.com/settlus/chain/tools/interop-node/repository"
	"github.com/settlus/chain/tools/interop-node/types"
)

type EthereumSubscriber struct {
	chainId string
	repo    repository.Repository
	client  *http.Client
	rpcUrl  string
	logger  log.Logger
	dbCh    chan *types.BlockEventData
}

var _ Subscriber = (*EthereumSubscriber)(nil)

func NewEthereumSubscriber(chainId string, rpcUrl string, logger log.Logger, repo repository.Repository) (*EthereumSubscriber, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	dbCh := make(chan *types.BlockEventData)

	return &EthereumSubscriber{
		chainId: chainId,
		repo:    repo,
		client:  client,
		rpcUrl:  rpcUrl,
		logger:  logger,
		dbCh:    dbCh,
	}, nil
}

func (sub *EthereumSubscriber) Id() string {
	return sub.chainId
}

func (sub *EthereumSubscriber) Start(ctx context.Context) {
	go sub.dbLoop(ctx)
	go sub.fetchLoop(ctx)
	go sub.fillDB()
}

func (sub *EthereumSubscriber) Stop() {
	sub.client.CloseIdleConnections()
}

// OwnerOf returns the owner of the given NFT
func (sub *EthereumSubscriber) OwnerOf(_ context.Context, nftAddressHex string, tokenIdHex string, blockHash string) (string, error) {
	ownerHex, err := sub.findOwnerFromDb(nftAddressHex, tokenIdHex, blockHash)

	switch err.(type) {
	case *types.NotFoundError:
		sub.logger.Info("Not found in db, try to find from network", "error", err)
		ownerHex, err := sub.findOwnerFromNetwork(nftAddressHex, tokenIdHex, blockHash)
		if err == nil {
			// lazily register ownership data to db
			go func() {
				block, err := sub.blockByHash(blockHash)
				if err != nil {
					sub.logger.Error("Failed to get block", "error", err)
					return
				}

				events := []*types.OwnershipTransferEvent{{
					ContractAddr: common.HexToAddress(nftAddressHex).Bytes(),
					To:           common.HexToAddress(ownerHex).Bytes(),
					TokenId:      common.HexToHash(tokenIdHex).Big(),
				}}

				if err := sub.repo.PutBlockData(block.Hash().Bytes(), block.Number().Bytes(), block.Header().Time*1000, events); err != nil {
					sub.logger.Error("Failed to putBlockData", "error", err)
				} else {
					sub.logger.Info("put new block", "number", block.Number(), "hash", blockHash, "timestamp", block.Header().Time*1000)
				}
			}()
			return ownerHex, nil
		}
	default:
		return ownerHex, err
	}

	return ownerHex, err
}

// GetOldestBlock returns the latest block data
func (sub *EthereumSubscriber) GetOldestBlock(timestamp uint64) (oracletypes.BlockData, error) {
	block, err := sub.repo.GetOldestBlock(timestamp)
	if err != nil {
		return oracletypes.BlockData{}, err
	}

	return oracletypes.BlockData{
		ChainId:     sub.chainId,
		BlockNumber: common.BytesToHash(block.Number).Big().Int64(),
		BlockHash:   common.Bytes2Hex(block.Hash),
	}, nil
}

type EthRpcInput struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
}

type OwnerOfOutput struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func (sub *EthereumSubscriber) findOwnerFromNetwork(nftAddressHex string, tokenIdHex string, blockHash string) (string, error) {
	addr := common.HexToAddress(nftAddressHex)
	data := common.Hex2Bytes(fmt.Sprintf("6352211e%064x", common.HexToHash(tokenIdHex))) // ownerOf = 0x6352211e
	msg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}

	body, err := json.Marshal(EthRpcInput{
		JsonRpc: "2.0",
		Method:  "eth_call",
		Params:  []interface{}{msg, blockHash},
		Id:      "1",
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

	var owner OwnerOfOutput
	if err := json.Unmarshal(resp, &owner); err != nil {
		return "failed to unmarshal json", err
	}

	return owner.Result, nil
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

			if err := sub.repo.PutBlockData(event.BlockHash, event.BlockNumber.Bytes(), event.Timestamp, event.NftTransferred); err != nil {
				sub.logger.Error("Failed to putBlockData", "error", err)
				continue
			}

			sub.logger.Info("put new block", "number", event.BlockNumber.Int64(), "hash", common.Bytes2Hex(event.BlockHash), "timestamp", event.Timestamp)

		case <-ctx.Done():
			close(sub.dbCh)
			return
		}
	}
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
			block, err := sub.blockByNumber(blockNumber)
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
			return
		}
	}

}

// parseBlock parses the block and returns the block event data
func (sub *EthereumSubscriber) parseBlock(block *ethtypes.Header) (*types.BlockEventData, error) {
	erc721Transferred, err := sub.getTransferEventsFromBlock(block)
	if err != nil {
		return nil, err
	}

	return &types.BlockEventData{
		BlockNumber:    block.Number,
		BlockHash:      block.Hash().Bytes(),
		Timestamp:      block.Time * 1000,
		NftTransferred: erc721Transferred,
	}, nil
}

// toFilterArg converts ethereum.FilterQuery to a filter argument for eth_getFilterLogs JSON RPC
func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = hexutil.EncodeBig(q.FromBlock)
		}
		arg["toBlock"] = hexutil.EncodeBig(q.ToBlock)
	}
	return arg, nil
}

type FilterLogOutput struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      string         `json:"id"`
	Result  []ethtypes.Log `json:"result"`
}

// getTransferEventsFromBlock returns all transfer events from a block
func (sub *EthereumSubscriber) getTransferEventsFromBlock(block *ethtypes.Header) ([]*types.OwnershipTransferEvent, error) {
	filterArgs, err := toFilterArg(ethereum.FilterQuery{
		FromBlock: block.Number,
		ToBlock:   block.Number,
		Topics:    [][]common.Hash{{common.HexToHash(types.EventTransferSignature)}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create filter args: %w", err)
	}

	body, err := json.Marshal(EthRpcInput{
		JsonRpc: "2.0",
		Method:  "eth_getFilterLogs",
		Params:  []interface{}{filterArgs},
		Id:      "1",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	res, err := sub.client.Post(sub.rpcUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to get filter logs: %w", err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var logs FilterLogOutput

	if err := json.Unmarshal(respBody, &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	erc721events := make([]*types.OwnershipTransferEvent, 0)
	for _, vLog := range logs.Result {
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

type BlockNumberOutput struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
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

	var blockNumber BlockNumberOutput
	if err := json.Unmarshal(respBody, &blockNumber); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return common.HexToHash(blockNumber.Result).Big(), nil
}

type BlockByHashOutput struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Result  *ethtypes.Block `json:"result"`
}

// blockByHash returns the block by hash
func (sub *EthereumSubscriber) blockByHash(blockHash string) (*ethtypes.Block, error) {
	body, err := json.Marshal(EthRpcInput{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByHash",
		Params:  []interface{}{common.HexToHash(blockHash), false},
		Id:      "1",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	res, err := sub.client.Post(sub.rpcUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var block BlockByHashOutput
	if err := json.Unmarshal(respBody, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return block.Result, nil
}

type BlockByNumberOutput struct {
	Jsonrpc string           `json:"jsonrpc"`
	Id      string           `json:"id"`
	Result  *ethtypes.Header `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
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

// blockByNumber returns the block by number
func (sub *EthereumSubscriber) blockByNumber(blockNumber *big.Int) (*ethtypes.Header, error) {
	body, err := json.Marshal(EthRpcInput{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{toBlockNumArg(blockNumber), false},
		Id:      "1",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	res, err := sub.client.Post(sub.rpcUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var block BlockByNumberOutput
	if err := json.Unmarshal(respBody, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	if block.Error != nil {
		return nil, fmt.Errorf("failed to get block: %s", block.Error.Message)
	}

	return block.Result, nil
}
