package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	ctypes "github.com/settlus/chain/types"
)

func TestStringToOwnershipData(t *testing.T) {
	type Result struct {
		nft   ctypes.Nft
		owner ctypes.HexAddressString
		err   bool
	}

	testCases := []struct {
		name       string
		voteString string
		expected   struct {
			nft   ctypes.Nft
			owner ctypes.HexAddressString
			err   bool
		}
	}{
		{
			name:       "Valid vote string",
			voteString: "1/0x123/0x0:0x777",
			expected: Result{
				nft: ctypes.Nft{
					ChainId:      "1",
					ContractAddr: ctypes.NormalizeHexAddress("0x123"),
					TokenId:      ctypes.NormalizeHexAddress("0x0"),
				},
				owner: ctypes.NormalizeHexAddress("0x777"),
				err:   false,
			},
		},
		{
			name:       "Valid vote string",
			voteString: "1/123/0a:777",
			expected: Result{
				nft: ctypes.Nft{
					ChainId:      "1",
					ContractAddr: ctypes.NormalizeHexAddress("0x123"),
					TokenId:      ctypes.NormalizeHexAddress("0xa"),
				},
				owner: ctypes.NormalizeHexAddress("777"),
				err:   false,
			},
		},
		{
			name:       "Invalid nftId - no colone",
			voteString: "1/0x123/0x0/0x777",
			expected: Result{
				nft:   ctypes.Nft{},
				owner: ctypes.HexAddressString(""),
				err:   true,
			},
		},
		{
			name:       "Invalid nftId - no slash",
			voteString: "10x123/0x0:0x777",
			expected: Result{
				nft:   ctypes.Nft{},
				owner: ctypes.HexAddressString(""),
				err:   true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nft, owner, err := StringToOwnershipData(tc.voteString)

			if tc.expected.err {
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected.nft, nft)
				require.Equal(t, tc.expected.owner, owner)
			}
		})
	}
}
