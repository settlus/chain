package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/x/settlement/types"
)

func (suite *SettlementTestSuite) TestParams_GetParams() {
	params := types.DefaultParams()
	s.keeper.SetParams(s.ctx, params)
	params = suite.keeper.GetParams(suite.ctx)
	suite.EqualValues(types.DefaultParams(), params)
}

func (suite *SettlementTestSuite) TestParams_SetParams() {
	params := suite.keeper.GetParams(suite.ctx)
	params.OracleFeePercentage = sdk.NewDecWithPrec(1, 2)
	suite.keeper.SetParams(suite.ctx, params)

	newParams := suite.keeper.GetParams(suite.ctx)
	suite.EqualValues(params, newParams)
}
