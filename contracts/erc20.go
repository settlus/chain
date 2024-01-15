package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/settlus/chain/evmos/x/evm/types"

	"github.com/settlus/chain/evmos/x/erc20/types"
)

var (
	//go:embed compiled_contracts/ERC20MinterPauserBurnerDecimals.json
	ERC20JSON []byte //nolint: golint

	// ERC20Contract is the compiled erc20 contract
	ERC20Contract evmtypes.CompiledContract

	// ERC20Address is the erc20 module address
	ERC20Address common.Address
)

func init() {
	ERC20Address = types.ModuleAddress

	err := json.Unmarshal(ERC20JSON, &ERC20Contract)
	if err != nil {
		panic(err)
	}

	if len(ERC20Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
