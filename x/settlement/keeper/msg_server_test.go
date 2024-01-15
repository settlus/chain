package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/settlus/chain/testutil/sample"
	utiltx "github.com/settlus/chain/testutil/tx"
	settlustypes "github.com/settlus/chain/types"

	"github.com/settlus/chain/x/settlement/types"
)

func (suite *SettlementTestSuite) TestMsgServer_Record() {
	suite.deploySampleContract()
	res, err := suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.appAdmin.String(),
		TenantId:        1,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: suite.sampleContract.String(),
		TokenIdHex:      "0x0",
	})

	suite.NoError(err)
	suite.NotNil(res)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 1, "request-1")
	suite.NotNil(utxr)
	suite.EqualValues(utxr.Amount.Amount.Int64(), 10)
	suite.EqualValues(utxr.Recipient, suite.nftOwner)
}

func (suite *SettlementTestSuite) TestMsgServer_Record_CreateAccount() {
	suite.deploySampleContract()
	newAccountEVM := sample.EthAddress()
	newAccount := newAccountEVM.Bytes()
	err := suite.MintNFT(suite.sampleContract, newAccountEVM)
	suite.NoError(err)

	acc := suite.app.AccountKeeper.GetAccount(suite.ctx, newAccount)
	suite.Nil(acc)

	res, err := suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.appAdmin.String(),
		TenantId:        1,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: suite.sampleContract.String(),
		TokenIdHex:      "0x1",
	})

	suite.NoError(err)
	suite.NotNil(res)

	acc = suite.app.AccountKeeper.GetAccount(suite.ctx, newAccount)
	suite.NotNil(acc)
}

func (suite *SettlementTestSuite) TestMsgServer_Record_NonAdminMustFail() {
	suite.deploySampleContract()
	res, err := suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.creator.String(),
		TenantId:        1,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: suite.sampleContract.String(),
		TokenIdHex:      "0x0",
	})

	suite.Error(err)
	suite.Nil(res)
}

func (suite *SettlementTestSuite) TestMsgServer_Record_InvalidTenantID() {
	suite.deploySampleContract()
	res, err := suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.appAdmin.String(),
		TenantId:        100,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: suite.sampleContract.String(),
		TokenIdHex:      "0x0",
	})

	suite.Error(err)
	suite.Nil(res)
}

func (suite *SettlementTestSuite) TestMsgServer_Record_UniqueRequestID() {
	suite.deploySampleContract()
	res, err := suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.appAdmin.String(),
		TenantId:        1,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: suite.sampleContract.String(),
		TokenIdHex:      "0x0",
	})

	suite.NoError(err)
	suite.NotNil(res)

	res, err = suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.appAdmin.String(),
		TenantId:        1,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: suite.sampleContract.String(),
		TokenIdHex:      "0x0",
	})

	suite.Error(err)
	suite.Nil(res)

	res, err = suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.appAdmin.String(),
		TenantId:        2,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: suite.sampleContract.String(),
		TokenIdHex:      "0x0",
	})

	suite.NoError(err)
	suite.NotNil(res)

	res, err = suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.appAdmin.String(),
		TenantId:        2,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: suite.sampleContract.String(),
		TokenIdHex:      "0x0",
	})

	suite.Error(err)
	suite.Nil(res)
}

func (suite *SettlementTestSuite) TestMsgServer_Record_IDInOrder() {
	suite.deploySampleContract()
	for i := 1; i <= 5; i++ {
		res, err := suite.msgServer.Record(suite.ctx, &types.MsgRecord{
			Sender:          suite.appAdmin.String(),
			TenantId:        1,
			RequestId:       fmt.Sprintf("request-%d", i),
			Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
			ChainId:         "settlus_5371-1",
			ContractAddress: suite.sampleContract.String(),
			TokenIdHex:      "0x0",
		})

		suite.NoError(err)
		suite.NotNil(res)
	}

	utxrs := suite.keeper.GetAllUTXRWithTenantAndID(suite.ctx)
	for i, utxr := range utxrs {
		suite.EqualValues(utxr.Id, i+1)
	}
}

func (suite *SettlementTestSuite) TestMsgServer_Record_NftNotMinted() {
	ownerAddress := utiltx.GenerateAddress()
	contractAddress := s.DeployAndMintSampleContract(ownerAddress)
	res, err := suite.msgServer.Record(suite.ctx, &types.MsgRecord{
		Sender:          suite.appAdmin.String(),
		TenantId:        1,
		RequestId:       "request-1",
		Amount:          sdk.NewCoin("uusdc", math.NewInt(10)),
		ChainId:         "settlus_5371-1",
		ContractAddress: contractAddress.String(),
		TokenIdHex:      "0xa",
	})

	suite.Error(err)
	suite.True(types.ErrEVMCallFailed.Is(err))
	suite.Nil(res)
}

func (suite *SettlementTestSuite) TestMsgServer_Record_InvalidChainId() {
	// TODO
}

func (suite *SettlementTestSuite) TestMsgServer_Record_InvalidContractAddress() {
	// TODO
}

func (suite *SettlementTestSuite) TestMsgServer_Record_InvalidTokenId() {
	// TODO
}

func (suite *SettlementTestSuite) TestMsgServer_AddTenantAdmin() {
	newAdmin := sample.AccAddress()
	res, err := suite.msgServer.AddTenantAdmin(suite.ctx, &types.MsgAddTenantAdmin{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		NewAdmin: newAdmin,
	})

	suite.NoError(err)
	suite.NotNil(res)

	tenant := suite.keeper.GetTenant(suite.ctx, 1)
	for _, admin := range tenant.Admins {
		if admin == newAdmin {
			return
		}
	}
	suite.Error(fmt.Errorf("new admin %s not found after AddTenantAdmin", newAdmin))
}

func (suite *SettlementTestSuite) TestMsgServer_AddTenantAdmin_DuplicateAdmin() {
	newAdmin := sample.AccAddress()
	res, err := suite.msgServer.AddTenantAdmin(suite.ctx, &types.MsgAddTenantAdmin{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		NewAdmin: newAdmin,
	})

	suite.NoError(err)
	suite.NotNil(res)

	res, err = suite.msgServer.AddTenantAdmin(suite.ctx, &types.MsgAddTenantAdmin{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		NewAdmin: newAdmin,
	})

	suite.Error(err)
	suite.True(types.ErrInvalidAdmin.Is(err))
	suite.Nil(res)
}

func (suite *SettlementTestSuite) TestMsgServer_RemoveTenantAdmin() {
	newAdmin := sample.AccAddress()
	res, err := suite.msgServer.AddTenantAdmin(suite.ctx, &types.MsgAddTenantAdmin{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		NewAdmin: newAdmin,
	})

	suite.NoError(err)
	suite.NotNil(res)

	removeRes, err := suite.msgServer.RemoveTenantAdmin(suite.ctx, &types.MsgRemoveTenantAdmin{
		Sender:        suite.appAdmin.String(),
		TenantId:      1,
		AdminToRemove: newAdmin,
	})

	suite.NoError(err)
	suite.NotNil(removeRes)

	tenant := suite.keeper.GetTenant(suite.ctx, 1)
	for _, admin := range tenant.Admins {
		if admin == newAdmin {
			suite.Error(fmt.Errorf("new admin %s found after RemoveTenantAdmin", newAdmin))
		}
	}
}

func (suite *SettlementTestSuite) TestMsgServer_RemoveTenantAdmin_LastAdmin() {
	res, err := suite.msgServer.RemoveTenantAdmin(suite.ctx, &types.MsgRemoveTenantAdmin{
		Sender:        suite.appAdmin.String(),
		TenantId:      1,
		AdminToRemove: suite.appAdmin.String(),
	})

	suite.Error(err)
	suite.True(types.ErrCannotRemoveAdmin.Is(err))
	suite.Nil(res)
}

func (suite *SettlementTestSuite) TestMsgServer_UpdateTenantPayoutPeriod() {
	res, err := suite.msgServer.UpdateTenantPayoutPeriod(suite.ctx, &types.MsgUpdateTenantPayoutPeriod{
		Sender:       suite.appAdmin.String(),
		TenantId:     1,
		PayoutPeriod: 1000,
	})

	suite.NoError(err)
	suite.NotNil(res)

	tenant := suite.keeper.GetTenant(suite.ctx, 1)
	suite.EqualValues(tenant.PayoutPeriod, 1000)
}

func (suite *SettlementTestSuite) TestMsgServer_UpdateTenantPayoutPeriod_InvalidTenantId() {
	res, err := suite.msgServer.UpdateTenantPayoutPeriod(suite.ctx, &types.MsgUpdateTenantPayoutPeriod{
		Sender:       suite.appAdmin.String(),
		TenantId:     100,
		PayoutPeriod: 1000,
	})

	suite.Error(err)
	suite.Nil(res)
}

func (suite *SettlementTestSuite) TestMsgServer_DepositToTreasury() {
	_, err := suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		Amount:   sdk.NewCoin("uusdc", math.NewInt(10)),
	})

	suite.NoError(err)

	tenantAccount := types.GetTenantTreasuryAccount(1)
	coins := suite.app.BankKeeper.SpendableCoins(suite.ctx, tenantAccount)
	totalBalance := coins[0].Amount.Int64()

	suite.EqualValues(10, totalBalance)
	suite.EqualValues("uusdc", coins[0].Denom)
}

func (suite *SettlementTestSuite) TestMsgServer_DepositToTreasury_InvalidTenantId() {
	_, err := suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		Sender:   suite.appAdmin.String(),
		TenantId: 100,
		Amount:   sdk.NewCoin("uusdc", math.NewInt(10)),
	})

	suite.Error(err)

	coins := suite.app.BankKeeper.SpendableCoins(suite.ctx, suite.app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Empty(coins)
}

func (suite *SettlementTestSuite) TestMsgServer_DepositToTreasury_InsufficientBalance() {
	_, err := suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		Amount:   sdk.NewCoin("ueurc", math.NewInt(10000000000)),
	})

	suite.Error(err)
}

func (suite *SettlementTestSuite) TestMsgServer_DepositToTreasury_MultipleTenants() {
	_, err := suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		Amount:   sdk.NewCoin("uusdc", math.NewInt(10)),
	})
	suite.NoError(err)

	_, err = suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		Sender:   suite.appAdmin.String(),
		TenantId: 2,
		Amount:   sdk.NewCoin("uusdc", math.NewInt(20)),
	})
	suite.NoError(err)

	tenant0Account := types.GetTenantTreasuryAccount(1)
	suite.NotNil(tenant0Account)
	tenant1Account := types.GetTenantTreasuryAccount(2)
	suite.NotNil(tenant1Account)

	suite.NotEqual(tenant0Account, tenant1Account)

	coins0 := suite.app.BankKeeper.SpendableCoins(suite.ctx, tenant0Account)
	coins1 := suite.app.BankKeeper.SpendableCoins(suite.ctx, tenant1Account)

	suite.EqualValues(coins0[0].Amount.Int64(), 10)
	suite.EqualValues(coins1[0].Amount.Int64(), 20)
}

func (suite *SettlementTestSuite) TestMsgServer_Cancel() {
	requestId := "request-cancel"
	_, err := suite.keeper.CreateUTXR(suite.ctx, 1, &types.UTXR{
		RequestId:   requestId,
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", math.NewInt(10)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 1, requestId)
	suite.NotNil(utxr)

	_, err = suite.msgServer.Cancel(suite.ctx, &types.MsgCancel{
		Sender:    suite.appAdmin.String(),
		TenantId:  1,
		RequestId: requestId,
	})
	suite.NoError(err)

	utxr = suite.keeper.GetUTXRByRequestId(suite.ctx, 1, requestId)
	suite.Nil(utxr)
}

func (suite *SettlementTestSuite) TestMsgServer_Cancel_InvalidTenant() {
	requestId := "request-cancel"
	_, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   requestId,
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", math.NewInt(10)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)

	_, err = suite.msgServer.Cancel(suite.ctx, &types.MsgCancel{
		Sender:    suite.appAdmin.String(),
		TenantId:  100,
		RequestId: requestId,
	})

	suite.Error(err)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 0, requestId)
	suite.NotNil(utxr)

	utxr = suite.keeper.GetUTXRByRequestId(suite.ctx, 1, requestId)
	suite.Nil(utxr)
}

func (suite *SettlementTestSuite) TestMsgServer_Cancel_NonAdmin() {
	requestId := "request-cancel"
	_, err := suite.keeper.CreateUTXR(suite.ctx, 0, &types.UTXR{
		RequestId:   requestId,
		Recipient:   settlustypes.NewHexAddressString(suite.appAdmin),
		Amount:      sdk.NewCoin("uusdc", math.NewInt(10)),
		PayoutBlock: uint64(100),
	})
	suite.NoError(err)

	_, err = suite.msgServer.Cancel(suite.ctx, &types.MsgCancel{
		Sender:    suite.creator.String(),
		TenantId:  3,
		RequestId: requestId,
	})
	suite.Error(err)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 0, requestId)
	suite.NotNil(utxr)
}
