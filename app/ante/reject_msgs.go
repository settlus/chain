package ante

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// RejectMessagesDecorator prevents invalid msg types from being executed
type RejectMessagesDecorator struct {
	stakingKeeper StakingKeeper
}

func NewRejectMessagesDecorator(sk StakingKeeper) RejectMessagesDecorator {
	return RejectMessagesDecorator{stakingKeeper: sk}
}

// AnteHandle rejects messages that requires ethereum-specific authentication.
// For example `MsgEthereumTx` requires fee to be deducted in the antehandler in
// order to perform the refund.
func (rmd RejectMessagesDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
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

		if sdk.MsgTypeURL(msg) == "/cosmos.staking.v1beta1.MsgDelegate" {
			probono, err := rmd.CheckProbono(ctx, msg)
			if err != nil {
				return ctx, err
			}

			if probono {
				return ctx, errorsmod.Wrapf(
					errortypes.ErrInvalidType,
					"Probono validator can delegate only through governance",
				)
			}
		}
	}

	return next(ctx, tx, simulate)
}

func (rmd RejectMessagesDecorator) CheckProbono(ctx sdk.Context, msg sdk.Msg) (bool, error) {
	delegateMsg, ok := msg.(*stakingtypes.MsgDelegate)
	if !ok {
		return false, errorsmod.Wrapf(
			errortypes.ErrInvalidType, "Invalid message type",
		)
	}

	valAddr, err := sdk.ValAddressFromBech32(delegateMsg.ValidatorAddress)
	if err != nil {
		return false, err
	}

	val, found := rmd.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return false, err
	}

	return val.IsProbono(), nil
}
