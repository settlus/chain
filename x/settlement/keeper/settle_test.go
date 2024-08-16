package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v19/x/erc20/types"

	"github.com/evmos/evmos/v19/contracts"

	"github.com/settlus/chain/testutil/sample"
	utiltx "github.com/settlus/chain/testutil/tx"
	"github.com/settlus/chain/x/settlement/types"
)

func (suite *SettlementTestSuite) TestSettle_Settle_Native() {
	_, err := suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		TenantId: 1,
		Amount:   sdk.NewCoin("uusdc", math.NewInt(1000)),
		Sender:   suite.appAdmin.String(),
	})
	suite.NoError(err)

	requestId := "request-1"
	recipient := sdk.MustAccAddressFromBech32(sample.AccAddress())
	_, err = suite.keeper.CreateUTXR(suite.ctx, 1, &types.UTXR{
		RequestId:  requestId,
		Recipients: types.SingleRecipients(recipient),
		Amount:     sdk.NewCoin("uusdc", math.NewInt(10)),
		CreatedAt:  10,
	})
	suite.NoError(err)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 1, requestId)
	suite.NotNil(utxr)
	suite.EqualValues(10, utxr.Amount.Amount.Int64())

	initialBalance := s.app.BankKeeper.GetBalance(suite.ctx, recipient, "uusdc")
	suite.EqualValues(0, initialBalance.Amount.Int64())

	suite.keeper.Settle(suite.ctx.WithBlockHeight(11), 1)

	utxr = suite.keeper.GetUTXRByRequestId(suite.ctx.WithBlockHeight(12), 1, requestId)
	suite.Nil(utxr)

	balance := s.app.BankKeeper.GetBalance(suite.ctx, recipient, "uusdc")
	suite.EqualValues(10, balance.Amount.Int64())
}

func (suite *SettlementTestSuite) TestSettle_Settle_ERC20_Conversion() {
	tenantId := uint64(2)
	amount := math.NewInt(100)
	treasuryAddr := types.GetTenantTreasuryAccount(tenantId)

	suite.MintERC20Token(
		suite.erc20TokenPair.GetERC20Contract(),
		erc20types.ModuleAddress,
		common.Address(treasuryAddr),
		amount.BigInt(),
	)
	suite.Commit()

	treasuryBalance := suite.app.Erc20Keeper.BalanceOf(
		suite.ctx,
		contracts.ERC20MinterBurnerDecimalsContract.ABI,
		suite.erc20TokenPair.GetERC20Contract(),
		common.Address(treasuryAddr),
	)
	suite.Require().EqualValues(amount.Int64(), treasuryBalance.Int64())

	requestId := "req-1"
	recipient := sdk.MustAccAddressFromBech32(sample.AccAddress())
	_, err := suite.keeper.CreateUTXR(suite.ctx, tenantId, &types.UTXR{
		RequestId:  requestId,
		Recipients: types.SingleRecipients(recipient),
		Amount:     sdk.NewCoin(suite.erc20TokenPair.Denom, amount),
		CreatedAt:  10,
	})
	suite.NoError(err)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, tenantId, requestId)
	suite.NotNil(utxr)

	initialBalance := s.app.BankKeeper.GetBalance(suite.ctx, recipient, suite.erc20TokenPair.Denom)
	suite.EqualValues(0, initialBalance.Amount.Int64())

	suite.keeper.Settle(suite.ctx.WithBlockHeight(11), tenantId)

	utxr = suite.keeper.GetUTXRByRequestId(suite.ctx.WithBlockHeight(12), tenantId, requestId)
	suite.Nil(utxr)

	balance := s.app.BankKeeper.GetBalance(suite.ctx, recipient, suite.erc20TokenPair.Denom)
	suite.EqualValues(amount.Int64(), balance.Amount.Int64())
}

func (suite *SettlementTestSuite) TestSettle_Settle_InsufficientTreasuryBalance() {
	// give only 50
	_, err := suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		TenantId: 1,
		Amount:   sdk.NewCoin("uusdc", math.NewInt(50)),
		Sender:   suite.appAdmin.String(),
	})
	suite.NoError(err)
	s.Commit()

	recipient := sdk.AccAddress(utiltx.GenerateAddress().Bytes())

	// total amount to settle is 100
	for i := 0; i < 10; i++ {
		res, err := suite.keeper.CreateUTXR(suite.ctx, 1, &types.UTXR{
			RequestId:  fmt.Sprintf("request-%d", i),
			Recipients: types.SingleRecipients(recipient),
			Amount:     sdk.NewCoin("uusdc", math.NewInt(10)),
			CreatedAt:  uint64(suite.ctx.BlockHeight()),
		})

		suite.NoError(err)
		suite.NotNil(res)
	}

	suite.keeper.Settle(suite.ctx.WithBlockHeight(100), 1)

	// first utxr should be settled and deleted
	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 1, "request-1")
	suite.Nil(utxr)

	// this should not
	utxr = suite.keeper.GetUTXRByRequestId(suite.ctx, 1, "request-9")
	suite.NotNil(utxr)

	// treasury balance should be 0
	treasury := types.GetTenantTreasuryAccount(1)
	coins := suite.app.BankKeeper.SpendableCoins(suite.ctx, treasury)
	suite.EqualValues(0, len(coins))

	// recipient balance should be 50
	balance := s.app.BankKeeper.GetBalance(suite.ctx, recipient, "uusdc")
	suite.EqualValues(50, balance.Amount.Int64())
}

func (suite *SettlementTestSuite) TestSettle_Settle_TopUpTreasuryBalance() {
	_, err := suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		Amount:   sdk.NewCoin("uusdc", math.NewInt(50)),
	})
	suite.NoError(err)

	recipient := sdk.AccAddress(utiltx.GenerateAddress().Bytes())

	_, err = suite.keeper.CreateUTXR(suite.ctx, 1, &types.UTXR{
		RequestId:  "request-1",
		Recipients: types.SingleRecipients(recipient),
		Amount:     sdk.NewCoin("uusdc", math.NewInt(100)),
		CreatedAt:  uint64(suite.ctx.BlockHeight()),
	})
	suite.NoError(err)

	utxr := suite.keeper.GetUTXRByRequestId(suite.ctx, 1, "request-1")
	suite.NotNil(utxr)

	suite.keeper.Settle(suite.ctx.WithBlockHeight(100), 1)

	utxr = suite.keeper.GetUTXRByRequestId(suite.ctx, 1, "request-1")
	suite.NotNil(utxr)

	_, err = suite.msgServer.DepositToTreasury(suite.ctx, &types.MsgDepositToTreasury{
		Sender:   suite.appAdmin.String(),
		TenantId: 1,
		Amount:   sdk.NewCoin("uusdc", math.NewInt(50)),
	})
	suite.NoError(err)
	suite.Commit()

	suite.keeper.Settle(suite.ctx.WithBlockHeight(100), 1)

	utxr = suite.keeper.GetUTXRByRequestId(suite.ctx, 1, "request-1")
	suite.Nil(utxr)

	balance := s.app.BankKeeper.GetBalance(suite.ctx, recipient, "uusdc")
	suite.EqualValues(100, balance.Amount.Int64())
}
