package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var (
	KeyAllowedChainIds = []byte("AllowedChainIds")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		AllowedChainIds: []string{},
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAllowedChainIds, &p.AllowedChainIds, validateAllowedChainIds),
	}
}

func validateAllowedChainIds(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, chainId := range v {
		if chainId == "" {
			return errorsmod.Wrapf(ErrInvalidChainId, "empty chain id")
		}
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	chainIds := make(map[string]bool)
	for _, chainId := range p.AllowedChainIds {
		if chainId == "" {
			return errorsmod.Wrapf(ErrInvalidChainId, "empty chain id")
		}

		if _, ok := chainIds[chainId]; ok {
			return errorsmod.Wrapf(ErrInvalidChainId, fmt.Sprintf("duplicate chain id: %s", chainId))
		}
		chainIds[chainId] = true
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
