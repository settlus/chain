package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/app/ante"
)

var _ sdk.PostDecorator = &SettlementDecorator{}

// SettlementDecorator is used for settlement transactions.
type SettlementDecorator struct {
}

// NewSettlementDecorator creates a new instance of the SettlementDecorator.
func NewSettlementDecorator() sdk.PostDecorator {
	return &SettlementDecorator{}
}

func (sd SettlementDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	if ante.IsSettlementTx(tx) {
		ctx.GasMeter().RefundGas(ctx.GasMeter().GasConsumed(), "reset the gas count")
		ctx.GasMeter().ConsumeGas(ante.CalculateGasCost(tx), "apply settlement tx")

		return next(ctx, tx, simulate, success)
	}

	return next(ctx, tx, simulate, success)
}
