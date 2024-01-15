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

package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/evmos/encoding"
)

// NextFn is a no-op function that returns the context and no error in order to mock
// the next function in the AnteHandler chain.
//
// It can be used in unit tests when calling a decorator's AnteHandle method, e.g.
// `dec.AnteHandle(ctx, tx, false, NextFn)`
func NextFn(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
	return ctx, nil
}

// ValidateAnteForMsgs is a helper function, which takes in an AnteDecorator as well as 1 or
// more messages, builds a transaction containing these messages, and returns any error that
// the AnteHandler might return.
func ValidateAnteForMsgs(ctx sdk.Context, dec sdk.AnteDecorator, msgs ...sdk.Msg) error {
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	txBuilder := encodingConfig.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return err
	}

	tx := txBuilder.GetTx()

	// Call Ante decorator
	_, err = dec.AnteHandle(ctx, tx, false, NextFn)
	return err
}
