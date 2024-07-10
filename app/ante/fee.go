package ante

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	oracletypes "github.com/settlus/chain/x/oracle/types"
)

// DeductFeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeDecorator struct {
	accountKeeper    authante.AccountKeeper
	bankKeeper       authtypes.BankKeeper
	feegrantKeeper   authante.FeegrantKeeper
	settlementKeeper SettlementKeeper
	txFeeChecker     authante.TxFeeChecker
}

func NewDeductFeeDecorator(ak authante.AccountKeeper, bk authtypes.BankKeeper, fk authante.FeegrantKeeper, sk SettlementKeeper) DeductFeeDecorator {
	tfc := newSettlementFeeChecker(sk)

	return DeductFeeDecorator{
		accountKeeper:    ak,
		bankKeeper:       bk,
		feegrantKeeper:   fk,
		settlementKeeper: sk,
		txFeeChecker:     tfc,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if isOracleTx(tx) {
		return next(ctx, tx, simulate)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidGasLimit, "must provide positive gas")
	}

	var (
		priority int64
		err      error
	)

	fee := feeTx.GetFee()
	if !simulate {
		fee, priority, err = dfd.txFeeChecker(ctx, tx)
		if err != nil {
			return ctx, err
		}
	}
	if err := dfd.checkDeductFee(ctx, tx, fee); err != nil {
		return ctx, err
	}

	newCtx := ctx.WithPriority(priority)

	return next(newCtx, tx, simulate)
}

func (dfd DeductFeeDecorator) checkDeductFee(ctx sdk.Context, sdkTx sdk.Tx, fee sdk.Coins) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := dfd.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName); addr == nil {
		return fmt.Errorf("fee collector module account (%s) has not been set", authtypes.FeeCollectorName)
	}

	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if dfd.feegrantKeeper == nil {
			return sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, sdkTx.GetMsgs())
			if err != nil {
				return errorsmod.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return sdkerrors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
	}

	// deduct the fees
	if !fee.IsZero() {
		err := DeductFees(dfd.bankKeeper, dfd.settlementKeeper, ctx, deductFeesFromAcc, fee)
		if err != nil {
			return err
		}
	}

	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, deductFeesFrom.String()),
		),
	}
	ctx.EventManager().EmitEvents(events)

	return nil
}

// DeductFees deducts fees from the given account.
func DeductFees(bankKeeper authtypes.BankKeeper, settlementKeeper SettlementKeeper, ctx sdk.Context, acc authtypes.AccountI, fees sdk.Coins) error {
	if !fees.IsValid() {
		return errorsmod.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	params := settlementKeeper.GetParams(ctx)
	ofp := params.OracleFeePercentage
	gasFees, oracleFees := CalculateFees(ofp, fees)

	if err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), authtypes.FeeCollectorName, gasFees); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	if err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), oracletypes.ModuleName, oracleFees); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return nil
}

// CalculateFees calculates gas fees and oracle fees from the given fees.
func CalculateFees(ofp sdk.Dec, fees sdk.Coins) (gasFees, oracleFees sdk.Coins) {
	for _, fee := range fees {
		gasFee, _ := sdk.NewDecCoinFromDec(
			fee.Denom,
			sdk.NewDecFromInt(fee.Amount).Mul(sdk.NewDec(1).Sub(ofp)),
		).TruncateDecimal()
		gasFees = gasFees.Add(gasFee)

		oracleFee, _ := sdk.NewDecCoinFromDec(
			fee.Denom,
			sdk.NewDecFromInt(fee.Amount).Mul(ofp),
		).TruncateDecimal()
		oracleFees = oracleFees.Add(oracleFee)
	}

	return
}

// SettlementGasConsumeDecorator is used as a post-handler for settlement tx,
// designed to allocate a constant amount of gas regardless of the gas consumed during execution.
type SettlementGasConsumeDecorator struct{}

func NewSettlementGasConsumeDecorator() SettlementGasConsumeDecorator {
	return SettlementGasConsumeDecorator{}
}

func (SettlementGasConsumeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	ctx.GasMeter().RefundGas(ctx.GasMeter().GasConsumed(), "reset the gas count")
	ctx.GasMeter().ConsumeGas(calculateGasCost(tx), "apply settlement tx")

	return next(ctx, tx, simulate)
}

type SettlusSetUpContextDecorator struct{}

func NewSettlusSetUpContextDecorator() SettlusSetUpContextDecorator {
	return SettlusSetUpContextDecorator{}
}

func (SettlusSetUpContextDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	gasTx, ok := tx.(authante.GasTx)
	if !ok {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid transaction type %T, expected GasTx", tx)
	}

	// To ignore gas costs from the KV Store
	newCtx := ctx.WithGasMeter(sdk.NewGasMeter(gasTx.GetGas())).WithKVGasConfig(storetypes.GasConfig{}).WithTransientKVGasConfig(storetypes.GasConfig{})

	return next(newCtx, tx, simulate)
}

type SettlusValidatorCheckDecorator struct {
	ork OracleKeeper
}

func NewSettlusValidatorCheckDecorator(oracleKeeper OracleKeeper) SettlusValidatorCheckDecorator {
	return SettlusValidatorCheckDecorator{
		ork: oracleKeeper,
	}
}

func (ovcd SettlusValidatorCheckDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if isSettlementTx(tx) {
		return next(ctx, tx, simulate)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()

	msgs := tx.GetMsgs()
	if len(msgs) > 1 {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Oracle tx should contain one msg per tx")
	}

	valAddr, err := getValidatorFromOracleMsg(msgs[0])
	if err != nil {
		return ctx, err
	}

	ok, err = ovcd.ork.ValidateFeeder(ctx, feePayer.String(), valAddr.String())
	if err != nil && !ok {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func getValidatorFromOracleMsg(msg sdk.Msg) (sdk.ValAddress, error) {
	switch msg := msg.(type) {
	case *oracletypes.MsgVote:
		val, err := sdk.ValAddressFromBech32(msg.Validator)
		if err != nil {
			return nil, err
		}
		return val, nil
	case *oracletypes.MsgPrevote:
		val, err := sdk.ValAddressFromBech32(msg.Validator)
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid oracle msg type")
	}
}
