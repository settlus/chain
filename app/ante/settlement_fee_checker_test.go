package ante

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/settlus/chain/evmos/encoding"
	settlementtypes "github.com/settlus/chain/x/settlement/types"
)

var _ SettlementKeeper = MockSettlementKeeper{}

type MockSettlementKeeper struct{}

func (m MockSettlementKeeper) GetParams(_ sdk.Context) settlementtypes.Params {
	return settlementtypes.DefaultParams()
}

func Test_SettlementFeeChecker(t *testing.T) {
	encodingConfig := encoding.MakeConfig(module.NewBasicManager())
	txCtx := sdk.NewContext(nil, tmproto.Header{Height: 1}, false, log.NewNopLogger())

	testCases := []struct {
		name       string
		ctx        sdk.Context
		buildTx    func() sdk.FeeTx
		expFees    string
		expSuccess bool
	}{
		{
			"success, record tx",
			txCtx,
			func() sdk.FeeTx {
				txBuilder := encodingConfig.TxConfig.NewTxBuilder()
				msg := &settlementtypes.MsgRecord{}
				err := txBuilder.SetMsgs(msg)
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uusdc", sdk.NewInt(10000))))
				require.NoError(t, err)

				return txBuilder.GetTx()
			},
			"10000uusdc",
			true,
		},
		{
			"fail, insufficient fees for record tx",
			txCtx,
			func() sdk.FeeTx {
				txBuilder := encodingConfig.TxConfig.NewTxBuilder()
				msg := &settlementtypes.MsgRecord{}
				err := txBuilder.SetMsgs(msg)
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uusdc", sdk.NewInt(1))))
				require.NoError(t, err)

				return txBuilder.GetTx()
			},
			"",
			false,
		},
		{
			"success, create tenant",
			txCtx,
			func() sdk.FeeTx {
				txBuilder := encodingConfig.TxConfig.NewTxBuilder()
				msg := &settlementtypes.MsgCreateTenant{}
				err := txBuilder.SetMsgs(msg)
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uusdc", sdk.NewInt(1010000))))
				require.NoError(t, err)

				return txBuilder.GetTx()
			},
			"1010000uusdc",
			true,
		},
		{
			"fail, insufficient fees for create tenant",
			txCtx,
			func() sdk.FeeTx {
				txBuilder := encodingConfig.TxConfig.NewTxBuilder()
				msg := &settlementtypes.MsgCreateTenant{}
				err := txBuilder.SetMsgs(msg)
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uusdc", sdk.NewInt(1000000))))
				require.NoError(t, err)

				return txBuilder.GetTx()
			},
			"",
			false,
		},
		{
			"fail, insufficient fees for create tenant-mc",
			txCtx,
			func() sdk.FeeTx {
				txBuilder := encodingConfig.TxConfig.NewTxBuilder()
				msg := &settlementtypes.MsgCreateTenantWithMintableContract{}
				err := txBuilder.SetMsgs(msg)
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uusdc", sdk.NewInt(1000000))))
				require.NoError(t, err)

				return txBuilder.GetTx()
			},
			"",
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fees, _, err := newSettlementFeeChecker(MockSettlementKeeper{})(tc.ctx, tc.buildTx())
			if tc.expSuccess {
				require.Equal(t, tc.expFees, fees.String())
			} else {
				require.Error(t, err)
			}
		})
	}
}
