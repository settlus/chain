package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRegisterAccount = "register_account"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
)

// NewMsgRegisterAccount returns a new MsgRegisterAccount
func NewMsgRegisterAccount(sender, port, channelId string, timeoutTimestamp uint64) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		Sender:           sender,
		Port:             port,
		ChannelId:        channelId,
		TimeoutTimestamp: timeoutTimestamp,
	}
}

func (msg *MsgRegisterAccount) Route() string {
	return RouterKey
}

func (msg *MsgRegisterAccount) Type() string {
	return TypeMsgRegisterAccount
}

func (msg *MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.Port == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "port cannot be empty")
	}

	if msg.ChannelId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "channel id cannot be empty")
	}

	return nil
}
