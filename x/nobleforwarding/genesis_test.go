package nobleforwarding_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/settlus/chain/testutil/keeper"
	"github.com/settlus/chain/testutil/nullify"
	"github.com/settlus/chain/x/nobleforwarding"
	"github.com/settlus/chain/x/nobleforwarding/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.NobleForwardingKeeper(t)
	nobleforwarding.InitGenesis(ctx, *k, genesisState)
	got := nobleforwarding.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
