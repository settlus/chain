package ante

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// RejectMessagesDecorator prevents invalid msg types from being executed
type RejectMessagesDecorator struct{}

// AnteHandle rejects messages that requires ethereum-specific authentication.
// For example `MsgEthereumTx` requires fee to be deducted in the antehandler in
// order to perform the refund.
func (RejectMessagesDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		// check if msgtypeurl starts with '/settlus.settlement'
		if strings.HasPrefix(sdk.MsgTypeURL(msg), "/settlus.settlement") {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrInvalidType,
				"Settlment Msg can only be processed in the Settlement ante handler",
			)
		}

		if sdk.MsgTypeURL(msg) == "/cosmos.staking.v1beta1.MsgCreateValidator" && ctx.BlockHeight() != 0 {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrInvalidType,
				"You can create validator only through governance or gentx",
			)
		}
	}

	return next(ctx, tx, simulate)
}
