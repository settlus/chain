package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ctypes "github.com/settlus/chain/types"
	"sigs.k8s.io/yaml"
)

var (
	KeyGasPrices           = []byte("GasPrices")
	KeyOracleFeePercentage = []byte("OracleFeePercentage")
	KeySupportedChains     = []byte("SupportedChains")
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
		SupportedChains: []*ctypes.Chain{
			{
				ChainId:   "1",
				ChainName: "Ethereum",
				ChainUrl:  "https://ethereum.org",
			},
		},
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyGasPrices, &p.GasPrices, validateGasPrices),
		paramtypes.NewParamSetPair(KeyOracleFeePercentage, &p.OracleFeePercentage, validateOracleFeePercentage),
		paramtypes.NewParamSetPair(KeySupportedChains, &p.SupportedChains, validateSupportedChains),
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

func validateSupportedChains(i interface{}) error {
	chains, ok := i.([]*ctypes.Chain)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	chainIds := make(map[string]bool)
	chainNames := make(map[string]bool)
	for _, chain := range chains {
		if strings.TrimSpace(chain.ChainId) == "" {
			return fmt.Errorf("empty chain id")
		}

		if strings.TrimSpace(chain.ChainName) == "" {
			return fmt.Errorf("empty chain name")
		}

		if strings.TrimSpace(chain.ChainUrl) == "" {
			return fmt.Errorf("empty chain url")
		}

		if _, ok := chainIds[chain.ChainId]; ok {
			return fmt.Errorf("duplicate chain id %s", chain.ChainId)
		}

		if _, ok := chainNames[chain.ChainName]; ok {
			return fmt.Errorf("duplicate chain name %s", chain.ChainName)
		}

		chainIds[chain.ChainId] = true
		chainNames[chain.ChainName] = true
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.OracleFeePercentage.IsNegative() || p.OracleFeePercentage.GT(sdk.OneDec()) {
		return fmt.Errorf("oracle fee percentage should be between 0 and 1")
	}

	if err := validateSupportedChains(p.SupportedChains); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
