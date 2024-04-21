package types

import (
	"testing"

	ctypes "github.com/settlus/chain/types"
	"github.com/stretchr/testify/require"
)

func TestParseBlockDataString(t *testing.T) {
	tests := []struct {
		name         string
		blockData    BlockData
		blockDatastr string
		wantErr      bool
	}{
		{
			name: "valid block data string",
			blockData: BlockData{
				ChainId:     "1",
				BlockNumber: 100,
				BlockHash:   "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
			},
			blockDatastr: "1:100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
			wantErr:      false,
		}, {
			name:         "invalid block data string: no colon",
			blockData:    BlockData{},
			blockDatastr: "1100/315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd32",
			wantErr:      true,
		},
		{
			name:         "invalid block data string: no slash",
			blockData:    BlockData{},
			blockDatastr: "1:100315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd32",
			wantErr:      true,
		},
		{
			name:         "invalid block data string: non-integer block number",
			blockData:    BlockData{},
			blockDatastr: "1:abc:315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, bd, err := StringToBlockData(tt.blockDatastr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.blockData, bd)
			}
		})
	}
}

func TestStringToOwnershipData(t *testing.T) {
	type Result struct {
		nft   Nft
		owner ctypes.HexAddressString
		err   bool
	}

	testCases := []struct {
		name       string
		voteString string
		expected   struct {
			nft   Nft
			owner ctypes.HexAddressString
			err   bool
		}
	}{
		{
			name:       "Valid vote string",
			voteString: "1/0x123/0x0:0x777",
			expected: Result{
				nft: Nft{
					ChainId:      "1",
					ContractAddr: "0x123",
					TokenId:      "0x0",
				},
				owner: ctypes.HexAddressString("0x777"),
				err:   false,
			},
		},
		{
			name:       "Valid vote string",
			voteString: "1/123/0a:777",
			expected: Result{
				nft: Nft{
					ChainId:      "1",
					ContractAddr: "123",
					TokenId:      "0a",
				},
				owner: ctypes.HexAddressString("777"),
				err:   false,
			},
		},
		{
			name:       "Invalid nftId - no colone",
			voteString: "1/0x123/0x0/0x777",
			expected: Result{
				nft:   Nft{},
				owner: ctypes.HexAddressString(""),
				err:   true,
			},
		},
		{
			name:       "Invalid nftId - no slash",
			voteString: "10x123/0x0:0x777",
			expected: Result{
				nft:   Nft{},
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
