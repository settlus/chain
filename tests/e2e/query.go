package e2e

import (
	"errors"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	tenderminttypes "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/evmos/evmos/v19/x/evm/types"

	settlementtypes "github.com/settlus/chain/x/settlement/types"
)

// if coin is zero, return empty coin.
func getSpecificBalance(endpoint, addr, denom string) (amt sdk.Coin, err error) {
	balances, err := querySettlusAllBalances(endpoint, addr)
	if err != nil {
		return amt, err
	}
	for _, c := range balances {
		if strings.Contains(c.Denom, denom) {
			amt = c
			break
		}
	}
	return amt, nil
}

func getEthBalance(endpoint, addr string) (uint64, error) {
	body, err := httpGet(fmt.Sprintf("%s/evmos/evm/v1/balances/%s", endpoint, addr))
	if err != nil {
		return 0, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var res evmtypes.QueryBalanceResponse
	if err := cdc.UnmarshalJSON(body, &res); err != nil {
		return 0, err
	}

	val, ok := sdkmath.NewIntFromString(res.Balance)
	if !ok {
		return 0, errors.New("invalid balance")
	}

	return val.Uint64(), nil
}

func queryTenants(endpoint string) ([]settlementtypes.TenantWithTreasury, error) {
	body, err := httpGet(fmt.Sprintf("%s/settlus/settlement/v1alpha1/tenants", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	fmt.Println(string(body))
	var tenantsResp settlementtypes.QueryTenantsResponse
	if err := cdc.UnmarshalJSON(body, &tenantsResp); err != nil {
		return nil, err
	}

	return tenantsResp.Tenants, nil
}

func queryTenant(endpoint string, tenantId uint64) (*settlementtypes.TenantWithTreasury, error) {
	body, err := httpGet(fmt.Sprintf("%s/settlus/settlement/v1alpha1/tenant/%d", endpoint, tenantId))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var tenantResp settlementtypes.QueryTenantResponse
	if err := cdc.UnmarshalJSON(body, &tenantResp); err != nil {
		return nil, err
	}

	return &tenantResp.Tenant, nil
}

func queryUtxr(endpoint string, tenantId uint64, requestId string) (*settlementtypes.UTXR, error) {
	body, err := httpGet(fmt.Sprintf("%s/settlus/settlement/v1alpha1/utxr/%d/%s", endpoint, tenantId, requestId))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var utxrResp settlementtypes.QueryUTXRResponse
	if err := cdc.UnmarshalJSON(body, &utxrResp); err != nil {
		return nil, err
	}

	return &utxrResp.Utxr, nil
}

func querySettlusAllBalances(endpoint, addr string) (sdk.Coins, error) {
	body, err := httpGet(fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", endpoint, addr))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var balancesResp banktypes.QueryAllBalancesResponse
	if err := cdc.UnmarshalJSON(body, &balancesResp); err != nil {
		return nil, err
	}

	return balancesResp.Balances, nil
}

func queryLatestBlockId(endpoint string) (string, error) {
	body, err := httpGet(fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/blocks/latest", endpoint))
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var blockResp tenderminttypes.GetLatestBlockResponse
	if err := cdc.UnmarshalJSON(body, &blockResp); err != nil {
		return "", err
	}

	return blockResp.BlockId.String(), nil
}
