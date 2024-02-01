package types

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params:     DefaultParams(),
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
