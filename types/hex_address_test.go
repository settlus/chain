package types_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/settlus/chain/types"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Bech32 string
	Hex    string
	Bytes  []byte
}

var testCases []TestCase = []TestCase{
	{
		Bech32: "settlus18p89lltyryx6w9xzjv0gy2pwqfx4rk5zlc2lf6",
		Hex:    "0x384e5FfD64190DA714c2931e82282e024d51Da82",
		Bytes:  []byte{0x38, 0x4e, 0x5f, 0xfd, 0x64, 0x19, 0xd, 0xa7, 0x14, 0xc2, 0x93, 0x1e, 0x82, 0x28, 0x2e, 0x2, 0x4d, 0x51, 0xda, 0x82},
	},
}

func Test_NewHexAddressString(t *testing.T) {
	for _, testCase := range testCases {
		_, bytes, err := bech32.DecodeAndConvert(testCase.Bech32)
		require.NoError(t, err)
		address := types.NewHexAddrFromAccAddr(bytes)
		require.Equal(t, testCase.Hex, address.String())
	}
}

func Test_HexAddressString_Marshal(t *testing.T) {
	for _, testCase := range testCases {
		address := types.HexAddressString(testCase.Hex)
		data, err := address.Marshal()
		require.NoError(t, err)
		require.Equal(t, testCase.Bytes, data)
	}
}

func TestHexAddressString_MarshalTo(t *testing.T) {
	for _, testCase := range testCases {
		address := types.HexAddressString(testCase.Hex)
		data := make([]byte, address.Size())

		n, err := address.MarshalTo(data)
		require.NoError(t, err)
		require.Equal(t, address.Size(), n)
		require.Equal(t, testCase.Bytes, data)
	}
}

func TestHexAddressString_Unmarshal(t *testing.T) {
	for _, testCase := range testCases {
		address := types.HexAddressString("")
		data := testCase.Bytes

		err := address.Unmarshal(data)
		require.NoError(t, err)

		expectedAddress := types.HexAddressString(testCase.Hex)
		require.Equal(t, expectedAddress, address)
	}
}

func TestHexAddressString_Normalize(t *testing.T) {
	cases := []struct {
		before string
		after  string
	}{
		{
			before: "0x1",
			after:  "0x0000000000000000000000000000000000000001",
		},
		{
			before: "a8",
			after:  "0x00000000000000000000000000000000000000a8",
		},
	}

	for _, tc := range cases {
		hexAddr := types.HexAddressString(tc.before)
		temp := make([]byte, hexAddr.Size())
		size, err := hexAddr.MarshalTo(temp)
		require.NoError(t, err)
		require.Equal(t, 20, size)

		var normHexAddr types.HexAddressString
		err = (&normHexAddr).Unmarshal(temp)
		require.NoError(t, err)
		require.Equal(t, tc.after, normHexAddr.String())
	}
}
