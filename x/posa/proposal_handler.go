package poa

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/settlus/chain/x/posa/keeper"
	"github.com/settlus/chain/x/posa/types"
)

func NewPoSAProposalHandler(cdc codec.Codec, k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CreateValidatorProposal:
			return handleCreateValidatorProposal(ctx, cdc, k, c)

		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized create validator proposal content type: %T", c)
		}
	}
}

// handleCreateValidatorProposal is a handler for executing a passed create validator proposal
func handleCreateValidatorProposal(ctx sdk.Context, cdc codec.Codec, k *keeper.Keeper, p *types.CreateValidatorProposal) error {
	amount, err := sdk.ParseCoinNormalized(p.Info.Amount)
	if err != nil {
		return err
	}

	minSelfDelegation, ok := sdk.NewIntFromString(p.Info.MinSelfDelegation)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	maxDelegation, ok := sdk.NewIntFromString(p.Info.MaxDelegation)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max delegation must be a positive integer or zero")
	}

	var pk cryptotypes.PubKey
	if err := cdc.UnmarshalInterfaceJSON([]byte(p.Info.Pubkey), &pk); err != nil {
		return err
	}

	commission := stakingtypes.CommissionRates{
		Rate:          sdk.NewDec(0),
		MaxRate:       sdk.NewDec(0),
		MaxChangeRate: sdk.NewDec(0),
	}

	validatorDescription := stakingtypes.NewDescription(
		p.Info.Moniker,
		"",
		"",
		"",
		"",
	)

	acc, err := sdk.AccAddressFromBech32(p.Info.DelegatorAddress)
	if err != nil {
		return err
	}
	
	valAddr := sdk.ValAddress(acc)

	newMsg, err := stakingtypes.NewMsgCreateValidator(
		valAddr,
		pk,
		amount,
		validatorDescription,
		commission,
		minSelfDelegation,
		maxDelegation,
		p.Info.IsProbono,
	)

	if err != nil {
		return err
	}

	err = k.CreateValidator(ctx, newMsg)
	if err != nil {
		return err
	}

	logger := k.Logger(ctx)
	logger.Info("Created validator by proposal")

	return nil
}
