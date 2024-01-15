package types

import (
	"strings"

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
