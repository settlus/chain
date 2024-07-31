package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/settlus/chain/contracts"
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

	period := k.GetPayoutPeriod(ctx, tenantId)
	for ; iterator.Valid(); iterator.Next() {
		var utxr types.UTXR
		k.cdc.MustUnmarshal(iterator.Value(), &utxr)

		utxrId := sdk.BigEndianToUint64(iterator.Key())
		payoutBlock := utxr.CreatedAt + period
		if payoutBlock > uint64(ctx.BlockHeight()) {
			logger.Debug("skip payout", "tenant", tenantId, "payoutBlock", payoutBlock, "currentBlock", ctx.BlockHeight())
			break
		}

		valid, err := k.tryPayout(ctx, tenantId, &utxr)
		if err != nil {
			logger.Error("failed to payout", "tenant", tenantId, "recipient", utxr.Recipients, "amount", utxr.Amount.String(), "error", err)
			break
		}

		if valid {
			if err := ctx.EventManager().EmitTypedEvents(&types.EventSettled{
				Tenant: tenantId,
				UtxrId: utxrId,
			}); err != nil {
				logger.Error("failed to emit EventSettled event", "tenant", tenantId, "error", err)
			}
		} else {
			if err := ctx.EventManager().EmitTypedEvents(&types.EventCancel{
				Tenant: tenantId,
				UtxrId: utxrId,
			}); err != nil {
				logger.Error("failed to emit EventCancel event", "tenant", tenantId, "utxr", utxrId, "error", err)
			}
		}

		if err := k.deleteUTXR(ctx, tenantId, utxrId); err != nil {
			// if delete fails after sending coins, return error.
			// this should not happen since we already checked the spendable balance.
			return fmt.Errorf("failed to delete utxr [%d] for tenant [%d]: %w", utxrId, tenantId, err)
		}
	}

	return nil
}

func (k SettlementKeeper) tryPayout(ctx sdk.Context, tenantId uint64, utxr *types.UTXR) (bool, error) {
	treasuryAddr := types.GetTenantTreasuryAccount(tenantId)
	tenant := k.GetTenant(ctx, tenantId)
	if tenant == nil {
		return false, fmt.Errorf("tenant [%d] not found", tenantId)
	}

	var validRecipients []*types.Recipient = make([]*types.Recipient, 0)
	var totalWeight uint32 = 0
	for _, recipient := range utxr.Recipients {
		if !recipient.Address.IsNull() {
			totalWeight += recipient.Weight
			validRecipients = append(validRecipients, recipient)
		}
	}

	if len(validRecipients) == 0 {
		return false, nil
	}

	for _, recipient := range validRecipients {
		recipientCosmosAddr := sdk.AccAddress(common.FromHex(recipient.Address.String()))

		// if total weight is 0, payout should be distributed equally
		// otherwise it will be distributed based on the weight
		amount := sdk.Coin{
			Denom:  utxr.Amount.Denom,
			Amount: math.ZeroInt(),
		}

		if totalWeight == 0 {
			amount.Amount = utxr.Amount.Amount.Quo(sdk.NewInt(int64(len(validRecipients))))
		} else {
			amount.Amount = utxr.Amount.Amount.Mul(sdk.NewInt(int64(recipient.Weight))).Quo(sdk.NewInt(int64(totalWeight)))
		}

		var err error
		switch payoutMethod := tenant.PayoutMethod; {
		case payoutMethod == types.PayoutMethod_Native:
			if k.erc20k.IsDenomRegistered(ctx, amount.Denom) {
				// convert from erc20 to Coin
				// TODO
				//pair := erc20types.NewTokenPair()
				//_, err = k.erc20k.ConvertCoinNativeERC20(
				//	ctx, erc20types.MsgConvertERC20(amount)
				//)
			} else {
				err = k.bk.SendCoins(ctx, treasuryAddr, recipientCosmosAddr, sdk.NewCoins(amount))
			}
		case payoutMethod == types.PayoutMethod_MintContract:
			contractAddr := tenant.GetContractAddress()
			_, err = k.evmk.CallEVM(ctx, contracts.SBTContract.ABI, common.BytesToAddress(treasuryAddr), common.HexToAddress(contractAddr),
				true, "mint", common.BytesToAddress(recipientCosmosAddr), amount.Amount.BigInt())
		default:
			return true, fmt.Errorf("invalid payout method: %s", payoutMethod)
		}
		if err != nil {
			return true, err
		}
	}

	return true, nil
}
