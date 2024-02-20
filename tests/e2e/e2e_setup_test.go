package e2e

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up e2e integration test suite...")
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e integration test suite...")
}

func (s *IntegrationTestSuite) TestBasic() {
	s.Run("test basic", func() {
		s.T().Log("testing basic functionality...")
	})
}
