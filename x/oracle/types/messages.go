package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgVote                    = "vote"
	TypeMsgPrevote                 = "prevote"
	TypeMsgFeederDelegationConsent = "feeder_delegation_consent"
)

var (
	_ sdk.Msg = &MsgVote{}
	_ sdk.Msg = &MsgPrevote{}
	_ sdk.Msg = &MsgFeederDelegationConsent{}
)

func NewMsgVote(feeder, validator, blockDataString, salt string, roundId uint64) *MsgVote {
	return &MsgVote{
		Feeder:          feeder,
		Validator:       validator,
		BlockDataString: blockDataString,
		Salt:            salt,
		RoundId:         roundId,
	}
}

func (msg *MsgVote) Route() string {
	return RouterKey
}

func (msg *MsgVote) Type() string {
	return TypeMsgVote
}

func (msg *MsgVote) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	_, err = ParseBlockDataString(msg.BlockDataString)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block data string (%s)", err)
	}

	return nil
}

func NewMsgPrevote(feeder string, validator string, hash string, roundId uint64) *MsgPrevote {
	return &MsgPrevote{
		Feeder:    feeder,
		Validator: validator,
		Hash:      hash,
		RoundId:   roundId,
	}
}

func (msg *MsgPrevote) Route() string {
	return RouterKey
}

func (msg *MsgPrevote) Type() string {
	return TypeMsgPrevote
}

func (msg *MsgPrevote) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgPrevote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPrevote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	return nil
}

func NewMsgFeederDelegationConsent(validator string, feederAddress string) *MsgFeederDelegationConsent {
	return &MsgFeederDelegationConsent{
		Validator:     validator,
		FeederAddress: feederAddress,
	}
}

func (msg *MsgFeederDelegationConsent) Route() string {
	return RouterKey
}

func (msg *MsgFeederDelegationConsent) Type() string {
	return TypeMsgFeederDelegationConsent
}

func (msg *MsgFeederDelegationConsent) GetSigners() []sdk.AccAddress {
	senderValidator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	sender := sdk.AccAddress(senderValidator)
	return []sdk.AccAddress{sender}
}

func (msg *MsgFeederDelegationConsent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFeederDelegationConsent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FeederAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	return nil
}
