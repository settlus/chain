package ante_test

import (
	evmosante "github.com/evmos/evmos/v19/app/ante"
	ethante "github.com/evmos/evmos/v19/app/ante/evm"
	evmosencoding "github.com/evmos/evmos/v19/encoding"
	evmostypes "github.com/evmos/evmos/v19/types"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/app/ante"
)

func (suite *AnteTestSuite) TestValidateHandlerOptions() {
	cases := []struct {
		name    string
		options ante.HandlerOptions
		expPass bool
	}{
		{
			"fail - empty options",
			ante.HandlerOptions{},
			false,
		},
		{
			"fail - empty account keeper",
			ante.HandlerOptions{
				Cdc:           suite.app.AppCodec(),
				AccountKeeper: nil,
			},
			false,
		},
		{
			"fail - empty bank keeper",
			ante.HandlerOptions{
				Cdc:           suite.app.AppCodec(),
				AccountKeeper: suite.app.AccountKeeper,
				BankKeeper:    nil,
			},
			false,
		},
		{
			"fail - empty distribution keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: nil,

				IBCKeeper: nil,
			},
			false,
		},
		{
			"fail - empty IBC keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,

				IBCKeeper: nil,
			},
			false,
		},
		{
			"fail - empty staking keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,

				IBCKeeper:     suite.app.IBCKeeper,
				StakingKeeper: nil,
			},
			false,
		},
		{
			"fail - empty fee market keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,

				IBCKeeper:       suite.app.IBCKeeper,
				StakingKeeper:   suite.app.StakingKeeper,
				FeeMarketKeeper: nil,
			},
			false,
		},
		{
			"fail - empty EVM keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				IBCKeeper:          suite.app.IBCKeeper,
				StakingKeeper:      suite.app.StakingKeeper,
				FeeMarketKeeper:    suite.app.FeeMarketKeeper,
				EvmKeeper:          nil,
			},
			false,
		},
		{
			"fail - empty signature gas consumer",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				IBCKeeper:          suite.app.IBCKeeper,
				StakingKeeper:      suite.app.StakingKeeper,
				FeeMarketKeeper:    suite.app.FeeMarketKeeper,
				EvmKeeper:          suite.app.EvmKeeper,
				SigGasConsumer:     nil,
			},
			false,
		},
		{
			"fail - empty signature mode handler",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				IBCKeeper:          suite.app.IBCKeeper,
				StakingKeeper:      suite.app.StakingKeeper,
				FeeMarketKeeper:    suite.app.FeeMarketKeeper,
				EvmKeeper:          suite.app.EvmKeeper,
				SigGasConsumer:     evmosante.SigVerificationGasConsumer,
				SignModeHandler:    nil,
			},
			false,
		},
		{
			"fail - empty tx fee checker",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				IBCKeeper:          suite.app.IBCKeeper,
				StakingKeeper:      suite.app.StakingKeeper,
				FeeMarketKeeper:    suite.app.FeeMarketKeeper,
				EvmKeeper:          suite.app.EvmKeeper,
				SigGasConsumer:     evmosante.SigVerificationGasConsumer,
				SignModeHandler:    suite.app.GetTxConfig().SignModeHandler(),
				TxFeeChecker:       nil,
			},
			false,
		},
		{
			"fail - empty settlement keeper",
			ante.HandlerOptions{
				Cdc:                    suite.app.AppCodec(),
				AccountKeeper:          suite.app.AccountKeeper,
				BankKeeper:             suite.app.BankKeeper,
				DistributionKeeper:     suite.app.DistrKeeper,
				ExtensionOptionChecker: evmostypes.HasDynamicFeeExtensionOption,
				EvmKeeper:              suite.app.EvmKeeper,
				StakingKeeper:          suite.app.StakingKeeper,
				FeegrantKeeper:         suite.app.FeeGrantKeeper,
				IBCKeeper:              suite.app.IBCKeeper,
				FeeMarketKeeper:        suite.app.FeeMarketKeeper,
				SignModeHandler:        evmosencoding.MakeConfig(app.ModuleBasics).TxConfig.SignModeHandler(),
				SigGasConsumer:         evmosante.SigVerificationGasConsumer,
				MaxTxGasWanted:         40000000,
				TxFeeChecker:           ethante.NewDynamicFeeChecker(suite.app.EvmKeeper),
				SettlementKeeper:       nil,
			},
			false,
		},
		{
			"fail - empty oracle keeper",
			ante.HandlerOptions{
				Cdc:                    suite.app.AppCodec(),
				AccountKeeper:          suite.app.AccountKeeper,
				BankKeeper:             suite.app.BankKeeper,
				DistributionKeeper:     suite.app.DistrKeeper,
				ExtensionOptionChecker: evmostypes.HasDynamicFeeExtensionOption,
				EvmKeeper:              suite.app.EvmKeeper,
				StakingKeeper:          suite.app.StakingKeeper,
				FeegrantKeeper:         suite.app.FeeGrantKeeper,
				IBCKeeper:              suite.app.IBCKeeper,
				FeeMarketKeeper:        suite.app.FeeMarketKeeper,
				SignModeHandler:        evmosencoding.MakeConfig(app.ModuleBasics).TxConfig.SignModeHandler(),
				SigGasConsumer:         evmosante.SigVerificationGasConsumer,
				MaxTxGasWanted:         40000000,
				TxFeeChecker:           ethante.NewDynamicFeeChecker(suite.app.EvmKeeper),
				SettlementKeeper:       suite.app.SettlementKeeper,
				OracleKeeper:           nil,
			},
			false,
		},
		{
			"success - default app options",
			ante.HandlerOptions{
				Cdc:                    suite.app.AppCodec(),
				AccountKeeper:          suite.app.AccountKeeper,
				BankKeeper:             suite.app.BankKeeper,
				DistributionKeeper:     suite.app.DistrKeeper,
				ExtensionOptionChecker: evmostypes.HasDynamicFeeExtensionOption,
				EvmKeeper:              suite.app.EvmKeeper,
				StakingKeeper:          suite.app.StakingKeeper,
				FeegrantKeeper:         suite.app.FeeGrantKeeper,
				IBCKeeper:              suite.app.IBCKeeper,
				FeeMarketKeeper:        suite.app.FeeMarketKeeper,
				SignModeHandler:        evmosencoding.MakeConfig(app.ModuleBasics).TxConfig.SignModeHandler(),
				SigGasConsumer:         evmosante.SigVerificationGasConsumer,
				MaxTxGasWanted:         40000000,
				TxFeeChecker:           ethante.NewDynamicFeeChecker(suite.app.EvmKeeper),
				SettlementKeeper:       suite.app.SettlementKeeper,
				OracleKeeper:           suite.app.OracleKeeper,
			},
			true,
		},
	}

	for _, tc := range cases {
		err := tc.options.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
