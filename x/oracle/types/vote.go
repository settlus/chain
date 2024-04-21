package types

import (
	"crypto/sha256"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	Weight  int64
	Miss    bool
	Abstain bool
}

// GetAggregateVoteHash returns the hash of the aggregate vote from the block data string and salt
func GetAggregateVoteHash(voteData []*VoteData, salt string) (string, error) {
	var sb strings.Builder
	sb.WriteString(salt)
	for _, vd := range voteData {
		for _, d := range vd.Data {
			sb.WriteString(d)
		}
	}
	h := sha256.New()
	if _, err := h.Write([]byte(sb.String())); err != nil {
		return "", err
	}
	bs := h.Sum(nil)
	return fmt.Sprintf("%X", bs), nil
}

func BlockDataToVoteData(bd *BlockData) []*VoteData {
	return []*VoteData{
		{
			Topic: OralceTopic_Block,
			Data:  []string{fmt.Sprintf("%s:%d/%s", bd.ChainId, bd.BlockNumber, bd.BlockHash)},
		},
	}
}

// IsLastBlockOfSlashWindow returns true if we are at the last block of the slash slashWindow
func IsLastBlockOfSlashWindow(ctx sdk.Context, slashWindow uint64) bool {
	if slashWindow == 0 {
		return false
	}

	return (uint64)(ctx.BlockHeight())%slashWindow == 0
}
