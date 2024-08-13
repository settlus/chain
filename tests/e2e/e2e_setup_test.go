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
	uusdcDenom     = "uusdc"
)

var (
	standardFees = sdk.NewCoin(asetlDenom, sdk.NewInt(300000000000000))
	faucetAddr   = "settlus10z74aw2m660tuezej4w5zr35zye6684t5ejjmk" // faucet address specified in the config.yml
	bobAddr      = "settlus1vfhltz7wr4ca862xd0azjuap4tupwgyzk7qukp" // bob address specified in the config.yml
)

type IntegrationTestSuite struct {
	suite.Suite
	ethClient       *ethclient.Client
	internalNftAddr string
	admin           string
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up e2e integration test suite...")

	s.SetupSettlementTestSuite()
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
