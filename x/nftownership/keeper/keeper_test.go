package keeper_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/settlus/chain/evmos/crypto/ethsecp256k1"
	evmtypes "github.com/settlus/chain/evmos/x/evm/types"
	feemarkettypes "github.com/settlus/chain/evmos/x/feemarket/types"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/cmd/settlusd/config"
	"github.com/settlus/chain/contracts"
	"github.com/settlus/chain/testutil"
	utiltx "github.com/settlus/chain/testutil/tx"
	"github.com/settlus/chain/x/nftownership/types"
)

type NftOwnershipTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	app            *app.App
	queryClient    types.QueryClient
	queryClientEvm evmtypes.QueryClient

	creator     common.Address
	address     common.Address
	consAddress sdk.ConsAddress
	clientCtx   client.Context //nolint:unused
	ethSigner   ethtypes.Signer
	priv        cryptotypes.PrivKey
	validator   stakingtypes.Validator

	signer keyring.Signer
}

var (
	s *NftOwnershipTestSuite
)

func TestKeeperTestSuite(t *testing.T) {
	s = new(NftOwnershipTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *NftOwnershipTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func (suite *NftOwnershipTestSuite) DoSetupTest(t require.TestingT) {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.priv = priv
	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = utiltx.NewSigner(priv)

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())
	suite.consAddress = consAddress

	// init app
	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState())
	header := testutil.NewHeader(
		1, time.Now().UTC(), "settlus_5371-1", consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	// query clients
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.NftOwnershipKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evmtypes.NewQueryClient(queryHelperEvm)

	// bond denom
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

	creator := utiltx.GenerateAddress()
	suite.creator = creator
	creatorCosmosAddress := sdk.AccAddress(creator.Bytes())
	err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, creatorCosmosAddress, sdk.NewCoins(testutil.NewSetl(1000000000000000000), testutil.NewMicroUSDC(10000)))
	require.NoError(t, err)

	addressCosmosAddress := sdk.AccAddress(suite.address.Bytes())
	err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addressCosmosAddress, sdk.NewCoins(testutil.NewSetl(1000000000000000000), testutil.NewMicroUSDC(10000)))
	require.NoError(t, err)

	// Set Validator
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, privCons.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	err = suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
	require.NoError(t, err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)

	// fund signer acc to pay for tx fees
	amt := sdk.NewInt(int64(math.Pow10(18) * 2))
	err = testutil.FundAccount(
		suite.ctx,
		suite.app.BankKeeper,
		suite.priv.PubKey().Address().Bytes(),
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amt)),
	)
	suite.Require().NoError(err)

	// fund nft license module account
	err = testutil.FundModuleAccount(
		suite.ctx,
		suite.app.BankKeeper,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amt)),
	)
	suite.Require().NoError(err)

	// TODO change to setup with 1 validator
	validators := s.app.StakingKeeper.GetValidators(s.ctx, 2)
	// set a bonded validator that takes part in consensus
	if validators[0].Status == stakingtypes.Bonded {
		suite.validator = validators[0]
	} else {
		suite.validator = validators[1]
	}

	suite.ethSigner = ethtypes.LatestSignerForChainID(s.app.EvmKeeper.ChainID())
}

// Commit commits and starts a new block with an updated context.
func (suite *NftOwnershipTestSuite) Commit() {
	suite.CommitAndBeginBlockAfter(time.Hour * 1)
}

// Commit commits a block at a given time. Reminder: At the end of each
// Tendermint Consensus round the following methods are run
//  1. BeginBlock
//  2. DeliverTx
//  3. EndBlock
//  4. Commit
func (suite *NftOwnershipTestSuite) CommitAndBeginBlockAfter(t time.Duration) {
	var err error
	suite.ctx, err = testutil.Commit(suite.ctx, suite.app, t, nil)
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	suite.queryClientEvm = evmtypes.NewQueryClient(queryHelper)
}

// DeployContract deploys the ERC721 contract with the provided name and symbol
func (suite *NftOwnershipTestSuite) DeployContract(name, symbol string) (common.Address, error) {
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

func (suite *NftOwnershipTestSuite) MintNFT(contractAddress, ownerAddress common.Address) error {
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

func (suite *NftOwnershipTestSuite) CheckNFTExists(contractAddress common.Address, tokenId *big.Int) (bool, error) {
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
