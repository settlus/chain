package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/settlus/chain/testutil/sample"
)

func TestMsgPrevote_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgPrevote
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgPrevote{
				Feeder:    "invalid_address",
				Validator: "invalid_address",
				Hash:      "invalid_hash",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgPrevote{
				Feeder:    sample.AccAddress(),
				Validator: sample.ValAddress(),
				Hash:      "valid_hash",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgVote_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgVote
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgVote{
				Feeder:    "invalid_address",
				Validator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgVote{
				Feeder:    sample.AccAddress(),
				Validator: sample.ValAddress(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgFeederDelegationConsent_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgFeederDelegationConsent
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgFeederDelegationConsent{
				FeederAddress: "invalid_address",
				Validator:     "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgFeederDelegationConsent{
				Validator:     sample.ValAddress(),
				FeederAddress: sample.AccAddress(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
