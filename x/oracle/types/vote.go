package types

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	Weight    int64
	MissCount int64
	Recipient sdk.ValAddress
	Abstain   bool
}

type BlockDataAndVoter struct {
	BlockData *BlockData
	Voter     sdk.ValAddress
}

type BlockDataAndWeight struct {
	BlockData *BlockData
	Weight    int64
}

// ParseBlockDataString parses a string of block data into a slice of BlockData.
// The string must be in the format of "{chain-id}:{height}:{hash},{chain-id}:{height}:{hash},..."
func ParseBlockDataString(blockDataStr string) ([]*BlockData, error) {
	var blockData []*BlockData
	for _, blockDataStr := range strings.Split(blockDataStr, ",") {
		if blockDataStr == "" {
			continue
		}
		blockDataFields := strings.Split(blockDataStr, ":")
		if len(blockDataFields) != 3 {
			return nil, fmt.Errorf("invalid block data string: %s", blockDataStr)
		}
		chainID := blockDataFields[0]
		height, err := strconv.ParseInt(blockDataFields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid block data string: %s", blockDataStr)
		}
		hash := blockDataFields[2]
		blockData = append(blockData, &BlockData{
			ChainId:     chainID,
			BlockNumber: height,
			BlockHash:   hash,
		})
	}
	return blockData, nil
}

// GetAggregateVoteHash returns the hash of the aggregate vote from the block data string and salt
func GetAggregateVoteHash(blockDataString, salt string) (string, error) {
	sourceStr := fmt.Sprintf("%s%s", salt, blockDataString)
	h := sha256.New()
	h.Write([]byte(sourceStr))
	bs := h.Sum(nil)
	return fmt.Sprintf("%X", bs), nil
}

// BlockDataToString returns a string representation of the given block data
func BlockDataToString(bd *BlockData) string {
	return fmt.Sprintf("%s:%d:%s", bd.ChainId, bd.BlockNumber, bd.BlockHash)
}

// IsLastBlockOfSlashWindow returns true if we are at the last block of the slash slashWindow
func IsLastBlockOfSlashWindow(ctx sdk.Context, slashWindow uint64) bool {
	if slashWindow == 0 {
		return false
	}

	return (uint64)(ctx.BlockHeight())%slashWindow == 0
}

// TallyVotes tallies votes by chain ID and picks the most voted block data
func TallyVotes(votesByChainId map[string][]*BlockDataAndVoter, validatorClaimMap map[string]Claim, thresholdVotes sdkmath.Int) map[string]*BlockData {
	blockData := make(map[string]*BlockData)
	for chainId, votes := range votesByChainId {
		var bds []*BlockDataAndWeight
		for _, vote := range votes {
			// Skip abstain votes
			if vote.BlockData.BlockNumber < 0 {
				continue
			}
			bds = append(bds, &BlockDataAndWeight{
				BlockData: vote.BlockData,
				Weight:    validatorClaimMap[vote.Voter.String()].Weight,
			})
		}

		blockData[chainId] = PicMostVotedBlockData(bds, thresholdVotes)
	}
	return blockData
}

// PicMostVotedBlockData picks the most voted block data from the given slice of block data
func PicMostVotedBlockData(votes []*BlockDataAndWeight, thresholdVotes sdkmath.Int) *BlockData {
	// Count votes
	voteCount := make(map[BlockData]int64)
	for _, vote := range votes {
		voteCount[*vote.BlockData] += vote.Weight
	}

	// Filter out votes that are below the threshold
	voteCountAboveThreshold := make(map[BlockData]int64)
	for bd, count := range voteCount {
		if sdkmath.NewInt(count).GTE(thresholdVotes) {
			voteCountAboveThreshold[bd] = count
		}
	}

	// Assuming the threshold is at least 50%, if there are more than one votes above the threshold, return nil
	if len(voteCountAboveThreshold) > 1 {
		return nil
	}

	// If there is only one vote above the threshold, return it
	for bd := range voteCountAboveThreshold {
		return &bd
	}

	// If there are no votes above the threshold, return nil
	return nil
}
