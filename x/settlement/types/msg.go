package types

import (
	"math/big"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgRecord{}
	_ sdk.Msg = &MsgCancel{}
	_ sdk.Msg = &MsgCreateTenant{}
	_ sdk.Msg = &MsgAddTenantAdmin{}
	_ sdk.Msg = &MsgRemoveTenantAdmin{}
	_ sdk.Msg = &MsgUpdateTenantPayoutPeriod{}
	_ sdk.Msg = &MsgDepositToTreasury{}
)

const (
	TypeMsgRecord                   = "record"
	TypeMsgCancel                   = "cancel"
	TypeMsgCreateTenant             = "create_tenant"
	TypeMsgAddTenantAdmin           = "add_tenant_admin"
	TypeMsgRemoveTenantAdmin        = "remove_tenant_admin"
	TypeMsgUpdateTenantPayoutPeriod = "update_tenant_payout_period"
	TypeMsgDepositToTreasury        = "deposit_to_treausry"
)

func NewMsgRecord(sender string, tenantId uint64, requestId string, amount sdk.Coin, chainId, contractAddress string, tokenIdHex string, metadata string) *MsgRecord {
	return &MsgRecord{
		Sender:          sender,
		TenantId:        tenantId,
		RequestId:       requestId,
		Amount:          amount,
		ChainId:         chainId,
		ContractAddress: contractAddress,
		TokenIdHex:      tokenIdHex,
		Metadata:        metadata,
	}
}

func (msg *MsgRecord) Route() string {
	return RouterKey
}

func (msg *MsgRecord) Type() string {
	return TypeMsgRecord
}

func (msg *MsgRecord) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRecord) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.Amount.IsZero() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "amount cannot be zero")
	}

	if !common.IsHexAddress(msg.ContractAddress) {
		return errorsmod.Wrapf(ErrInvalidContractAddress, "invalid contract hex address (%s)", msg.ContractAddress)
	}

	contractAddress := common.HexToAddress(msg.ContractAddress)

	if contractAddress == common.HexToAddress("") {
		return errorsmod.Wrapf(ErrInvalidContractAddress, "invalid contract address (%s)", msg.ContractAddress)
	}

	if msg.TokenIdHex == "" || len(msg.TokenIdHex) < 3 || !strings.HasPrefix(msg.TokenIdHex, "0x") {
		return errorsmod.Wrapf(ErrInvalidTokenId, "invalid token id (%s)", msg.TokenIdHex)
	}

	i := new(big.Int)
	_, ok := i.SetString(msg.TokenIdHex[2:], 16)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidTokenId, "invalid token id (%s)", msg.TokenIdHex)
	}

	return nil
}

func NewMsgCancel(sender string, tenantId uint64, requestId string) *MsgCancel {
	return &MsgCancel{
		Sender:    sender,
		TenantId:  tenantId,
		RequestId: requestId,
	}
}

func (msg *MsgCancel) Route() string {
	return RouterKey
}

func (msg *MsgCancel) Type() string {
	return TypeMsgCancel
}

func (msg *MsgCancel) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgCancel) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancel) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func NewMsgCreateTenant(sender string, denom string, payoutPeriod uint64) *MsgCreateTenant {
	return &MsgCreateTenant{
		Sender:       sender,
		Denom:        denom,
		PayoutPeriod: payoutPeriod,
	}
}

func (msg *MsgCreateTenant) Route() string {
	return RouterKey
}

func (msg *MsgCreateTenant) Type() string {
	return TypeMsgCreateTenant
}

func (msg *MsgCreateTenant) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgCreateTenant) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTenant) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.Denom == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom cannot be empty")
	}

	if msg.PayoutPeriod == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "payout period cannot be zero")
	}

	return nil
}

func NewMsgCreateTenantWithMintableContract(sender string, denom string, payoutPeriod uint64, contractAddress string) *MsgCreateTenantWithMintableContract {
	msg := &MsgCreateTenantWithMintableContract{
		Sender:       sender,
		Denom:        denom,
		PayoutPeriod: payoutPeriod,
	}

	if contractAddress != "" {
		msg.XContractAddress = &MsgCreateTenantWithMintableContract_ContractAddress{ContractAddress: contractAddress}
	}

	return msg
}

func (msg *MsgCreateTenantWithMintableContract) Route() string {
	return RouterKey
}

func (msg *MsgCreateTenantWithMintableContract) Type() string {
	return TypeMsgCreateTenant
}

func (msg *MsgCreateTenantWithMintableContract) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgCreateTenantWithMintableContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTenantWithMintableContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.Denom == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom cannot be empty")
	}

	if msg.PayoutPeriod == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "payout period cannot be zero")
	}

	if msg.GetContractAddress() != "" && !common.IsHexAddress(msg.GetContractAddress()) {
		return errorsmod.Wrapf(ErrInvalidContractAddress, "invalid contract hex address (%s)", msg.GetContractAddress())
	}

	return nil
}

func NewMsgAddTenantAdmin(sender string, tenantId uint64, newAdmin string) *MsgAddTenantAdmin {
	return &MsgAddTenantAdmin{
		Sender:   sender,
		TenantId: tenantId,
		NewAdmin: newAdmin,
	}
}

func (msg *MsgAddTenantAdmin) Route() string {
	return RouterKey
}

func (msg *MsgAddTenantAdmin) Type() string {
	return TypeMsgAddTenantAdmin
}

func (msg *MsgAddTenantAdmin) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgAddTenantAdmin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddTenantAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.NewAdmin == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "admin cannot be empty")
	}

	admin, err := sdk.AccAddressFromBech32(msg.NewAdmin)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid admin address (%s)", err)
	}

	if admin.Empty() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "admin cannot be empty")
	}

	return nil
}

func NewMsgRemoveTenantAdmin(sender string, tenantId uint64, admin string) *MsgRemoveTenantAdmin {
	return &MsgRemoveTenantAdmin{
		Sender:        sender,
		TenantId:      tenantId,
		AdminToRemove: admin,
	}
}

func (msg *MsgRemoveTenantAdmin) Route() string {
	return RouterKey
}

func (msg *MsgRemoveTenantAdmin) Type() string {
	return TypeMsgRemoveTenantAdmin
}

func (msg *MsgRemoveTenantAdmin) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgRemoveTenantAdmin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveTenantAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.AdminToRemove == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "adminToRemove cannot be empty")
	}

	adminToRemove, err := sdk.AccAddressFromBech32(msg.AdminToRemove)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid adminToRemove address (%s)", err)
	}

	if adminToRemove.Empty() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "adminToRemove cannot be empty")
	}

	return nil
}

func NewMsgUpdateTenantPayoutPeriod(sender string, tenantId uint64, payoutPeriod uint64) *MsgUpdateTenantPayoutPeriod {
	return &MsgUpdateTenantPayoutPeriod{
		Sender:       sender,
		TenantId:     tenantId,
		PayoutPeriod: payoutPeriod,
	}
}

func (msg *MsgUpdateTenantPayoutPeriod) Route() string {
	return RouterKey
}

func (msg *MsgUpdateTenantPayoutPeriod) Type() string {
	return TypeMsgUpdateTenantPayoutPeriod
}

func (msg *MsgUpdateTenantPayoutPeriod) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUpdateTenantPayoutPeriod) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateTenantPayoutPeriod) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.PayoutPeriod == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "payout period cannot be zero")
	}

	return nil
}

func NewMsgDepositToTreasury(sender string, tenantId uint64, amount sdk.Coin) *MsgDepositToTreasury {
	return &MsgDepositToTreasury{
		Sender:   sender,
		TenantId: tenantId,
		Amount:   amount,
	}
}

func (msg *MsgDepositToTreasury) Route() string {
	return RouterKey
}

func (msg *MsgDepositToTreasury) Type() string {
	return TypeMsgDepositToTreasury
}

func (msg *MsgDepositToTreasury) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgDepositToTreasury) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDepositToTreasury) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.Amount.IsZero() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "amount cannot be zero")
	}

	return nil
}
