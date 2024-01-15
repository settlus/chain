package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgPrevote{}, "oracle/Prevote", nil)
	cdc.RegisterConcrete(&MsgVote{}, "oracle/Vote", nil)
	cdc.RegisterConcrete(&MsgFeederDelegationConsent{}, "oracle/FeederDelegationConsent", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPrevote{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVote{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFeederDelegationConsent{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
