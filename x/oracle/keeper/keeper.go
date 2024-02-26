package keeper

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/settlus/chain/x/oracle/types"
)

const BlockTimestampMargin = time.Second

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace

		AccountKeeper      types.AccountKeeper
		BankKeeper         types.BankKeeper
		DistributionKeeper types.DistributionKeeper
		StakingKeeper      types.StakingKeeper

		distributionName string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distributionKeeper types.DistributionKeeper,
	stakingKeeper types.StakingKeeper,

	distributionName string,
) *Keeper {
	// ensure oracle module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,

		AccountKeeper:      accountKeeper,
		BankKeeper:         bankKeeper,
		DistributionKeeper: distributionKeeper,
		StakingKeeper:      stakingKeeper,

		distributionName: distributionName,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetCurrentRoundInfo(ctx sdk.Context) *types.RoundInfo {
	params := k.GetParams(ctx)
	blockHeight := ctx.BlockHeight()
	prevoteEnd, voteEnd := types.CalculateVotePeriod(blockHeight, params.VotePeriod)

	roundInfo := types.RoundInfo{
		Id:         types.CalculateRoundId(blockHeight, params.VotePeriod),
		PrevoteEnd: prevoteEnd,
		VoteEnd:    voteEnd,

		ChainIds:  params.GetWhitelistChainIds(),
		Timestamp: ctx.BlockHeader().Time.Add(-BlockTimestampMargin).UnixMilli(),
	}

	return &roundInfo
}

// GetBlockData returns block data of a chain
func (k Keeper) GetBlockData(ctx sdk.Context, chainId string) (*types.BlockData, error) {
	store := ctx.KVStore(k.storeKey)
	chain, err := k.GetChain(ctx, chainId)
	if err != nil {
		return nil, fmt.Errorf("chain not found: %w", err)
	}

	bd := store.Get(types.BlockDataKey(chain.ChainId))
	if bd == nil {
		return nil, types.ErrBlockDataNotFound
	}
	var blockData types.BlockData
	k.cdc.MustUnmarshal(bd, &blockData)
	return &blockData, nil
}

// GetAllBlockData returns block data of all chains
func (k Keeper) GetAllBlockData(ctx sdk.Context) []types.BlockData {
	store := ctx.KVStore(k.storeKey)
	var blockData []types.BlockData
	iterator := sdk.KVStorePrefixIterator(store, types.BlockDataKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var bd types.BlockData
		k.cdc.MustUnmarshal(iterator.Value(), &bd)
		blockData = append(blockData, bd)
	}
	return blockData
}

// SetBlockData sets block data of a chain
func (k Keeper) SetBlockData(ctx sdk.Context, blockData types.BlockData) {
	store := ctx.KVStore(k.storeKey)
	bd := k.cdc.MustMarshal(&blockData)
	store.Set(types.BlockDataKey(blockData.ChainId), bd)
}

// GetFeederDelegation returns feeder delegation of a validator
func (k Keeper) GetFeederDelegation(ctx sdk.Context, validatorAddr string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	fd := store.Get(types.FeederDelegationKey(validatorAddr))
	if fd == nil {
		// By default, validator is its own feeder
		validator, _ := sdk.ValAddressFromBech32(validatorAddr)
		return sdk.AccAddress(validator.Bytes())
	}
	return fd
}

// GetFeederDelegations returns feeder delegations of all validators
func (k Keeper) GetFeederDelegations(ctx sdk.Context) []types.FeederDelegation {
	store := ctx.KVStore(k.storeKey)
	var feederDelegations []types.FeederDelegation
	iterator := sdk.KVStorePrefixIterator(store, types.FeederDelegationKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validatorAddr := string(iterator.Key()[1:])
		feederAddr := string(iterator.Value())
		feederDelegations = append(feederDelegations, types.FeederDelegation{
			ValidatorAddress: validatorAddr,
			FeederAddress:    feederAddr,
		})
	}

	return feederDelegations
}

// SetFeederDelegation sets feeder delegation of a validator
func (k Keeper) SetFeederDelegation(ctx sdk.Context, validatorAddr string, feederAddr string) error {
	feeder, err := sdk.AccAddressFromBech32(feederAddr)
	if err != nil {
		return errorsmod.Wrapf(types.ErrInvalidFeeder, "feeder address %s is invalid", feederAddr)
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.FeederDelegationKey(validatorAddr), feeder.Bytes())
	return nil
}

// ValidateFeeder validates if the feeder has permission to vote for the validator
func (k Keeper) ValidateFeeder(ctx sdk.Context, feederAddr string, validatorAddr string) (bool, error) {
	validator, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return false, errorsmod.Wrapf(types.ErrInvalidValidator, "validator address %s is invalid", validatorAddr)
	}

	if val, found := k.StakingKeeper.GetValidator(ctx, validator); !found || !val.IsBonded() {
		return false, errorsmod.Wrapf(types.ErrValidatorNotFound, "validator %s is not active", validator.String())
	}

	feeder, err := sdk.AccAddressFromBech32(feederAddr)
	if err != nil {
		return false, errorsmod.Wrapf(types.ErrInvalidFeeder, "feeder address %s is invalid", feederAddr)
	}

	if !feeder.Equals(validator) {
		delegation := k.GetFeederDelegation(ctx, validatorAddr)
		if !delegation.Equals(feeder) {
			return false, errorsmod.Wrapf(types.ErrNoVotingPermission, "feeder %s has no permission to vote for validator %s", feeder, validator)
		}
	}

	return true, nil
}

// GetMissCount returns miss count of a validator
func (k Keeper) GetMissCount(ctx sdk.Context, validatorAddr string) uint64 {
	store := ctx.KVStore(k.storeKey)
	mc := store.Get(types.MissCountKey(validatorAddr))
	if mc == nil {
		return 0
	}
	return sdk.BigEndianToUint64(mc)
}

// GetMissCounts returns miss counts of all validators
func (k Keeper) GetMissCounts(ctx sdk.Context) []types.MissCount {
	var missCounts []types.MissCount
	k.IterateMissCount(ctx, func(validatorAddr string, missCount uint64) (stop bool) {
		missCounts = append(missCounts, types.MissCount{
			ValidatorAddress: validatorAddr,
			MissCount:        missCount,
		})
		return false
	})
	return missCounts
}

// SetMissCount sets miss count of a validator
func (k Keeper) SetMissCount(ctx sdk.Context, validatorAddr string, missCount uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MissCountKey(validatorAddr), sdk.Uint64ToBigEndian(missCount))
}

func (k Keeper) DeleteMissCount(ctx sdk.Context, validatorAddr string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.MissCountKey(validatorAddr))
}

func (k Keeper) IterateMissCount(ctx sdk.Context, handler func(validatorAddr string, missCount uint64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.MissCountKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		validatorAddr := string(iter.Key()[1:])
		missCount := sdk.BigEndianToUint64(iter.Value())

		if handler(validatorAddr, missCount) {
			break
		}
	}
}

// GetRewardPool returns the current reward pool balance
func (k Keeper) GetRewardPool(ctx sdk.Context) sdk.Coins {
	addr := k.AccountKeeper.GetModuleAddress(types.ModuleName)
	if addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return k.BankKeeper.GetAllBalances(ctx, addr)
}

/********************* Prevote *********************/

// GetAggregatePrevote returns aggregate prevote of a validator
func (k Keeper) GetAggregatePrevote(ctx sdk.Context, validatorAddr string) *types.AggregatePrevote {
	store := ctx.KVStore(k.storeKey)
	ap := store.Get(types.AggregatePrevoteKey(validatorAddr))
	if ap == nil {
		return nil
	}
	var aggregatePrevote types.AggregatePrevote
	k.cdc.MustUnmarshal(ap, &aggregatePrevote)
	return &aggregatePrevote
}

// GetAggregatePrevotes returns aggregate prevotes of all validators
func (k Keeper) GetAggregatePrevotes(ctx sdk.Context) []types.AggregatePrevote {
	store := ctx.KVStore(k.storeKey)
	var aggregatePrevotes []types.AggregatePrevote
	iterator := sdk.KVStorePrefixIterator(store, types.AggregatePrevoteKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ap types.AggregatePrevote
		k.cdc.MustUnmarshal(iterator.Value(), &ap)
		aggregatePrevotes = append(aggregatePrevotes, ap)
	}
	return aggregatePrevotes
}

// GetAggregatePrevoteStore returns a new KV store from aggregate prevote prefix
func (k Keeper) GetAggregatePrevoteStore(ctx sdk.Context) sdk.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.AggregatePrevoteKeyPrefix)
}

// SetAggregatePrevote sets aggregate prevote of a validator
func (k Keeper) SetAggregatePrevote(ctx sdk.Context, aggregatePrevote types.AggregatePrevote) {
	store := ctx.KVStore(k.storeKey)
	ap := k.cdc.MustMarshal(&aggregatePrevote)
	store.Set(types.AggregatePrevoteKey(aggregatePrevote.Voter), ap)
}

// DeleteAggregatePrevote deletes aggregate prevote of a validator
func (k Keeper) DeleteAggregatePrevote(ctx sdk.Context, validatorAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AggregatePrevoteKey(validatorAddress))
}

// IterateAggregatePrevotes iterates over prevotes in the store
func (k Keeper) IterateAggregatePrevotes(ctx sdk.Context, handler func(voterAddr string, aggregatePrevote types.AggregatePrevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AggregatePrevoteKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		voterAddr := string(iter.Key()[1:])
		var aggregatePrevote types.AggregatePrevote
		k.cdc.MustUnmarshal(iter.Value(), &aggregatePrevote)

		if handler(voterAddr, aggregatePrevote) {
			break
		}
	}
}

/********************* Vote *********************/

// GetAggregateVote returns aggregate vote of a validator
func (k Keeper) GetAggregateVote(ctx sdk.Context, validatorAddr string) *types.AggregateVote {
	store := ctx.KVStore(k.storeKey)
	av := store.Get(types.AggregateVoteKey(validatorAddr))
	if av == nil {
		return nil
	}
	var aggregateVote types.AggregateVote
	k.cdc.MustUnmarshal(av, &aggregateVote)
	return &aggregateVote
}

// GetAggregateVotes returns aggregate votes of all validators
func (k Keeper) GetAggregateVotes(ctx sdk.Context) []types.AggregateVote {
	store := ctx.KVStore(k.storeKey)
	var aggregateVotes []types.AggregateVote
	iterator := sdk.KVStorePrefixIterator(store, types.AggregateVoteKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var av types.AggregateVote
		k.cdc.MustUnmarshal(iterator.Value(), &av)
		aggregateVotes = append(aggregateVotes, av)
	}
	return aggregateVotes
}

// GetAggregateVoteStore returns a new KV store from aggregate vote prefix
func (k Keeper) GetAggregateVoteStore(ctx sdk.Context) sdk.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.AggregateVoteKeyPrefix)
}

// GroupVotesByChainId groups votes by chain ID
func (k Keeper) GroupVotesByChainId(ctx sdk.Context) map[string][]*types.BlockDataAndVoter {
	// Get all aggregate votes
	var aggregateVotes []types.AggregateVote
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AggregateVoteKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var av types.AggregateVote
		k.cdc.MustUnmarshal(iterator.Value(), &av)
		aggregateVotes = append(aggregateVotes, av)
	}

	// Group votes by chain ID
	groupedVotes := make(map[string][]*types.BlockDataAndVoter)
	for _, vote := range aggregateVotes {
		for _, bd := range vote.BlockData {
			voter, _ := sdk.ValAddressFromBech32(vote.Voter)
			groupedVotes[bd.ChainId] = append(groupedVotes[bd.ChainId], &types.BlockDataAndVoter{
				BlockData: bd,
				Voter:     voter,
			})
		}
	}
	return groupedVotes
}

// SetAggregateVote sets aggregate vote of a validator
func (k Keeper) SetAggregateVote(ctx sdk.Context, aggregateVote types.AggregateVote) {
	store := ctx.KVStore(k.storeKey)
	av := k.cdc.MustMarshal(&aggregateVote)
	store.Set(types.AggregateVoteKey(aggregateVote.Voter), av)
}

// DeleteAggregateVote deletes aggregate vote of a validator
func (k Keeper) DeleteAggregateVote(ctx sdk.Context, validator string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AggregateVoteKey(validator))
}

// IterateAggregateVotes iterates rate over votes in the store
func (k Keeper) IterateAggregateVotes(ctx sdk.Context, handler func(voterAddr string, aggregateVote types.AggregateVote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AggregateVoteKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		voterAddr := string(iter.Key()[1:])
		var aggregateVote types.AggregateVote
		k.cdc.MustUnmarshal(iter.Value(), &aggregateVote)

		if handler(voterAddr, aggregateVote) {
			break
		}
	}
}

/********************* Ballots & Validators *********************/

// RewardBallotWinners distributes rewards to validators who voted for the winning block data
func (k Keeper) RewardBallotWinners(ctx sdk.Context, validatorClaimMap *map[string]types.Claim) error {
	weightSum := int64(0)
	for _, claim := range *validatorClaimMap {
		if claim.MissCount == 0 && !claim.Abstain {
			weightSum += claim.Weight
		}
	}

	if weightSum == 0 {
		return nil
	}

	// distribute rewards in proportion to the voting power
	rewards := sdk.NewDecCoinsFromCoins(k.GetRewardPool(ctx)...)
	if rewards.IsZero() {
		return nil
	}

	logger := k.Logger(ctx)
	logger.Debug("RewardBallotWinner", "rewards", rewards)

	var distributedReward sdk.Coins

	for _, voter := range *validatorClaimMap {
		// skip if the validator abstained or missed the vote
		if voter.Abstain || voter.MissCount > 0 {
			logger.Debug(fmt.Sprintf("no reward %s(%s)",
				voter.Recipient.String(),
				voter.Recipient.String()),
				"miss count in the current ballot", voter.MissCount,
				"abstain", voter.Abstain,
			)
			continue
		}

		// multiply the reward by the weight of the validator
		// rewardCoins = rewards * (weight / weightSum)
		rewardCoins, _ := rewards.MulDec(sdk.NewDec(voter.Weight).QuoInt64(weightSum)).TruncateDecimal()

		// distribute reward to the validator
		receiverVal, ok := k.StakingKeeper.GetValidator(ctx, voter.Recipient)
		if !ok {
			return fmt.Errorf("validator not found: %s", voter.Recipient)
		}

		if !rewardCoins.IsZero() {
			k.DistributionKeeper.AllocateTokensToValidator(ctx, receiverVal, sdk.NewDecCoinsFromCoins(rewardCoins...))
			distributedReward = distributedReward.Add(rewardCoins...)
		} else {
			logger.Debug(fmt.Sprintf("no reward %s(%s)",
				receiverVal.GetMoniker(),
				receiverVal.GetOperator().String()),
				"weight", voter.Weight,
			)
		}
	}

	// Move distributed reward to distribution module
	if err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distributionName, distributedReward); err != nil {
		return fmt.Errorf("failed to move distributed reward to distribution module: %w", err)
	}

	return nil
}

func (k Keeper) ClearBallots(ctx sdk.Context) {
	// Clear all aggregate prevotes that are older than the current block height
	k.IterateAggregatePrevotes(ctx, func(validatorAddress string, aggregatePrevote types.AggregatePrevote) (stop bool) {
		if ctx.BlockHeight() >= int64(aggregatePrevote.SubmitBlock) {
			k.DeleteAggregatePrevote(ctx, validatorAddress)
		}

		return false
	})

	// Clear all aggregate votes
	k.IterateAggregateVotes(ctx, func(validatorAddress string, aggregateVote types.AggregateVote) (stop bool) {
		k.DeleteAggregateVote(ctx, validatorAddress)
		return false
	})
}

// SlashValidatorsAndResetMissCount slashes validators who missed the vote more than MaxMissCountPerSlashWindow
// and resets miss count of all validators
func (k Keeper) SlashValidatorsAndResetMissCount(ctx sdk.Context) {
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()
	distributionHeight := height - sdk.ValidatorUpdateDelay - 1
	powerReduction := k.StakingKeeper.PowerReduction(ctx)
	params := k.GetParams(ctx)
	MaxMissCountPerSlashWindow := params.MaxMissCountPerSlashWindow
	slashFraction := params.SlashFraction

	k.IterateMissCount(ctx, func(operatorAddress string, missCount uint64) (stop bool) {
		if missCount > MaxMissCountPerSlashWindow {
			// slash validator if the validator missed the vote more than MaxMissCountPerSlashWindow
			validatorAddress, err := sdk.ValAddressFromBech32(operatorAddress)
			if err != nil {
				panic(fmt.Errorf("failed to parse validator address from store: %w", err))
			}
			validator, ok := k.StakingKeeper.GetValidator(ctx, validatorAddress)
			if !ok {
				logger.Debug(fmt.Sprintf("validator not found: %s", operatorAddress))
			}

			if validator.IsBonded() && !validator.IsJailed() {
				consAddr, err := validator.GetConsAddr()
				if err != nil {
					panic(fmt.Errorf("failed to get consensus address from validator: %w", err))
				}

				k.StakingKeeper.Slash(
					ctx, consAddr,
					distributionHeight, validator.GetConsensusPower(powerReduction), slashFraction,
				)
				k.StakingKeeper.Jail(ctx, consAddr)
			}
		}

		k.DeleteMissCount(ctx, operatorAddress)
		return false
	})
}
