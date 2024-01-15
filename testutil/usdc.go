package testutil

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	microUSDCDenom = "uusdc"
)

func NewMicroUSDC(amt int64) sdk.Coin {
	return sdk.NewInt64Coin(microUSDCDenom, amt)
}

func NewSetl(amt int64) sdk.Coin {
	return sdk.NewInt64Coin("setl", amt)
}
