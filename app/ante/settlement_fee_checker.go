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
	SettlementCreateTenantGasCost uint64 = 100000000
)

func newSettlementFeeChecker(k SettlementKeeper) ante.TxFeeChecker {
	return func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return nil, 0, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
		}

		feeCoins := feeTx.GetFee()
		params := k.GetParams(ctx)
		gasPrice := params.GetGasPrice()
		gasRequired := calculateGasCost(tx)
		requiredFees := gasPrice.Amount.Mul(sdk.NewInt(int64(gasRequired)))

		if feeCoins.AmountOf(gasPrice.Denom).LT(requiredFees) {
			return nil, 0, errorsmod.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s, required: %s", feeCoins, requiredFees.String())
		}

		return sdk.NewCoins(sdk.NewCoin(gasPrice.Denom, requiredFees)), 0, nil
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
