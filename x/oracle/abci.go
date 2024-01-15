package oracle

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/oracle/keeper"
	"github.com/settlus/chain/x/oracle/types"
)

// EndBlocker runs at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	params := k.GetParams(ctx)
	if !types.IsLastBlockOfVotePeriod(ctx, params.VotePeriod) {
		return
	}

	logger := k.Logger(ctx)

	// Build claim map over all validators in active set
	validatorClaimMap := make(map[string]types.Claim)

	whitelistChainCount := len(params.Whitelist)
	maxValidators := k.StakingKeeper.MaxValidators(ctx)
	iterator := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()

	// Add all active validators to the claim map
	i := 0
	for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
		validator, found := k.StakingKeeper.GetValidator(ctx, iterator.Value())
		if !found {
			continue
		}

		// Exclude inactive validator or jailed validator
		if validator.IsBonded() && !validator.IsJailed() {
			valAddr := validator.GetOperator()
			validatorClaimMap[valAddr.String()] = types.Claim{
				Weight:    validator.GetConsensusPower(k.StakingKeeper.PowerReduction(ctx)),
				MissCount: int64(whitelistChainCount),
				Recipient: valAddr,
				Abstain:   false,
			}
			i++
		}
	}

	// calculate threshold power for a block to be considered as a winner
	// threshold = total power * params.VoteThreshold (0.5 by default)
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx), k.StakingKeeper.PowerReduction(ctx))
	voteThreshold := params.VoteThreshold
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()

	// Organize votes by chain IDs
	votesByChainId := k.GroupVotesByChainId(ctx)
	voteResults := types.TallyVotes(votesByChainId, validatorClaimMap, thresholdVotes)

	// set block data
	for chainId, blockData := range voteResults {
		if blockData == nil {
			logger.Debug("consensus failed")
			if err := ctx.EventManager().EmitTypedEvent(&types.EventOracleConsensusFailed{
				ChainId:     chainId,
				BlockHeight: ctx.BlockHeight(),
			}); err != nil {
				panic(fmt.Errorf("failed to emit event (%s)", err))
			}

			continue
		}
		k.SetBlockData(ctx, *blockData)
	}

	// increase win count for validators who voted for the winning block data
	toleratedErrorBand := params.ToleratedErrorBand
	for chainId, vote := range votesByChainId {
		correctBlockData := voteResults[chainId]

		if correctBlockData == nil {
			continue
		}

		for _, voteData := range vote {
			claim, ok := validatorClaimMap[voteData.Voter.String()]
			if !ok {
				// TODO: handle error
				panic(voteData.Voter.String() + " not found in validatorClaimMap")
			}

			// if abstain flag is set, continue
			if claim.Abstain {
				continue
			}

			// if the validator voted abstain, set abstain flag to true and continue
			if voteData.BlockData.BlockNumber < 0 {
				claim.Abstain = true
				validatorClaimMap[voteData.Voter.String()] = claim
				continue
			}

			// if the validator did not abstain and voted for the block number outside the tolerated error band, increase miss count
			if *voteData.BlockData == *correctBlockData || (voteData.BlockData.BlockNumber > 0 && Abs(voteData.BlockData.BlockNumber-correctBlockData.BlockNumber) <= int64(toleratedErrorBand)) {
				claim.MissCount--
				validatorClaimMap[voteData.Voter.String()] = claim
			}
		}
	}

	// do miss counting
	for _, claim := range validatorClaimMap {
		if claim.MissCount > 0 {
			// get miss count and increase it by 1
			k.SetMissCount(ctx, claim.Recipient.String(), k.GetMissCount(ctx, claim.Recipient.String())+1)
		}
	}

	// distribute rewards to winners
	if err := k.RewardBallotWinners(ctx, &validatorClaimMap); err != nil {
		panic(err)
	}

	// clear ballots
	k.ClearBallots(ctx)

	// if we are at the last block of slash window, slash validators and reset miss count
	if types.IsLastBlockOfSlashWindow(ctx, params.SlashWindow) {
		k.SlashValidatorsAndResetMissCount(ctx)
	}
}

// Abs returns the absolute value of x.
// (How is this not included in the Go standard library? Math.abs() is only for float64.)
func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
