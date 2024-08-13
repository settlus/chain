package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/oracle module sentinel errors
var (
	ErrChainNotFound         = errorsmod.Register(ModuleName, 1000, "chain not found")
	ErrNoVotingPermission    = errorsmod.Register(ModuleName, 1002, "no voting permission")
	ErrValidatorNotFound     = errorsmod.Register(ModuleName, 1003, "invalid validator")
	ErrRevealPeriodMissMatch = errorsmod.Register(ModuleName, 1004, "reveal period of submitted vote do not match with registered prevote")
	ErrInvalidVote           = errorsmod.Register(ModuleName, 1005, "invalid vote")
	ErrVotePeriodIsZero      = errorsmod.Register(ModuleName, 1006, "vote period is zero")
	ErrPrevotesNotAccepted   = errorsmod.Register(ModuleName, 1007, "prevotes are not accepted in this period")
	ErrInvalidParams         = errorsmod.Register(ModuleName, 1008, "invalid params")
	ErrInvalidValidator      = errorsmod.Register(ModuleName, 1009, "invalid validator")
	ErrInvalidFeeder         = errorsmod.Register(ModuleName, 1010, "invalid feeder")
)
