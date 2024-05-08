package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"sigs.k8s.io/yaml"
)

var (
	KeyGasPrices           = []byte("GasPrices")
	KeyOracleFeePercentage = []byte("OracleFeePercentage")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		GasPrices: sdk.NewDecCoins(
			sdk.DecCoin{
				Denom:  "uusdc",
				Amount: sdk.NewDec(1)},
			sdk.DecCoin{
				Denom:  "setl",
				Amount: sdk.NewDecWithPrec(1, 4),
			}),
		OracleFeePercentage: sdk.NewDec(1),
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyGasPrices, &p.GasPrices, validateGasPrices),
		paramtypes.NewParamSetPair(KeyOracleFeePercentage, &p.OracleFeePercentage, validateOracleFeePercentage),
	}
}

func validateGasPrices(i interface{}) error {
	_, ok := i.(sdk.DecCoins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateOracleFeePercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() || v.GT(sdk.OneDec()) {
		return fmt.Errorf("oracle fee percentage should be between 0 and 1")
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.OracleFeePercentage.IsNegative() || p.OracleFeePercentage.GT(sdk.OneDec()) {
		return fmt.Errorf("oracle fee percentage should be between 0 and 1")
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
