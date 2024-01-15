package keeper_test

func (suite *SettlementTestSuite) TestKeeper_GetLargestTenantId() {
	id := suite.keeper.GetLargestTenantId(suite.ctx)
	suite.Equal(uint64(2), id)

	tenantId, err := suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(3), tenantId)

	id = suite.keeper.GetLargestTenantId(suite.ctx)
	suite.Equal(uint64(3), id)

	tenantId, err = suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(4), tenantId)

	id = suite.keeper.GetLargestTenantId(suite.ctx)
	suite.Equal(uint64(4), id)
}

func (suite *SettlementTestSuite) TestKeeper_GetTenant() {
	tenant := suite.keeper.GetTenant(suite.ctx, 0)
	suite.Nil(tenant)

	tenantId, err := suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(3), tenantId)

	tenant = suite.keeper.GetTenant(suite.ctx, tenantId)
	suite.NotNil(tenant)
	suite.Equal(uint64(3), tenant.Id)
	suite.Equal([]string{suite.appAdmin.String()}, tenant.Admins)
	suite.Equal(uint64(100), tenant.PayoutPeriod)

	tenantId, err = suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(4), tenantId)

	tenant = suite.keeper.GetTenant(suite.ctx, tenantId)
	suite.NotNil(tenant)
	suite.Equal(uint64(4), tenant.Id)
	suite.Equal([]string{suite.appAdmin.String()}, tenant.Admins)
	suite.Equal(uint64(100), tenant.PayoutPeriod)
}

func (suite *SettlementTestSuite) TestKeeper_CreateNewTenant() {
	tenantId, err := suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(3), tenantId)

	tenant := suite.keeper.GetTenant(suite.ctx, tenantId)
	suite.NotNil(tenant)
	suite.Equal(uint64(3), tenant.Id)
	suite.Equal([]string{suite.appAdmin.String()}, tenant.Admins)
	suite.Equal(uint64(100), tenant.PayoutPeriod)

	tenantId, err = suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(4), tenantId)

	tenant = suite.keeper.GetTenant(suite.ctx, tenantId)
	suite.NotNil(tenant)
	suite.Equal(uint64(4), tenant.Id)
	suite.Equal([]string{suite.appAdmin.String()}, tenant.Admins)
	suite.Equal(uint64(100), tenant.PayoutPeriod)
}

func (suite *SettlementTestSuite) TestKeeper_CheckAdminPermission() {
	tenantId, err := suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(3), tenantId)

	has := suite.keeper.CheckAdminPermission(suite.ctx, tenantId, suite.appAdmin.String())
	suite.True(has)

	has = suite.keeper.CheckAdminPermission(suite.ctx, tenantId, suite.creator.String())
	suite.False(has)
}

func (suite *SettlementTestSuite) TestKeeper_GetPayoutPeriod() {
	tenantId, err := suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(3), tenantId)

	payoutPeriod := suite.keeper.GetPayoutPeriod(suite.ctx, tenantId)
	suite.Equal(uint64(100), payoutPeriod)
}

func (suite *SettlementTestSuite) TestKeeper_CheckTenantExist() {
	tenantId, err := suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(3), tenantId)

	has := suite.keeper.CheckTenantExist(suite.ctx, tenantId)
	suite.True(has)

	has = suite.keeper.CheckTenantExist(suite.ctx, 100)
	suite.False(has)
}

func (suite *SettlementTestSuite) TestKeeper_SetTenant() {
	tenantId, err := suite.keeper.CreateNewTenant(suite.ctx, suite.appAdmin.String(), "uusdc", 100, "native", "")
	suite.NoError(err)
	suite.Equal(uint64(3), tenantId)

	tenant := suite.keeper.GetTenant(suite.ctx, tenantId)
	suite.NotNil(tenant)
	suite.Equal(uint64(3), tenant.Id)
	suite.Equal([]string{suite.appAdmin.String()}, tenant.Admins)
	suite.Equal(uint64(100), tenant.PayoutPeriod)

	tenant.Admins = append(tenant.Admins, suite.creator.String())
	tenant.PayoutPeriod = 200
	suite.keeper.SetTenant(suite.ctx, tenant)

	tenant = suite.keeper.GetTenant(suite.ctx, tenantId)
	suite.NotNil(tenant)
	suite.Equal(uint64(3), tenant.Id)
	suite.Equal([]string{suite.appAdmin.String(), suite.creator.String()}, tenant.Admins)
	suite.Equal(uint64(200), tenant.PayoutPeriod)
}
