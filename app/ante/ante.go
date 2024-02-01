package ante

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// NewAnteHandler returns an ante handler responsible for attempting to route an
// Ethereum or SDK transaction to an internal ante handler for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
func NewAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		if isSettlementTx(tx) {
			return newSettlementAnteHandler(options)(ctx, tx, sim)
		}

		if isCreateValidatorTx(tx) && ctx.BlockHeight() != 0 {
				return ctx, errorsmod.Wrap(errortypes.ErrInvalidRequest, "You can create validator only through governance or gentx")
		}

		if txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx); ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 0 {
				switch typeURL := opts[0].GetTypeUrl(); typeURL {
				case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
					// handle as *evmtypes.MsgEthereumTx
					anteHandler = newEVMAnteHandler(options)
				default:
					return ctx, errorsmod.Wrapf(
						errortypes.ErrUnknownExtensionOptions,
						"rejecting tx with unsupported extension option: %s", typeURL,
					)
				}

				return anteHandler(ctx, tx, sim)
			}
		}

		// handle as totally normal Cosmos SDK tx
		switch tx.(type) {
		case sdk.Tx:
			anteHandler = newCosmosAnteHandler(options)
		default:
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}
}

func isSettlementTx(tx sdk.Tx) bool {
	if len(tx.GetMsgs()) == 0 {
		return false
	}

	for _, msg := range tx.GetMsgs() {
		// EIP-712 Msg can't be built with ExtensionOptions, so we filter settlement messages by MsgTypeURL
		if !strings.HasPrefix(sdk.MsgTypeURL(msg), "/settlus.settlement") {
			return false
		}
	}

	return true
}

func isCreateValidatorTx(tx sdk.Tx) bool {
	if len(tx.GetMsgs()) == 0 {
		return false
	}

	for _, msg := range tx.GetMsgs() {
		if !strings.HasPrefix(sdk.MsgTypeURL(msg), "/cosmos.staking.v1beta1.MsgCreateValidator") {
			return false
		}
	}

	return true
}

func NewPostHandler(options HandlerOptions) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, sim bool) (sdk.Context, error) {
		if isSettlementTx(tx) {
			return sdk.ChainAnteDecorators(
				NewSettlementGasConsumeDecorator(),
			)(ctx, tx, sim)
		}

		return ctx, nil
	}
}
