package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"

	"github.com/settlus/chain/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) Prevote(goCtx context.Context, prevote *types.MsgPrevote) (*types.MsgPrevoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	period := m.GetParams(ctx).VotePeriod
	if types.IsLastBlockOfVotePeriod(ctx, period) {
		return nil, errorsmod.Wrapf(types.ErrPrevotesNotAccepted, "current block height (%d) is the last block of the period", ctx.BlockHeight())
	}

	// validator and feeder addresses are validated in the ValidateBasic()
	if v, err := m.ValidateFeeder(ctx, prevote.Feeder, prevote.Validator); err != nil && v {
		return nil, err
	}

	aggregatePrevote := types.AggregatePrevote{
		Hash:        prevote.Hash,
		Voter:       prevote.Validator,
		SubmitBlock: uint64(ctx.BlockHeight()),
	}
	m.SetAggregatePrevote(ctx, aggregatePrevote)

	if err := ctx.EventManager().EmitTypedEvent(&types.EventPrevote{
		Feeder:    prevote.Feeder,
		Validator: prevote.Validator,
		Hash:      prevote.Hash,
	}); err != nil {
		return nil, fmt.Errorf("failed to emit event (%s)", err)
	}

	return &types.MsgPrevoteResponse{}, nil
}

func (m msgServer) Vote(goCtx context.Context, vote *types.MsgVote) (*types.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validator and feeder addresses are validated in the ValidateBasic()
	if v, err := m.ValidateFeeder(ctx, vote.Feeder, vote.Validator); err != nil && v {
		return nil, err
	}

	// Check if the prevote exists
	aggregatePrevote := m.GetAggregatePrevote(ctx, vote.Validator)
	if aggregatePrevote == nil {
		return nil, fmt.Errorf("aggregate prevote not found")
	}

	// Check if the msg is submitted in the proper period
	params := m.GetParams(ctx)

	if params.VotePeriod == 0 {
		return nil, types.ErrVotePeriodIsZero
	}

	if (uint64(ctx.BlockHeight())/params.VotePeriod)-(aggregatePrevote.SubmitBlock/params.VotePeriod) != 1 {
		return nil, types.ErrRevealPeriodMissMatch
	}

	// Check if the vote matches the aggregate prevote hash
	hash, err := types.GetAggregateVoteHash(vote.BlockDataString, vote.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregate vote hash (%s)", err)
	}
	if aggregatePrevote.Hash != hash {
		return nil, errorsmod.Wrapf(types.ErrInvalidVote, "hash submitted in prevote (%s) does not match the hash of the vote (%s)", aggregatePrevote.Hash, hash)
	}

	// Check if all whitelisted chains are included in the vote
	blockData, err := types.ParseBlockDataString(vote.BlockDataString)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidVote, "failed to parse block data string (%s)", err)
	}
	var whiteListedChainIds []string
	for _, chain := range params.GetWhitelist() {
		whiteListedChainIds = append(whiteListedChainIds, chain.ChainId)
	}
	var voteChainIds []string
	for _, chain := range blockData {
		voteChainIds = append(voteChainIds, chain.ChainId)
	}
	if !slices.Equal(whiteListedChainIds, voteChainIds) {
		return nil, errorsmod.Wrapf(types.ErrInvalidVote, "chain ids are not matched with the whitelist")
	}

	aggregateVote := types.AggregateVote{
		BlockData: blockData,
		Voter:     vote.Validator,
	}

	m.SetAggregateVote(ctx, aggregateVote)
	m.DeleteAggregatePrevote(ctx, vote.Validator)

	err = ctx.EventManager().EmitTypedEvent(&types.EventVote{
		Feeder:    vote.Feeder,
		Validator: vote.Validator,
		BlockData: blockData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to emit event (%s)", err)
	}

	return &types.MsgVoteResponse{}, nil
}

func (m msgServer) FeederDelegationConsent(goCtx context.Context, consent *types.MsgFeederDelegationConsent) (*types.MsgFeederDelegationConsentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validator and feeder addresses are validated in the ValidateBasic()
	validator, _ := sdk.ValAddressFromBech32(consent.Validator)
	if val, found := m.StakingKeeper.GetValidator(ctx, validator); !found || !val.IsBonded() {
		return nil, errorsmod.Wrapf(types.ErrValidatorNotFound, "validator %s is not active", validator.String())
	}

	err := m.SetFeederDelegation(ctx, consent.Validator, consent.FeederAddress)
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventFeederDelegationConsent{
		Feeder:    consent.FeederAddress,
		Validator: consent.Validator,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to emit event (%s)", err)
	}

	return &types.MsgFeederDelegationConsentResponse{}, nil
}
