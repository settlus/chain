package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

func EthAddress() common.Address {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return common.BytesToAddress(addr.Bytes())
}

// ValAddress returns a sample validator address
func ValAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.ValAddress(addr).String()
}
