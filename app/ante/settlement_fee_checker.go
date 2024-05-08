package ante

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

const (
	SettlementBasicGasCost        uint64 = 10000
	SettlementCreateTenantGasCost uint64 = 1000000000000
)

func newSettlementFeeChecker(k SettlementKeeper) ante.TxFeeChecker {
	return func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return nil, 0, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
		}

		feeCoins := feeTx.GetFee()
		params := k.GetParams(ctx)
		gasPrices := params.GetGasPrices()
		gasRequired := calculateGasCost(tx)
		for _, gasPrice := range gasPrices {
			gasPrice := sdk.NormalizeDecCoin(gasPrice)
			requiredFees := gasPrice.Amount.Mul(sdk.NewDec(int64(gasRequired))).TruncateInt()

			if feeCoins.AmountOf(gasPrice.Denom).GTE(requiredFees) {
				return sdk.NewCoins(sdk.NewCoin(gasPrice.Denom, requiredFees)), 0, nil
			}
		}

		return nil, 0, errorsmod.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s", feeCoins)
	}
}

func calculateGasCost(tx sdk.Tx) (gas uint64) {
	for _, msg := range tx.GetMsgs() {
		gas += SettlementBasicGasCost
		url := sdk.MsgTypeURL(msg)
		if strings.HasSuffix(url, "MsgCreateTenant") || strings.HasSuffix(url, "MsgCreateTenantWithMintableContract") {
			gas += SettlementCreateTenantGasCost
		}
	}

	return gas
}
