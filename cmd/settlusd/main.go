package main

import (
	"os"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/cmd/settlusd/cmd"
	cmdcfg "github.com/settlus/chain/cmd/settlusd/config"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	setConfig()
	cmdcfg.RegisterDenom()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, cmd.EnvPrefix, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}

func setConfig() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()
}
