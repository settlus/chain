// Copyright 2022 Evmos Foundation
// This file is part of the Evmos Network packages.
//
// Evmos is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Evmos packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Evmos packages. If not, see https://github.com/evmos/evmos/blob/main/LICENSE

package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
)

// GetCoinbaseAddress returns the block proposer's validator operator address.
func (k Keeper) GetCoinbaseAddress(ctx sdk.Context, proposerAddress sdk.ConsAddress) (common.Address, error) {
	validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, GetProposerAddress(ctx, proposerAddress))
	if !found {
		return common.Address{}, errorsmod.Wrapf(
			stakingtypes.ErrNoValidatorFound,
			"failed to retrieve validator from block proposer address %s",
			proposerAddress.String(),
		)
	}

	coinbase := common.BytesToAddress(validator.GetOperator())
	return coinbase, nil
}

// GetProposerAddress returns current block proposer's address when provided proposer address is empty.
func GetProposerAddress(ctx sdk.Context, proposerAddress sdk.ConsAddress) sdk.ConsAddress {
	if len(proposerAddress) == 0 {
		proposerAddress = ctx.BlockHeader().ProposerAddress
	}
	return proposerAddress
}
