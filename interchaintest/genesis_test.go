package interchaintest_test

import (
	"context"

	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

type genesisWrapper struct {
	chain          *cosmos.CosmosChain
	tfRoles        NobleRoles
	fiatTfRoles    NobleRoles
	paramAuthority ibc.Wallet
	extraWallets   ExtraWallets
}

func settlusChainSpec(
	ctx context.Context,
	gw *genesisWrapper,
	chainID string,
	nv, nf int,
	minSetupTf, minSetupFiatTf bool,
	minModifyTf, minModifyFiatTf bool,
) *interchaintest.ChainSpec {
	return &interchaintest.ChainSpec{
		NumValidators: &nv,
		NumFullNodes:  &nf,
		ChainConfig: ibc.ChainConfig{
			Type:           "cosmos",
			Name:           "noble",
			ChainID:        chainID,
			Bin:            "nobled",
			Denom:          "token",
			Bech32Prefix:   "noble",
			CoinType:       "118",
			GasPrices:      "0.0token",
			GasAdjustment:  1.1,
			TrustingPeriod: "504h",
			NoHostMount:    false,
			Images:         nobleImageInfo,
			EncodingConfig: NobleEncoding(),
			PreGenesis:     preGenesisAll(ctx, gw, minSetupTf, minSetupFiatTf),
			ModifyGenesis:  modifyGenesisAll(gw, minModifyTf, minModifyFiatTf),
		},
	}
}
