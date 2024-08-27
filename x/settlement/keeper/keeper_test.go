package keeper_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	erc20types "github.com/evmos/evmos/v19/x/erc20/types"
	evmtypes "github.com/evmos/evmos/v19/x/evm/types"

	utiltx "github.com/settlus/chain/testutil/tx"
	"github.com/settlus/chain/utils"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	evmoscontracts "github.com/evmos/evmos/v19/contracts"
	"github.com/evmos/evmos/v19/crypto/ethsecp256k1"
	feemarkettypes "github.com/evmos/evmos/v19/x/feemarket/types"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/cmd/settlusd/config"
	"github.com/settlus/chain/contracts"
	"github.com/settlus/chain/testutil"
	"github.com/settlus/chain/x/settlement/keeper"
	"github.com/settlus/chain/x/settlement/types"
)

const (
	stakeDenom = config.BaseDenom
)

func newCoin(amt int64) sdk.Coin {
	return sdk.NewInt64Coin(stakeDenom, amt)
}

type SettlementTestSuite struct {
	suite.Suite

	app            *app.SettlusApp
	keeper         *keeper.SettlementKeeper
	msgServer      types.MsgServer
	ctx            sdk.Context
	validator      stakingtypes.Validator
	queryClient    types.QueryClient
	queryClientEvm evmtypes.QueryClient

	priv           cryptotypes.PrivKey
	appAdmin       sdk.AccAddress
	evmAddress     common.Address
	consAddress    sdk.ConsAddress
	erc20TokenPair *erc20types.TokenPair
	creator        sdk.AccAddress
	nftOwner       string
	sampleContract common.Address

	signer keyring.Signer
}

var (
	s *SettlementTestSuite
)

func (suite *SettlementTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func (suite *SettlementTestSuite) DoSetupTest(t require.TestingT) {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.appAdmin = priv.PubKey().Address().Bytes()
	suite.priv = priv
	suite.evmAddress = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = utiltx.NewSigner(priv)
	suite.creator = utiltx.GenerateAddress().Bytes()

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())
	suite.consAddress = consAddress

	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState(), utils.MainnetChainID)
	header := testutil.NewHeader(
		1, time.Now().UTC(), utils.MainnetChainID, consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	// query clients
	querier := keeper.Querier{SettlementKeeper: suite.app.SettlementKeeper}
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	suite.queryClient = types.NewQueryClient(queryHelper)
	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evmtypes.NewQueryClient(queryHelperEvm)

	// bond denoms
	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = config.BaseDenom
	if err = suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams); err != nil {
		panic(fmt.Errorf("failed to set staking params: %w", err))
	}

	evmParams := suite.app.EvmKeeper.GetParams(suite.ctx)
	evmParams.EvmDenom = config.BaseDenom
	err = suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)
	require.NoError(t, err)

	// keeper
	suite.keeper = suite.app.SettlementKeeper

	// msg server
	suite.msgServer = keeper.NewMsgServerImpl(suite.app.SettlementKeeper)

	// fund accounts
	err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, suite.creator, sdk.NewCoins(newCoin(1000000000000000000), testutil.NewMicroUSDC(100000000)))
	require.NoError(t, err)
	err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, suite.appAdmin, sdk.NewCoins(newCoin(1000000000000000000), testutil.NewMicroUSDC(100000000)))
	require.NoError(t, err)
	err = testutil.FundAccount(
		suite.ctx,
		suite.app.BankKeeper,
		suite.priv.PubKey().Address().Bytes(),
		sdk.NewCoins(newCoin(1000000000000000000), testutil.NewMicroUSDC(10000)),
	)
	require.NoError(t, err)

	// Set Validator
	valAddr := sdk.ValAddress(suite.appAdmin.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, privCons.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	err = suite.app.StakingKeeper.Hooks().AfterValidatorCreated(suite.ctx, validator.GetOperator())
	require.NoError(t, err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)

	// register USDC as an ERC20 token
	usdcContractAddr, err := suite.app.Erc20Keeper.DeployERC20Contract(suite.ctx, banktypes.Metadata{
		Name:        "uusdc",
		Symbol:      "uusdc",
		Description: "USDC",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "uusdc",
				Exponent: 6,
			},
		},
	})
	suite.Require().NoError(err)
	suite.Commit()

	tokenPair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, usdcContractAddr)
	suite.Require().NoError(err)
	suite.Commit()
	suite.erc20TokenPair = tokenPair

	// TODO change to setup with 1 validator
	validators := suite.app.StakingKeeper.GetValidators(s.ctx, 2)
	// set a bonded validator that takes part in consensus
	if validators[0].Status == stakingtypes.Bonded {
		suite.validator = validators[0]
	} else {
		suite.validator = validators[1]
	}

	suite.NoError(err)

	suite.addTestTenants()
}

func (suite *SettlementTestSuite) addTestTenants() {
	tenant1 := &types.Tenant{
		Id:           1,
		Admins:       []string{suite.appAdmin.String()},
		Denom:        "uusdc",
		PayoutPeriod: 1,
		PayoutMethod: "native",
	}
	tenant2 := &types.Tenant{
		Id:           2,
		Admins:       []string{suite.appAdmin.String()},
		Denom:        suite.erc20TokenPair.Denom,
		PayoutPeriod: 1,
		PayoutMethod: "native",
	}
	res, err := suite.msgServer.CreateTenant(suite.ctx, &types.MsgCreateTenant{Sender: suite.appAdmin.String(), Denom: tenant1.Denom, PayoutPeriod: tenant1.PayoutPeriod})
	suite.Require().NoError(err)
	suite.Require().Equal(tenant1.Id, res.TenantId)

	res, err = suite.msgServer.CreateTenant(suite.ctx, &types.MsgCreateTenant{Sender: suite.appAdmin.String(), Denom: tenant2.Denom, PayoutPeriod: tenant2.PayoutPeriod})
	suite.Require().NoError(err)
	suite.Require().Equal(tenant2.Id, res.TenantId)
}

func (suite *SettlementTestSuite) deploySampleContract() {
	ownerAddress := utiltx.GenerateAddress()
	contractAddress := suite.DeployAndMintSampleContract(ownerAddress)

	suite.nftOwner = ownerAddress.Hex()
	suite.sampleContract = contractAddress
}

// Commit commits and starts a new block with an updated context.
func (suite *SettlementTestSuite) Commit() {
	suite.CommitAndBeginBlockAfter(time.Hour * 1)
}

// Commit commits a block at a given time. Reminder: At the end of each
// Tendermint Consensus round the following methods are run
//  1. BeginBlock
//  2. DeliverTx
//  3. EndBlock
//  4. Commit
func (suite *SettlementTestSuite) CommitAndBeginBlockAfter(t time.Duration) {
	var err error
	suite.ctx, err = testutil.Commit(suite.ctx, suite.app, t, nil)
	suite.Require().NoError(err)

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evmtypes.NewQueryClient(queryHelperEvm)
}

// DeployAndMintSampleContract deploys the ERC721 contract with the provided name and symbol
func (suite *SettlementTestSuite) DeployAndMintSampleContract(ownerAddress common.Address) common.Address {
	// TODO: remove excessive Commit() calls
	suite.Commit()
	addr, err := suite.DeployNFTContract("Bored Ape Yacht Club", "BAYC")
	suite.NoError(err)
	suite.Commit()

	suite.Commit()
	err = suite.MintNFT(addr, ownerAddress)
	suite.NoError(err)
	suite.Commit()

	suite.Commit()
	err = suite.MintNFT(addr, ownerAddress)
	suite.NoError(err)
	suite.Commit()

	exists, err := suite.CheckNFTExists(addr, big.NewInt(0))
	suite.NoError(err)
	suite.True(exists)

	return addr
}

// DeployNFTContract deploys the ERC721 contract with the provided name and symbol
func (suite *SettlementTestSuite) DeployNFTContract(name, symbol string) (common.Address, error) {
	suite.Commit()
	addr, err := testutil.DeployContract(
		suite.ctx,
		suite.app,
		suite.priv,
		suite.queryClientEvm,
		contracts.ERC721Contract,
		name, symbol,
	)
	suite.Commit()
	return addr, err
}

func (suite *SettlementTestSuite) MintNFT(contractAddress, ownerAddress common.Address) error {
	suite.Commit()
	err := testutil.MintNFT(
		suite.ctx,
		suite.app,
		suite.priv,
		contracts.ERC721Contract,
		contractAddress,
		ownerAddress,
	)
	suite.Commit()
	return err
}

func (suite *SettlementTestSuite) CheckNFTExists(contractAddress common.Address, tokenId *big.Int) (bool, error) {
	exists, err := testutil.CheckNFTExists(
		suite.ctx,
		suite.app,
		suite.priv,
		contracts.ERC721Contract,
		contractAddress,
		tokenId,
	)

	return exists, err
}

func (suite *SettlementTestSuite) MintERC20Token(contractAddr, from, to common.Address, amount *big.Int) {
	transferData, err := evmoscontracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)
	_, err = suite.app.EvmKeeper.CallEVMWithData(suite.ctx, from, &contractAddr, transferData, true)
	suite.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	s = new(SettlementTestSuite)
	suite.Run(t, s)
}
