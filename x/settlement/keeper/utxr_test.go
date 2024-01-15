package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	settlustypes "github.com/settlus/chain/types"
	"github.com/settlus/chain/x/settlement/types"
)

func (suite *SettlementTestSuite) TestKeeper_HasUTXRByRequestId() {
	utxrId, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-1",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(1), utxrId)

	has := suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-1")
	suite.True(has)

	has = suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-2")
	suite.False(has)
}

func (suite *SettlementTestSuite) TestKeeper_GetLatestUTXR() {
	id := suite.keeper.GetLargestUTXRId(suite.ctx, 0)
	suite.Equal(uint64(0), id)

	utxrId, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-1",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(1), utxrId)

	has := suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-1")
	suite.True(has)

	id = suite.keeper.GetLargestUTXRId(suite.ctx, 0)
	suite.Equal(uint64(1), id)

	utxrId, err = suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-2",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(2), utxrId)

	has = suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-2")
	suite.True(has)

	id = suite.keeper.GetLargestUTXRId(suite.ctx, 0)
	suite.Equal(uint64(2), id)

	deletedUtxrId, err := suite.keeper.DeleteUTXRByRequestId(suite.ctx, 0, "request-2")
	suite.NoError(err)
	suite.Equal(uint64(2), deletedUtxrId)

	has = suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-2")
	suite.False(has)

	id = suite.keeper.GetLargestUTXRId(suite.ctx, 0)
	suite.Equal(uint64(1), id)
}

func (suite *SettlementTestSuite) TestKeeper_GetLatestUTXR_MultipleTenants() {
	t0Id := suite.keeper.GetLargestUTXRId(suite.ctx, 0)
	suite.Equal(uint64(0), t0Id)

	t1Id := suite.keeper.GetLargestUTXRId(suite.ctx, 1)
	suite.Equal(uint64(0), t1Id)

	utxrId, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-1",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(1), utxrId)

	has := suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-1")
	suite.True(has)

	t0Id = suite.keeper.GetLargestUTXRId(suite.ctx, 0)
	suite.Equal(uint64(1), t0Id)

	utxrId, err = suite.keeper.CreateUTXR(suite.ctx, 1, &types.UTXR{
		RequestId:   "request-1",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(1), utxrId)

	has = suite.keeper.HasUTXRByRequestId(suite.ctx, 1, "request-1")
	suite.True(has)

	t1Id = suite.keeper.GetLargestUTXRId(suite.ctx, 1)
	suite.Equal(uint64(1), t1Id)

	utxrId, err = suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-2",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(2), utxrId)

	has = suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-2")
	suite.True(has)

	t0Id = suite.keeper.GetLargestUTXRId(suite.ctx, 0)
	suite.Equal(uint64(2), t0Id)

	t1Id = suite.keeper.GetLargestUTXRId(suite.ctx, 1)
	suite.Equal(uint64(1), t1Id)
}

func (suite *SettlementTestSuite) TestKeeper_CreateUTXR() {
	utxrId, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-1",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(1), utxrId)

	has := suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-1")
	suite.True(has)
}

func (suite *SettlementTestSuite) TestKeeper_GetUTXRByRequestId() {
	_, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-1",
		Recipient:   settlustypes.NewHexAddressString(suite.creator),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(10)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 0, "request-1")
	suite.NotNil(utxr)
	suite.EqualValues(10, utxr.Amount.Amount.Int64())
	suite.EqualValues(suite.creator, common.FromHex(utxr.Recipient.String()))
}

func (suite *SettlementTestSuite) TestKeeper_DeleteUTXRByRequestId() {
	_, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-1",
		Recipient:   settlustypes.NewHexAddressString(suite.creator),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(10)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)

	suite.NoError(err)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 0, "request-1")
	suite.NotNil(utxr)
	suite.EqualValues(10, utxr.Amount.Amount.Int64())
	suite.EqualValues(suite.creator, common.FromHex(utxr.Recipient.String()))

	utxrId, err := suite.keeper.DeleteUTXRByRequestId(suite.ctx, 0, "request-1")
	suite.NoError(err)
	suite.Equal(uint64(1), utxrId)

	has := suite.keeper.HasUTXRByRequestId(suite.ctx, 0, "request-1")
	suite.False(has)
}

func (suite *SettlementTestSuite) TestKeeper_GetAllUTXRWithTenantAndID() {
	utxrId, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-1",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(1), utxrId)

	utxrId, err = suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   "request-2",
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", sdk.NewInt(100)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)
	suite.Equal(uint64(2), utxrId)

	utxrs := suite.keeper.GetAllUTXRWithTenantAndID(suite.ctx)
	suite.Equal(2, len(utxrs))
	suite.Equal("request-1", utxrs[0].Utxr.RequestId)
	suite.Equal(uint64(1), utxrs[0].Id)
	suite.Equal("request-2", utxrs[1].Utxr.RequestId)
	suite.Equal(uint64(2), utxrs[1].Id)
}
