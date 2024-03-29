package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &GenesisState{
				Params: DefaultParams(),
				Tenants: []Tenant{
					{
						Id:           1,
						Admins:       []string{"admin1"},
						PayoutPeriod: 1,
					},
				},
				Utxrs: []UTXRWithTenantAndId{},
			},
			valid: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
