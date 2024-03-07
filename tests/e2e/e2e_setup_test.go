package e2e

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/suite"
)

const (
	settlusdBinary = "settlusd"
	txCommand      = "tx"
	queryCommand   = "query"
	keysCommand    = "keys"
	setlDenom      = "setl"
	asetlDenom     = "asetl"
)

var (
	standardFees = sdk.NewCoin(asetlDenom, sdk.NewInt(330000))
)

type IntegrationTestSuite struct {
	suite.Suite
	ethClient       *ethclient.Client
	internalNftAddr string
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up e2e integration test suite...")

	s.T().Log("connecting to Ethereum JSON-RPC...")
	ethClient, err := ethclient.Dial(ethAPIEndpoint)
	s.Require().NoError(err)
	s.ethClient = ethClient

	s.T().Log("mint NFT for settlement Test")
	contractAddr, err := mintNFTContract(ethClient)
	s.internalNftAddr = contractAddr

	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e integration test suite...")
}

func (s *IntegrationTestSuite) TestBasic() {
	s.Run("test basic", func() {
		s.T().Log("testing basic functionality...")
		s.Run("get latest block", func() {
			blockId, err := queryLatestBlockId(chainAPIEndpoint)

			s.Require().NoError(err)
			s.Require().NotEmpty(blockId)
		})
	})
}
