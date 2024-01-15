package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/settlement module sentinel errors
var (
	ErrNotAuthorized          = sdkerrors.Register(ModuleName, 1100, "not authorized")
	ErrNotEnoughBalance       = sdkerrors.Register(ModuleName, 1101, "not enough balance")
	ErrInvalidTenant          = sdkerrors.Register(ModuleName, 1102, "invalid tenant")
	ErrInvalidTxId            = sdkerrors.Register(ModuleName, 1103, "invalid tx id")
	ErrInvalidAccount         = sdkerrors.Register(ModuleName, 1104, "invalid account")
	ErrNotFound               = sdkerrors.Register(ModuleName, 1105, "account not found")
	ErrInvalidRequest         = sdkerrors.Register(ModuleName, 1106, "invalid request")
	ErrInvalidChainId         = sdkerrors.Register(ModuleName, 1107, "invalid chain id")
	ErrInvalidContractAddress = sdkerrors.Register(ModuleName, 1108, "invalid contract address")
	ErrInvalidTokenId         = sdkerrors.Register(ModuleName, 1109, "invalid token id")
	ErrEVMCallFailed          = sdkerrors.Register(ModuleName, 1110, "evm call failed")
	ErrEventCreationFailed    = sdkerrors.Register(ModuleName, 1111, "failed to emit event")
	ErrDuplicateRequestId     = sdkerrors.Register(ModuleName, 1112, "duplicate request id")
	ErrUTXRNotFound           = sdkerrors.Register(ModuleName, 1113, "utxr not found")
	ErrInvalidAdmin           = sdkerrors.Register(ModuleName, 1114, "invalid admin")
	ErrCannotRemoveAdmin      = sdkerrors.Register(ModuleName, 1115, "cannot remove admin")
)
