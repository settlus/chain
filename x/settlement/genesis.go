package settlement

import (
	"fmt"

	"github.com/settlus/chain/x/settlement/keeper"
	"github.com/settlus/chain/x/settlement/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k *keeper.SettlementKeeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	for _, utxrWithTenantAndId := range genState.Utxrs {
		if _, err := k.CreateUTXR(ctx, utxrWithTenantAndId.TenantId, &utxrWithTenantAndId.Utxr); err != nil {
			panic(fmt.Errorf("unable to create utxr during init genesis: %w", err))
		}
	}

	for _, tenant := range genState.Tenants {
		k.CreateTreasuryAccount(ctx, tenant.Id)
		k.SetTenant(ctx, &tenant)
	}

	k.InitAccountModule(ctx)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k *keeper.SettlementKeeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.Utxrs = k.GetAllUTXRWithTenantAndID(ctx)
	genesis.Tenants = k.GetAllTenants(ctx)

	return genesis
}
