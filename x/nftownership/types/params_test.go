package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	p1 := DefaultParams()
	err := p1.Validate()
	require.NoError(t, err)

	// empty string as chain id
	p2 := DefaultParams()
	p2.AllowedChainIds = []string{""}
	err = p2.Validate()
	require.Error(t, err)

	// duplicate chain id
	p3 := DefaultParams()
	p3.AllowedChainIds = []string{"chain1", "chain1"}
	err = p3.Validate()
	require.Error(t, err)
}
