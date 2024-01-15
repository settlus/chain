package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	p1 := DefaultParams()
	err := p1.Validate()
	require.NoError(t, err)

	// oracle fee percentage larger than 1
	p2 := DefaultParams()
	p2.OracleFeePercentage = sdk.NewDecWithPrec(101, 2)
	err = p2.Validate()
	require.Error(t, err)
}
