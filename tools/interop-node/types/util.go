package types

import (
	"fmt"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

func ValidateHexString(s string) bool {
	s = strings.TrimPrefix(s, "0x")

	if len(s)%2 == 1 {
		s = "0" + s
	}

	_, err := hexutil.Decode(s)
	return err == nil
}

// GetAddressFromPubKey returns the address of a public key
func GetAddressFromPubKey(pubKey cryptotypes.PubKey) (string, error) {
	acc := authtypes.NewBaseAccount(pubKey.Address().Bytes(), pubKey, 0, 0)
	if err := acc.Validate(); err != nil {
		return "", fmt.Errorf("failed to validate account: %w", err)
	}

	return acc.GetAddress().String(), nil
}
