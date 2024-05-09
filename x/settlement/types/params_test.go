package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/settlus/chain/types"
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

	// empty ChainName
	p3 := DefaultParams()
	p3.SupportedChains[0].ChainName = ""
	err = p3.Validate()
	require.Error(t, err)

	// empty ChainUrl
	p4 := DefaultParams()
	p4.SupportedChains[0].ChainUrl = ""
	err = p4.Validate()
	require.Error(t, err)

	// empty ChainId
	p5 := DefaultParams()
	p5.SupportedChains[0].ChainId = ""
	err = p5.Validate()
	require.Error(t, err)

	// duplicate ChainId
	p6 := DefaultParams()
	p6.SupportedChains = append(p6.SupportedChains, &ctypes.Chain{
		ChainId:   "1",
		ChainName: "Foo",
		ChainUrl:  "http://localhost:8545",
	})
	err = p6.Validate()
	require.Error(t, err)

	// duplicate ChainName
	p7 := DefaultParams()
	p7.SupportedChains = append(p7.SupportedChains, &ctypes.Chain{
		ChainId:   "2",
		ChainName: "Ethereum",
		ChainUrl:  "http://localhost:8545",
	})
	err = p7.Validate()
	require.Error(t, err)
}
