package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	evmtypes "github.com/evmos/evmos/v19/x/evm/types"
)

var (
	//go:embed compiled_contracts/ERC20NonTransferable.json
	SBTJSON []byte //nolint: golint

	// SBTContract is the compiled SBT contract
	SBTContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(SBTJSON, &SBTContract)
	if err != nil {
		panic(err)
	}

	if len(SBTContract.Bin) == 0 {
		panic("load contract failed")
	}
}
