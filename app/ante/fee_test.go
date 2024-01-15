package ante

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Test_CalculateFees(t *testing.T) {
	tests := []struct {
		name                string
		oracleFeePercentage sdk.Dec
		fees                sdk.Coins
		expectedGasFees     sdk.Coins
		expectedOracleFees  sdk.Coins
	}{
		{
			name:                "50/50 split",
			oracleFeePercentage: sdk.NewDecWithPrec(50, 2),
			fees:                sdk.NewCoins(sdk.NewInt64Coin("setl", 100)),
			expectedGasFees:     sdk.NewCoins(sdk.NewInt64Coin("setl", 50)),
			expectedOracleFees:  sdk.NewCoins(sdk.NewInt64Coin("setl", 50)),
		}, {
			name:                "100% oracle fee",
			oracleFeePercentage: sdk.NewDecWithPrec(100, 2),
			fees:                sdk.NewCoins(sdk.NewInt64Coin("setl", 100)),
			expectedGasFees:     sdk.NewCoins(),
			expectedOracleFees:  sdk.NewCoins(sdk.NewInt64Coin("setl", 100)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gasFees, oracleFees := CalculateFees(tt.oracleFeePercentage, tt.fees)
			if !gasFees.IsEqual(tt.expectedGasFees) {
				t.Errorf("expected gas fees %s, got %s", tt.expectedGasFees, gasFees)
			}
			if !oracleFees.IsEqual(tt.expectedOracleFees) {
				t.Errorf("expected oracle fees %s, got %s", tt.expectedOracleFees, oracleFees)
			}
		})
	}
}
