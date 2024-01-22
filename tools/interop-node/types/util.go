package types

import (
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

// ValidateHexString checks if a string is a valid hex string
func ValidateHexString(s string) bool {
	_, err := hexutil.Decode(s)
	return err == nil
}
