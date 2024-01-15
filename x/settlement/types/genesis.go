package types

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:  DefaultParams(),
		Tenants: []Tenant(nil),
		Utxrs:   []UTXRWithTenantAndId(nil),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	tenants := gs.Tenants
	tenantIds := make([]uint64, len(tenants))
	for i, tenant := range tenants {
		tenantIds[i] = tenant.Id
	}

	return gs.Params.Validate()
}
