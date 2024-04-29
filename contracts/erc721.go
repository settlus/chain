package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/settlus/chain/evmos/x/evm/types"

	"github.com/settlus/chain/x/settlement/types"
)

var (
	//go:embed compiled_contracts/ERC721.json
	ERC721JSON     []byte //nolint: golint
	ERC721Contract evmtypes.CompiledContract
	ERC721Address  common.Address
)

func init() {
	ERC721Address = types.ModuleAddress
	err := json.Unmarshal(ERC721JSON, &ERC721Contract)
	if err != nil {
		panic(err)
	}

	if len(ERC721Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
