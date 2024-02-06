package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/nobleforwarding module sentinel errors
var (
	ErrInvalidPacket            = errorsmod.Register(ModuleName, 1100, "invalid packet")
	ErrInvalidPacketTimeout     = errorsmod.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion           = errorsmod.Register(ModuleName, 1501, "invalid version")
	ErrInvalidChannelCapability = errorsmod.Register(ModuleName, 1502, "invalid channel capability")
	ErrFailedToSendPacket       = errorsmod.Register(ModuleName, 1503, "failed to send packet")
	ErrInvalidPacketData        = errorsmod.Register(ModuleName, 1504, "invalid packet data")
	ErrInvalidAckFormat         = errorsmod.Register(ModuleName, 1505, "invalid acknowledgement format")
)
