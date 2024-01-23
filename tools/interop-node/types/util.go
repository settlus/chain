package types

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
