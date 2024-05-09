package keeper

import (
	"context"
	"fmt"
	"strconv"

	ctypes "github.com/settlus/chain/types"
	settlustypes "github.com/settlus/chain/types"
	"github.com/settlus/chain/x/settlement/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	*SettlementKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *SettlementKeeper) types.MsgServer {
	return &msgServer{SettlementKeeper: keeper}
}

var _ types.MsgServer = msgServer{}

// Record implements MsgServer.Record
func (k msgServer) Record(goCtx context.Context, msg *types.MsgRecord) (*types.MsgRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	if !k.CheckAdminPermission(ctx, msg.TenantId, msg.Sender) {
		return nil, types.ErrNotAuthorized
	}

	tenant := k.GetTenant(ctx, msg.TenantId)
	if tenant.Denom != msg.Amount.Denom {
		return nil, types.ErrInvalidRequest
	}

	if tenant.PayoutPeriod == 0 {
		return nil, types.ErrInvalidTenant
	}

	recipients, err := k.GetRecipients(ctx, msg.ChainId, msg.ContractAddress, msg.TokenIdHex)
	if err != nil {
		return nil, err
	}

	payoutBlock := uint64(ctx.BlockHeight())

	utxrId, err := k.CreateUTXR(
		ctx,
		msg.TenantId,
		&types.UTXR{
			RequestId:  msg.RequestId,
			Recipients: recipients,
			Amount:     msg.Amount,
			Nft: &ctypes.Nft{
				ChainId:      msg.ChainId,
				ContractAddr: settlustypes.HexAddressString(msg.ContractAddress),
				TokenId:      settlustypes.HexAddressString(msg.TokenIdHex),
			},
			CreatedAt: payoutBlock,
		},
	)
	if err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvents(&types.EventRecord{
		Sender:    msg.Sender,
		Tenant:    msg.TenantId,
		UtxrId:    utxrId,
		RequestId: msg.RequestId,
		Amount:    msg.Amount,
		Nft: &ctypes.Nft{
			ChainId:      msg.ChainId,
			ContractAddr: settlustypes.HexAddressString(msg.ContractAddress),
			TokenId:      settlustypes.HexAddressString(msg.TokenIdHex),
		},
		Recipients: recipients,
		Metadata:   msg.Metadata,
		CreatedAt:  uint64(ctx.BlockHeight()),
	}); err != nil {
		return nil, errorsmod.Wrapf(types.ErrEventCreationFailed, "EventRecord event creation failed")
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, "utxr"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("tenant_id", strconv.Itoa(int(msg.TenantId))),
		},
	)

	return &types.MsgRecordResponse{
		UtxrId: utxrId,
	}, nil
}

// Cancel implements MsgServer.Cancel
func (k msgServer) Cancel(goCtx context.Context, msg *types.MsgCancel) (*types.MsgCancelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if exist := k.CheckTenantExist(ctx, msg.TenantId); !exist {
		return nil, types.ErrInvalidTenant
	}
	if !k.CheckAdminPermission(ctx, msg.TenantId, msg.Sender) {
		return nil, types.ErrNotAuthorized
	}

	utxrId, err := k.DeleteUTXRByRequestId(ctx, msg.TenantId, msg.RequestId)
	if err != nil {
		return nil, types.ErrInvalidTxId
	}

	if err := ctx.EventManager().EmitTypedEvents(&types.EventCancel{
		Tenant: msg.TenantId,
		UtxrId: utxrId,
	}); err != nil {
		return nil, errorsmod.Wrapf(types.ErrEventCreationFailed, "EventCancel event creation failed")
	}

	defer telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, "cancel"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("tenant_id", strconv.Itoa(int(msg.TenantId))),
		},
	)

	return &types.MsgCancelResponse{}, nil
}

// CreateTenant implements MsgServer.CreateTenant
func (k msgServer) CreateTenant(goCtx context.Context, msg *types.MsgCreateTenant) (*types.MsgCreateTenantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	tenantId, err := k.CreateNewTenant(ctx, msg.Sender, msg.Denom, msg.PayoutPeriod, types.PayoutMethod_Native, "")
	return &types.MsgCreateTenantResponse{TenantId: tenantId}, err
}

func (k msgServer) CreateTenantWithMintableContract(goCtx context.Context, msg *types.MsgCreateTenantWithMintableContract) (*types.MsgCreateTenantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	tenantId, err := k.CreateNewTenant(ctx, msg.Sender, msg.Denom, msg.PayoutPeriod, types.PayoutMethod_MintContract, msg.GetContractAddress())
	return &types.MsgCreateTenantResponse{TenantId: tenantId}, err
}

// AddTenantAdmin implements MsgServer.AddTenantAdmin
func (k msgServer) AddTenantAdmin(goCtx context.Context, msg *types.MsgAddTenantAdmin) (*types.MsgAddTenantAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if !k.CheckAdminPermission(ctx, msg.TenantId, msg.Sender) {
		return nil, types.ErrNotAuthorized
	}

	tenant := k.GetTenant(ctx, msg.TenantId)
	if tenant == nil {
		return nil, types.ErrInvalidTenant
	}

	for _, admin := range tenant.Admins {
		if admin == msg.NewAdmin {
			return nil, errorsmod.Wrapf(types.ErrInvalidAdmin, "admin %s already exists", msg.NewAdmin)
		}
	}

	tenant.Admins = append(tenant.Admins, msg.NewAdmin)
	k.SetTenant(ctx, tenant)

	return &types.MsgAddTenantAdminResponse{}, nil
}

// RemoveTenantAdmin implements MsgServer.RemoveTenantAdmin
func (k msgServer) RemoveTenantAdmin(goCtx context.Context, msg *types.MsgRemoveTenantAdmin) (*types.MsgRemoveTenantAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if !k.CheckAdminPermission(ctx, msg.TenantId, msg.Sender) {
		return nil, types.ErrNotAuthorized
	}

	tenant := k.GetTenant(ctx, msg.TenantId)
	if tenant == nil {
		return nil, types.ErrInvalidTenant
	}

	for i, admin := range tenant.Admins {
		if admin == msg.AdminToRemove {
			if len(tenant.Admins) == 1 {
				return nil, errorsmod.Wrapf(types.ErrCannotRemoveAdmin, "cannot remove the last admin")
			}
			tenant.Admins = append(tenant.Admins[:i], tenant.Admins[i+1:]...)
			k.SetTenant(ctx, tenant)
			return &types.MsgRemoveTenantAdminResponse{}, nil
		}
	}

	return nil, errorsmod.Wrapf(types.ErrInvalidAdmin, "admin %s does not exist", msg.AdminToRemove)
}

// UpdateTenantPayoutPeriod implements MsgServer.UpdateTenantPayoutPeriod
func (k msgServer) UpdateTenantPayoutPeriod(goCtx context.Context, msg *types.MsgUpdateTenantPayoutPeriod) (*types.MsgUpdateTenantPayoutPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if !k.CheckAdminPermission(ctx, msg.TenantId, msg.Sender) {
		return nil, types.ErrNotAuthorized
	}

	tenant := k.GetTenant(ctx, msg.TenantId)
	if tenant == nil {
		return nil, types.ErrInvalidTenant
	}

	tenant.PayoutPeriod = msg.PayoutPeriod
	k.SetTenant(ctx, tenant)

	return &types.MsgUpdateTenantPayoutPeriodResponse{}, nil
}

// DepositToTreasury implements MsgServer.DepositToTreasury
func (k msgServer) DepositToTreasury(goCtx context.Context, msg *types.MsgDepositToTreasury) (*types.MsgDepositToTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	tenant := k.GetTenant(ctx, msg.TenantId)
	if tenant == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidTenant, "tenant %d does not exist", msg.TenantId)
	}

	if tenant.PayoutMethod != types.PayoutMethod_Native {
		return nil, errorsmod.Wrapf(types.ErrInvalidTenant, "payout method is not native")
	}

	// already checked in ValidateBasic
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	account := k.ak.GetAccount(ctx, sender)
	if account == nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "account %s not found", msg.Sender)
	}

	coins := k.bk.SpendableCoins(ctx, account.GetAddress())
	if coins.AmountOf(msg.Amount.Denom).LT(msg.Amount.Amount) {
		return nil, types.ErrNotEnoughBalance
	}

	if err := k.bk.SendCoins(ctx, account.GetAddress(), types.GetTenantTreasuryAccount(msg.TenantId), sdk.NewCoins(msg.Amount)); err != nil {
		return nil, fmt.Errorf("failed to send coins to treasury account: %w", err)
	}

	defer telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, "deposit_to_treasury"},
		float32(msg.Amount.Amount.Int64()),
		[]metrics.Label{
			telemetry.NewLabel("tenant_id", strconv.Itoa(int(msg.TenantId))),
			telemetry.NewLabel("denom", msg.Amount.Denom),
		},
	)

	return &types.MsgDepositToTreasuryResponse{}, nil
}
