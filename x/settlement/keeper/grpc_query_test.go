package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	settlustypes "github.com/settlus/chain/types"
	"github.com/settlus/chain/x/settlement/types"
)

func (suite *SettlementTestSuite) TestKeeper_Params() {
	params := types.DefaultParams()
	s.keeper.SetParams(s.ctx, params)
	paramsResponse, err := s.app.SettlementKeeper.Params(s.ctx, &types.QueryParamsRequest{})
	suite.NoError(err)
	suite.Equal(types.DefaultParams(), paramsResponse.Params)
}

func (suite *SettlementTestSuite) TestKeeper_UTXR() {
	// Create UTXRs
	for i := 1; i < 4; i++ {
		utxrId, err := s.keeper.CreateUTXR(s.ctx, uint64(1), &types.UTXR{
			RequestId:   fmt.Sprintf("request-%d", i),
			Recipient:   settlustypes.NewHexAddressString(s.creator),
			Amount:      sdk.NewCoin("uusdc", sdk.NewInt(10)),
			PayoutBlock: 100,
		})
		s.NoError(err)
		s.Equal(uint64(i), utxrId)

		utxrId, err = s.keeper.CreateUTXR(s.ctx, uint64(2), &types.UTXR{
			RequestId:   fmt.Sprintf("request-%d", i),
			Recipient:   settlustypes.NewHexAddressString(s.creator),
			Amount:      sdk.NewCoin("uusdc", sdk.NewInt(10)),
			PayoutBlock: 100,
		})
		s.NoError(err)
		s.Equal(uint64(i), utxrId)
	}
	suite.Commit()

	// Get UTXRs
	for i := 1; i < 4; i++ {
		res, err := s.queryClient.UTXR(s.ctx, &types.QueryUTXRRRequest{
			TenantId:  uint64(1),
			RequestId: fmt.Sprintf("request-%d", i),
		})
		s.NoError(err)
		s.Equal(fmt.Sprintf("request-%d", i), res.Utxr.RequestId)

		res, err = s.queryClient.UTXR(s.ctx, &types.QueryUTXRRRequest{
			TenantId:  uint64(2),
			RequestId: fmt.Sprintf("request-%d", i),
		})
		s.NoError(err)
		s.Equal(fmt.Sprintf("request-%d", i), res.Utxr.RequestId)
	}
}

func (suite *SettlementTestSuite) TestKeeper_UTXRs() {
	// Create UTXRs
	for i := 1; i < 4; i++ {
		utxrId, err := s.keeper.CreateUTXR(s.ctx, 1, &types.UTXR{
			RequestId:   fmt.Sprintf("request-%d", i),
			Recipient:   settlustypes.NewHexAddressString(s.creator),
			Amount:      sdk.NewCoin("uusdc", sdk.NewInt(10)),
			PayoutBlock: 100,
		})
		s.NoError(err)
		s.Equal(uint64(i), utxrId)

		utxrId, err = s.keeper.CreateUTXR(s.ctx, 2, &types.UTXR{
			RequestId:   fmt.Sprintf("request-%d", i),
			Recipient:   settlustypes.NewHexAddressString(s.creator),
			Amount:      sdk.NewCoin("uusdc", sdk.NewInt(10)),
			PayoutBlock: 100,
		})
		s.NoError(err)
		s.Equal(uint64(i), utxrId)
	}

	// Get UTXRs
	res, err := s.queryClient.UTXRs(s.ctx, &types.QueryUTXRsRequest{
		TenantId: 1,
	})
	s.NoError(err)
	s.Equal(3, len(res.Utxrs))

	res, err = s.queryClient.UTXRs(s.ctx, &types.QueryUTXRsRequest{
		TenantId: 2,
	})
	s.NoError(err)
	s.Equal(3, len(res.Utxrs))
}

func (suite *SettlementTestSuite) TestKeeper_Tenant() {
	// Create Tenant
	tenant := &types.Tenant{
		Id:           3,
		Admins:       []string{s.appAdmin.String()},
		Denom:        "uusdc",
		PayoutMethod: "native",
		PayoutPeriod: 1,
	}
	suite.keeper.SetTenant(suite.ctx, tenant)

	// Get Tenant
	res, err := s.queryClient.Tenant(s.ctx, &types.QueryTenantRequest{TenantId: 3})
	s.NoError(err)
	s.Equal(tenant, res.Tenant.Tenant)
	s.Equal(types.GetTenantTreasuryAccount(3).String(), res.Tenant.Treasury.Address)
	s.Equal(sdk.NewCoin("uusdc", sdk.NewInt(0)), *res.Tenant.Treasury.Balance)

	// Deposit to Tenant Treasury
	newBalance := sdk.NewCoin("uusdc", sdk.NewInt(100))
	err = s.app.BankKeeper.SendCoins(s.ctx, s.creator, types.GetTenantTreasuryAccount(3), sdk.NewCoins(newBalance))
	s.NoError(err)

	// Check if the balance is updated
	res, err = s.queryClient.Tenant(s.ctx, &types.QueryTenantRequest{TenantId: 3})
	s.NoError(err)
	s.Equal(newBalance, *res.Tenant.Treasury.Balance)
}

func (suite *SettlementTestSuite) TestKeeper_Tenants() {
	// Create Tenant
	tenant := &types.Tenant{
		Id:           3,
		Admins:       []string{s.appAdmin.String()},
		PayoutPeriod: 1,
	}
	suite.keeper.SetTenant(suite.ctx, tenant)

	// Get Tenant
	res, err := s.queryClient.Tenants(s.ctx, &types.QueryTenantsRequest{})
	s.NoError(err)
	s.Equal(3, len(res.Tenants))
	s.Equal(tenant, res.Tenants[2].Tenant)
}
