package voteprocessor

import (
	"cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/oracle/types"
)

type IVoteProcessor interface {
	TallyVotes(ctx sdk.Context, validatorClaimMap map[string]types.Claim)
}

type VoteProcessor[Source comparable, Data comparable] struct {
	topic          types.OracleTopic
	aggregateVotes []types.AggregateVote
	thresholdVotes math.Int
	dataConverter  DataConverter[Source, Data]
	onConsensus    ConsensusHook[Source, Data]
}

func (vp *VoteProcessor[Source, Data]) TallyVotes(ctx sdk.Context, validatorClaimMap map[string]types.Claim) {
	logger := ctx.Logger()
	votes := vp.groupVotes(vp.aggregateVotes)
	voteResults := make(map[Source]Data)
	for source, votes := range votes {
		var dwwList []DataWithWeight[Data]
		for _, vote := range votes {
			dwwList = append(dwwList, DataWithWeight[Data]{
				Data:   vote.Data,
				Weight: validatorClaimMap[vote.Voter.String()].Weight,
			})
		}

		picked, ok := vp.pickMostVoted(dwwList)
		if ok {
			voteResults[source] = picked
		} else {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "consensus_failed"},
				1,
				[]metrics.Label{
					telemetry.NewLabel("topic", vp.topic.String()),
				},
			)
		}
	}

	vp.onConsensus(ctx, voteResults)

	// increase win count for validators who voted for the winning data
	for source, vote := range votes {
		majorData := voteResults[source]
		for _, voteData := range vote {
			claim, ok := validatorClaimMap[voteData.Voter.String()]
			if !ok {
				// if the validator is not in the active set or is jailed, skip
				logger.Info("validator not found in active set", "validator", voteData.Voter.String())
				continue
			}

			// if abstain flag is set, continue
			if claim.Abstain {
				continue
			}

			if voteData.Data != majorData {
				claim.Miss = claim.Miss || true
			}

			validatorClaimMap[voteData.Voter.String()] = claim
		}
	}
}

func (vp *VoteProcessor[Source, Data]) groupVotes(aggregateVotes []types.AggregateVote) map[Source][]DataWithVoter[Data] {
	groupedVotes := make(map[Source][]DataWithVoter[Data])
	for _, vote := range aggregateVotes {
		for _, vd := range vote.VoteData {
			voter, err := sdk.ValAddressFromBech32(vote.Voter)
			if err != nil {
				continue
			}

			for _, strData := range vd.Data {
				source, data, err := vp.dataConverter(strData)
				if err != nil {
					continue
				}

				groupedVotes[source] = append(groupedVotes[source], DataWithVoter[Data]{
					Data:  data,
					Voter: voter,
				})
			}
		}
	}
	return groupedVotes
}

// pickMostVotedData picks the most voted data from the given slice of data
func (vp *VoteProcessor[Source, Data]) pickMostVoted(dwwList []DataWithWeight[Data]) (ret Data, ok bool) {
	// Count votes
	voteCount := make(map[Data]int64)
	for _, dww := range dwwList {
		voteCount[dww.Data] += dww.Weight
	}

	// Filter out votes that are below the threshold
	voteCountAboveThreshold := make(map[Data]int64)
	for data, count := range voteCount {
		if math.NewInt(count).GTE(vp.thresholdVotes) {
			voteCountAboveThreshold[data] = count
		}
	}

	// Assuming the threshold is at least 50%, if there are more than one votes above the threshold, return nil
	if len(voteCountAboveThreshold) > 1 {
		return ret, false
	}

	// If there is only one vote above the threshold, return it
	for data := range voteCountAboveThreshold {
		return data, true
	}

	// If there are no votes above the threshold, return nil
	return ret, false
}
