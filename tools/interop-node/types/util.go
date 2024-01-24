package types

import (
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func PadBytes(pad int, b []byte) []byte {
	if len(b) == pad {
		return b
	}

	if len(b) > pad {
		return b[:32]
	}

	padded := make([]byte, pad)
	copy(padded[pad-len(b):], b)
	return padded
}

// GetAddressFromPubKey returns the address of a public key
func GetAddressFromPubKey(pubKey cryptotypes.PubKey) (string, error) {
	acc := authtypes.NewBaseAccount(pubKey.Address().Bytes(), pubKey, 0, 0)
	if err := acc.Validate(); err != nil {
		return "", fmt.Errorf("failed to validate account: %w", err)
	}

	return acc.GetAddress().String(), nil
}
