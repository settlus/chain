package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRecord{}, "settlement/MsgRecord", nil)
	cdc.RegisterConcrete(&MsgCancel{}, "settlement/MsgCancel", nil)
	cdc.RegisterConcrete(&MsgDepositToTreasury{}, "settlement/MsgDepositToTreasury", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRecord{},
		&MsgCancel{},
		&MsgDepositToTreasury{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	AminoCdc  = codec.NewAminoCodec(amino)
)
