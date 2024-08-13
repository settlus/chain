package oracle

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/oracle/keeper"
	"github.com/settlus/chain/x/oracle/types"
	"github.com/settlus/chain/x/oracle/voteprocessor"
)

// EndBlocker runs at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	roundInfo := k.CalculateNextRoundInfo(ctx)
	k.SetCurrentRoundInfo(ctx, roundInfo)

	params := k.GetParams(ctx)

	_, voteEnd := types.CalculateVotePeriod(ctx.BlockHeight(), params.VotePeriod)
	// Proceed only if we are at the end of the vote period
	if ctx.BlockHeight() != voteEnd {
		return
	}

	maxValidators := k.StakingKeeper.MaxValidators(ctx)
	iterator := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()

	// Add all active validators to the claim map
	i := 0
	// Build claim map over all validators in active set
	validatorClaimMap := make(map[string]types.Claim)
	for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
		validator, found := k.StakingKeeper.GetValidator(ctx, iterator.Value())
		if !found {
			continue
		}

		// Exclude inactive validator or jailed validator
		if validator.IsBonded() && !validator.IsJailed() {
			valAddr := validator.GetOperator()
			validatorClaimMap[valAddr.String()] = types.Claim{
				Weight:  validator.GetConsensusPower(k.StakingKeeper.PowerReduction(ctx)),
				Miss:    false,
				Abstain: false,
			}
			i++
		}
	}

	// calculate threshold power for a block to be considered as a winner
	// threshold = total power * params.VoteThreshold (0.5 by default)
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx), k.StakingKeeper.PowerReduction(ctx))
	voteThreshold := params.VoteThreshold
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()

	// Get all aggregate votes
	aggregateVotes := k.GetAggregateVotes(ctx)

	// Check abstain
	for _, vote := range aggregateVotes {
		if len(vote.VoteData) == 0 {
			claim := validatorClaimMap[vote.Voter]
			claim.Abstain = true
		}
	}

	vps := voteprocessor.NewSettlusVoteProcessors(k, aggregateVotes, thresholdVotes)
	for _, vp := range vps {
		vp.TallyVotes(ctx, validatorClaimMap)
	}

	// do miss counting
	for addr, claim := range validatorClaimMap {
		if claim.Miss {
			// get miss count and increase it by 1
			k.SetMissCount(ctx, addr, k.GetMissCount(ctx, addr)+1)
		}
	}

	// distribute rewards to winners
	if err := k.RewardBallotWinners(ctx, validatorClaimMap); err != nil {
		logger := k.Logger(ctx)
		logger.Error("failed to distribute rewards", "error", err)
	}

	// clear ballots
	k.ClearBallots(ctx)

	// if we are at the last block of slash window, slash validators and reset miss count
	if types.IsLastBlockOfSlashWindow(ctx, params.SlashWindow) {
		k.SlashValidatorsAndResetMissCount(ctx)
	}
}
