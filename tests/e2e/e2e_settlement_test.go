package e2e

import (
	"math/rand"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/settlus/chain/x/settlement/types"
)

func (s *IntegrationTestSuite) TestBasicTenant() {
	admin := "settlus1vfhltz7wr4ca862xd0azjuap4tupwgyzk7qukp"
	denom := "setl"
	var tenantId uint64 = 0

	s.Run("create_tenant", func() {
		period := 10
		beforeTenants, err := queryTenants(chainAPIEndpoint)
		s.Require().NoError(err)

		s.execCreateTenant(admin, denom, strconv.FormatInt(int64(period), 10), "1010000uusdc")
		s.Require().Eventually(
			func() bool {
				afterTenants, err := queryTenants(chainAPIEndpoint)
				s.Require().NoError(err)

				newTenant := afterTenants[len(afterTenants)-1]
				s.Require().Equal(denom, newTenant.Tenant.Denom)
				s.Require().Equal(period, int(newTenant.Tenant.PayoutPeriod))
				s.Require().Equal(types.PayoutMethod_Native, newTenant.Tenant.PayoutMethod)

				tenantId = newTenant.Tenant.Id

				return len(afterTenants) == len(beforeTenants)+1
			},
			time.Minute,
			2*time.Second,
		)
	})

	s.Run("record_external_nft_revenue", func() {
		requestId := "req-" + strconv.Itoa(rand.Intn(10000000000))
		revenue := sdk.NewCoin("setl", sdk.NewInt(10))
		s.execRecord(admin, tenantId, requestId, revenue.String(), "1", "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "0x0", "10000uusdc")
		s.Require().Eventually(
			func() bool {
				utxr, err := queryUtxr(chainAPIEndpoint, "1", requestId)

				s.Require().NoError(err)
				s.Require().Equal("0xf7801B8115f3Fe46AC55f8c0Fdb5243726bdb66A", utxr.Recipient)
				s.Require().Equal(revenue, utxr.Amount)

				return err == nil
			},
			time.Minute,
			2*time.Second,
		)
	})
}
