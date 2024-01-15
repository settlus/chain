package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/settlement/types"
)

type SettlementKeeper interface {
	GetParams(ctx sdk.Context) (params types.Params)
}
