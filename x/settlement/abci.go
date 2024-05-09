package settlement

import (
	"time"

	"github.com/settlus/chain/x/settlement/keeper"
	"github.com/settlus/chain/x/settlement/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlock(ctx sdk.Context, k *keeper.SettlementKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), "begin_blocker")

	for _, tenant := range k.GetAllTenants(ctx) {
		k.Settle(ctx, tenant.Id)
	}
}
