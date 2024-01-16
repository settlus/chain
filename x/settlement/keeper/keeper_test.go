package keeper_test

import (
	"math/big"
	"testing"
	"time"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	evmtypes "github.com/settlus/chain/evmos/x/evm/types"
	utiltx "github.com/settlus/chain/testutil/tx"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/settlus/chain/evmos/crypto/ethsecp256k1"
	erc20types "github.com/settlus/chain/evmos/x/erc20/types"
	feemarkettypes "github.com/settlus/chain/evmos/x/feemarket/types"

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

	app            *app.App
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
	erc20Address   common.Address
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
	chainId := "settlus_5371-1"

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

	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState())
	header := testutil.NewHeader(
		1, time.Now().UTC(), chainId, consAddress, nil, nil,
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
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	evmParams := suite.app.EvmKeeper.GetParams(suite.ctx)
	evmParams.EvmDenom = config.BaseDenom
	err = suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)
	require.NoError(t, err)

	// nft ownership module params
	nftownershipParams := suite.app.NftOwnershipKeeper.GetParams(suite.ctx)
	nftownershipParams.AllowedChainIds = []string{"settlus_5371-1"}
	suite.app.NftOwnershipKeeper.SetParams(suite.ctx, nftownershipParams)

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
	err = suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
	require.NoError(t, err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)

	tokenPair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, banktypes.Metadata{
		Description: "USDC",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "uusdc",
				Exponent: 0,
			},
		},
		Base:    "uusdc",
		Display: "uusdc",
		Name:    "uusdc",
		Symbol:  "uusdc",
	})

	// TODO change to setup with 1 validator
	validators := suite.app.StakingKeeper.GetValidators(s.ctx, 2)
	// set a bonded validator that takes part in consensus
	if validators[0].Status == stakingtypes.Bonded {
		suite.validator = validators[0]
	} else {
		suite.validator = validators[1]
	}

	suite.erc20Address = common.HexToAddress(tokenPair.Erc20Address)
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
		Denom:        "uusdc",
		PayoutPeriod: 1,
		PayoutMethod: "native",
	}
	suite.keeper.SetTenant(suite.ctx, tenant1)
	suite.keeper.SetTenant(suite.ctx, tenant2)
}

func (suite *SettlementTestSuite) deploySampleContract() {
	ownerAddress := utiltx.GenerateAddress()
	contractAddress := s.DeployAndMintSampleContract(ownerAddress)

	suite.nftOwner = ownerAddress.Hex()
	suite.sampleContract = contractAddress
}

func (suite *SettlementTestSuite) balanceOf(contract, account common.Address) any { //nolint
	erc20 := contracts.ERC20Contract.ABI

	// TODO: use erc20 BalanceOf keeper
	res, err := suite.app.EvmKeeper.CallEVM(suite.ctx, erc20, erc20types.ModuleAddress, contract, false, "balanceOf", account)
	if err != nil {
		return nil
	}

	unpacked, err := erc20.Unpack("balanceOf", res.Ret)
	if err != nil {
		return nil
	}
	if len(unpacked) == 0 {
		return nil
	}

	return unpacked[0]
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
	addr, err := s.DeployContract("Bored Ape Yatch Club", "BAYC")
	suite.NoError(err)

	err = s.MintNFT(addr, ownerAddress)
	suite.NoError(err)

	return addr
}

// DeployContract deploys the ERC721 contract with the provided name and symbol
func (suite *SettlementTestSuite) DeployContract(name, symbol string) (common.Address, error) {
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

func TestKeeperTestSuite(t *testing.T) {
	s = new(SettlementTestSuite)
	suite.Run(t, s)
}
