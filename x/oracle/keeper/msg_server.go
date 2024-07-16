package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	roundInfo := m.GetCurrentRoundInfo(ctx)

	if roundInfo == nil || roundInfo.Id != prevote.RoundId {
		return nil, errorsmod.Wrapf(types.ErrPrevotesNotAccepted, "round info not found for the round id (%d)", prevote.RoundId)
	}

	if ctx.BlockHeight() > roundInfo.PrevoteEnd {
		return nil, errorsmod.Wrapf(types.ErrPrevotesNotAccepted, "prevote period is over")
	}

	aggregatePrevote := types.AggregatePrevote{
		Hash:  prevote.Hash,
		Voter: prevote.Validator,
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

	roundInfo := m.GetCurrentRoundInfo(ctx)
	if roundInfo == nil || roundInfo.Id != vote.RoundId {
		return nil, errorsmod.Wrapf(types.ErrPrevotesNotAccepted, "round info not found for the round id (%d)", vote.RoundId)
	}

	if ctx.BlockHeight() > roundInfo.VoteEnd {
		return nil, errorsmod.Wrapf(types.ErrPrevotesNotAccepted, "vote period is over")
	}

	if !types.ValidateVoteData(vote.VoteData, m.SettlementKeeper.GetSupportedChainIds(ctx)) {
		return nil, errorsmod.Wrapf(types.ErrInvalidVote, "invalid vote data")
	}

	// Check if the prevote exists
	aggregatePrevote := m.GetAggregatePrevote(ctx, vote.Validator)
	if aggregatePrevote == nil {
		return nil, fmt.Errorf("aggregate prevote not found")
	}

	// Check if the vote matches the aggregate prevote hash
	hash, err := types.GetAggregateVoteHash(vote.VoteData, vote.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregate vote hash (%s)", err)
	}
	if aggregatePrevote.Hash != hash {
		return nil, errorsmod.Wrapf(types.ErrInvalidVote, "hash submitted in prevote (%s) does not match the hash of the vote (%s)", aggregatePrevote.Hash, hash)
	}

	aggregateVote := types.AggregateVote{
		VoteData: vote.VoteData,
		Voter:    vote.Validator,
	}

	m.SetAggregateVote(ctx, aggregateVote)
	m.DeleteAggregatePrevote(ctx, vote.Validator)

	err = ctx.EventManager().EmitTypedEvent(&types.EventVote{
		Feeder:    vote.Feeder,
		Validator: vote.Validator,
		VoteData:  vote.VoteData,
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
