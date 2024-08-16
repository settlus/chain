package types

import (
	"fmt"
	"strings"

	ctypes "github.com/settlus/chain/types"
)

func ValidateVoteData(voteData []*VoteData, chainList []string) bool {
	chainMap := make(map[string]struct{})
	for _, chainId := range chainList {
		chainMap[chainId] = struct{}{}
	}

	for _, vd := range voteData {
		switch vd.Topic {
		case OracleTopic_BLOCK:
			// deprecated
			return true

		case OracleTopic_OWNERSHIP:
			for _, data := range vd.Data {
				nft, _, err := StringToOwnershipData(data)
				if err != nil {
					return false
				}
				if _, ok := chainMap[nft.ChainId]; !ok {
					return false
				}
			}
		default:
			return false
		}
	}

	return true
}

func StringToOwnershipData(voteString string) (ctypes.Nft, ctypes.HexAddressString, error) {
	// voteString = chainId/contractAddr/tokenId:owner

	data := strings.Split(voteString, ":")
	nft, err := ctypes.ParseNftId(data[0])

	if !isValidHex(nft.ContractAddr.String()) || !isValidHex(nft.TokenId.String()) {
		return ctypes.Nft{}, "", fmt.Errorf("invalid nftId: %s", data[0])
	}

	owner := ctypes.NoramlizeHexAddress(data[1])

	return nft, owner, err
}

func isValidHex(s string) bool {
	if len(s) == 0 {
		return false
	}

	if len(s) > 2 && s[:2] == "0x" {
		s = s[2:]
	}

	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') {
			return false
		}
	}
	return true
}
