package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/settlus/chain/x/oracle/types"
)

func TestKeeper_GetParams(t *testing.T) {
	params := types.DefaultParams()
	err := s.app.OracleKeeper.SetParams(s.ctx, params)
	require.NoError(t, err)

	actual := s.app.OracleKeeper.GetParams(s.ctx)
	require.Equal(t, params, actual)
}

func TestKeeper_SetParams(t *testing.T) {
	params := types.DefaultParams()
	err := s.app.OracleKeeper.SetParams(s.ctx, params)
	require.NoError(t, err)

	require.Equal(t, params, s.app.OracleKeeper.GetParams(s.ctx))
}
