package e2e

import (
	"fmt"
	"strings"

	tenderminttypes "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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

func queryTenants(endpoint string) ([]settlementtypes.TenantWithTreasury, error) {
	body, err := httpGet(fmt.Sprintf("%s/settlus/settlement/tenants", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var tenantsResp settlementtypes.QueryTenantsResponse
	if err := cdc.UnmarshalJSON(body, &tenantsResp); err != nil {
		return nil, err
	}

	return tenantsResp.Tenants, nil
}

func queryUtxr(endpoint string, tenantId string, requestId string) (*settlementtypes.UTXR, error) {
	body, err := httpGet(fmt.Sprintf("%s/settlus/settlement/utxr/%s/%s", endpoint, tenantId, requestId))
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
