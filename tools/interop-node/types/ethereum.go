package types

import ethtypes "github.com/ethereum/go-ethereum/core/types"

// Header represents a block header in the Ethereum blockchain.
type Header struct {
	Hash      string `json:"hash"`
	Number    string `json:"number"`
	Timestamp string `json:"timestamp"`
}

// BlockByNumberOutput is a JSON-RPC response containing a block header.
// if an error occurred, the error field will be non-nil.
type BlockByNumberOutput struct {
	Jsonrpc string  `json:"jsonrpc"`
	Id      string  `json:"id"`
	Result  *Header `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// BlockByHashOutput is a JSON-RPC response containing a block header.
type BlockByHashOutput struct {
	Jsonrpc string  `json:"jsonrpc"`
	Id      int     `json:"id"`
	Result  *Header `json:"result"`
}

// BlockNumberOutput is a JSON-RPC response containing the current block number.
type BlockNumberOutput struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

type EthRpcInput struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
}

type OwnerOfOutput struct {
	JsonRpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  string `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type FilterLogOutput struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      string         `json:"id"`
	Result  []ethtypes.Log `json:"result"`
}
