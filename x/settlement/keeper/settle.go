package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/settlus/chain/contracts"
	erc20types "github.com/settlus/chain/evmos/x/erc20/types"
	"github.com/settlus/chain/x/settlement/types"
)

func (k SettlementKeeper) Settle(ctx sdk.Context, tenantId uint64) {
	logger := k.Logger(ctx)

	if err := k.settleUTXRs(ctx, tenantId); err != nil {
		if err := ctx.EventManager().EmitTypedEvents(&types.EventSettlementFailed{
			Tenant: tenantId,
			Reason: err.Error(),
		}); err != nil {
			logger.Error("failed to emit EventSettlementFailed event", "tenant", tenantId, "error", err)
		}
		// if settlement fails, panic. However, this should not happen.
		panic(fmt.Errorf("failed to settle: %w", err))
	}
}

func (k SettlementKeeper) settleUTXRs(ctx sdk.Context, tenantId uint64) error {
	logger := k.Logger(ctx)
	store := k.GetUTXRStore(ctx, tenantId)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var utxr types.UTXR
		k.cdc.MustUnmarshal(iterator.Value(), &utxr)

		utxrId := sdk.BigEndianToUint64(iterator.Key())

		if utxr.PayoutBlock > uint64(ctx.BlockHeight()) {
			logger.Debug("skip payout", "tenant", tenantId, "payoutBlock", utxr.PayoutBlock, "currentBlock", ctx.BlockHeight())
			break
		}

		if err := k.tryPayout(ctx, tenantId, &utxr); err != nil {
			logger.Error("failed to payout", "tenant", tenantId, "recipient", utxr.Recipient, "amount", utxr.Amount.String(), "error", err)
			break
		}

		recipientCosmosAddr := sdk.AccAddress(common.FromHex(utxr.Recipient.String()))
		if err := ctx.EventManager().EmitTypedEvents(&types.EventSettled{
			Tenant:    tenantId,
			UtxrId:    utxrId,
			RequestId: utxr.RequestId,
			Amount:    utxr.Amount,
			Recipient: recipientCosmosAddr.String(),
		}); err != nil {
			logger.Error("failed to emit EventSettled event", "tenant", tenantId, "error", err)
		}

		if err := k.deleteUTXR(ctx, tenantId, utxrId); err != nil {
			// if delete fails after sending coins, return error.
			// this should not happen since we already checked the spendable balance.
			return fmt.Errorf("failed to delete utxr [%d] for tenant [%d]: %w", utxrId, tenantId, err)
		}
	}

	return nil
}

func (k SettlementKeeper) tryPayout(ctx sdk.Context, tenantId uint64, utxr *types.UTXR) (err error) {
	treasuryAddr := types.GetTenantTreasuryAccount(tenantId)
	recipientCosmosAddr := sdk.AccAddress(common.FromHex(utxr.Recipient.String()))
	tenant := k.GetTenant(ctx, tenantId)
	if tenant == nil {
		return fmt.Errorf("tenant [%d] not found", tenantId)
	}

	switch payoutMethod := tenant.PayoutMethod; {
	case payoutMethod == types.PayoutMethod_Native:
		if k.erc20k.IsDenomRegistered(ctx, utxr.Amount.Denom) {
			msg := erc20types.NewMsgConvertCoin(utxr.Amount, common.BytesToAddress(recipientCosmosAddr), treasuryAddr)
			_, err = k.erc20k.ConvertCoin(ctx, msg)
		} else {
			err = k.bk.SendCoins(ctx, treasuryAddr, recipientCosmosAddr, sdk.NewCoins(utxr.Amount))
		}
	case payoutMethod == types.PayoutMethod_MintContract:
		contractAddr := tenant.GetContractAddress()
		_, err = k.evmk.CallEVM(ctx, contracts.SBTContract.ABI, common.BytesToAddress(treasuryAddr), common.HexToAddress(contractAddr),
			true, "mint", common.BytesToAddress(recipientCosmosAddr), utxr.Amount.Amount.BigInt())
	default:
		return fmt.Errorf("invalid payout method: %s", payoutMethod)
	}

	return err
}
