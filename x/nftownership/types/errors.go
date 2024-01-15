package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// ErrABIPack x/nftownership module sentinel errors
var (
	ErrABIPack        = errorsmod.Register(ModuleName, 9, "contract ABI pack failed")
	ErrInvalidChainId = errorsmod.Register(ModuleName, 10, "invalid chain id")
	ErrEVMCallFailed  = errorsmod.Register(ModuleName, 11, "evm call failed")
)
