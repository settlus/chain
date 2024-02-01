package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyAllowedAuthorities = []byte("AllowedAuthorities")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		AllowedAuthorities: []string{},
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAllowedAuthorities, &p.AllowedAuthorities, validateAllowedAuthorities),
	}
}

func validateAllowedAuthorities(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, chainId := range v {
		if chainId == "" {
			return errorsmod.Wrapf(ErrInvalidLengthParams, "empty chain id")
		}
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	authorities := make(map[string]bool)
	for _, authority := range p.AllowedAuthorities {
		if authority == "" {
			return errorsmod.Wrapf(ErrInvalidLengthParams, "empty chain id")
		}

		if _, ok := authorities[authority]; ok {
			return errorsmod.Wrapf(ErrInvalidLengthParams, fmt.Sprintf("duplicate chain id: %s", authority))
		}
		authorities[authority] = true
	}

	return nil
}
