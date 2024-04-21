package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/settlus/chain/types"
)

func SingleRecipients(creator sdk.AccAddress) []*Recipient {
	return []*Recipient{
		{
			Address: ctypes.NewHexAddrFromAccAddr(creator),
		},
	}
}
