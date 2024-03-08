package e2e

import (
	"github.com/cosmos/cosmos-sdk/codec"
	appparams "github.com/cosmos/cosmos-sdk/simapp/params"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/settlus/chain/app"
)

var (
	encodingConfig appparams.EncodingConfig
	cdc            codec.Codec
)

func init() {
	encodingConfig = app.MakeEncodingConfig()
	authvesting.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	stakingtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	evidencetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cdc = encodingConfig.Codec
}
