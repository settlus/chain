package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRecord_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRecord
		err  error
	}{
		{
			name: "invalid sender address",
			msg: MsgRecord{
				Sender:          "invalid_address",
				TenantId:        0,
				RequestId:       "request-1",
				Amount:          sdk.NewCoin("uusdc", sdk.NewInt(100)),
				ChainId:         "settlus_5371-1",
				ContractAddress: "0x0000000000000000000000000000000000000001",
				TokenIdHex:      "0x1",
				Metadata:        "metadata",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid contract address",
			msg: MsgRecord{
				Sender:          sample.AccAddress(),
				TenantId:        0,
				RequestId:       "request-1",
				Amount:          sdk.NewCoin("uusdc", sdk.NewInt(100)),
				ChainId:         "settlus_5371-1",
				ContractAddress: "invalid_contract_address",
				TokenIdHex:      "0x1",
				Metadata:        "metadata",
			},
			err: ErrInvalidContractAddress,
		}, {
			name: "invalid amount",
			msg: MsgRecord{
				Sender:          sample.AccAddress(),
				TenantId:        0,
				RequestId:       "request-1",
				Amount:          sdk.NewCoin("uusdc", sdk.NewInt(0)),
				ChainId:         "settlus_5371-1",
				ContractAddress: "0x0000000000000000000000000000000000000001",
				TokenIdHex:      "0x1",
				Metadata:        "metadata",
			},
			err: sdkerrors.ErrInvalidCoins,
		}, {
			name: "empty token id",
			msg: MsgRecord{
				Sender:          sample.AccAddress(),
				TenantId:        0,
				RequestId:       "request-1",
				Amount:          sdk.NewCoin("uusdc", sdk.NewInt(100)),
				ChainId:         "settlus_5371-1",
				ContractAddress: "0x0000000000000000000000000000000000000001",
				TokenIdHex:      "",
				Metadata:        "metadata",
			},
			err: ErrInvalidTokenId,
		}, {
			name: "invalid token id",
			msg: MsgRecord{
				Sender:          sample.AccAddress(),
				TenantId:        0,
				RequestId:       "request-1",
				Amount:          sdk.NewCoin("uusdc", sdk.NewInt(100)),
				ChainId:         "settlus_5371-1",
				ContractAddress: "0x0000000000000000000000000000000000000001",
				TokenIdHex:      "invalid_token_id",
				Metadata:        "metadata",
			},
			err: ErrInvalidTokenId,
		}, {
			name: "token id without 0x prefix",
			msg: MsgRecord{
				Sender:          sample.AccAddress(),
				TenantId:        0,
				RequestId:       "request-1",
				Amount:          sdk.NewCoin("uusdc", sdk.NewInt(100)),
				ChainId:         "settlus_5371-1",
				ContractAddress: "0x0000000000000000000000000000000000000001",
				TokenIdHex:      "1",
				Metadata:        "metadata",
			},
			err: ErrInvalidTokenId,
		}, {
			name: "very large token id",
			msg: MsgRecord{
				Sender:          sample.AccAddress(),
				TenantId:        0,
				RequestId:       "request-1",
				Amount:          sdk.NewCoin("uusdc", sdk.NewInt(100)),
				ChainId:         "settlus_5371-1",
				ContractAddress: "0x0000000000000000000000000000000000000001",
				TokenIdHex:      "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0",
				Metadata:        "metadata",
			},
			err: nil,
		}, {
			name: "valid msg",
			msg: MsgRecord{
				Sender:          sample.AccAddress(),
				TenantId:        0,
				RequestId:       "request-1",
				Amount:          sdk.NewCoin("uusdc", sdk.NewInt(100)),
				ChainId:         "settlus_5371-1",
				ContractAddress: "0x0000000000000000000000000000000000000001",
				TokenIdHex:      "0x1",
				Metadata:        "metadata",
			},
			err: nil,
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

func TestMsgCancel_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCancel
		err  error
	}{
		{
			name: "invalid sender address",
			msg: MsgCancel{
				Sender:    "invalid_address",
				TenantId:  0,
				RequestId: "request-1",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid msg",
			msg: MsgCancel{
				Sender:    sample.AccAddress(),
				TenantId:  0,
				RequestId: "request-1",
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

func TestMsgDepositToTreasury_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDepositToTreasury
		err  error
	}{
		{
			name: "invalid sender address",
			msg: MsgDepositToTreasury{
				Sender:   "invalid_address",
				TenantId: 0,
				Amount:   sdk.NewCoin("uusdc", sdk.NewInt(100)),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid amount",
			msg: MsgDepositToTreasury{
				Sender:   sample.AccAddress(),
				TenantId: 0,
				Amount:   sdk.NewCoin("uusdc", sdk.NewInt(0)),
			},
			err: sdkerrors.ErrInvalidCoins,
		}, {
			name: "valid address",
			msg: MsgDepositToTreasury{
				Sender:   sample.AccAddress(),
				TenantId: 0,
				Amount:   sdk.NewCoin("uusdc", sdk.NewInt(100)),
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
