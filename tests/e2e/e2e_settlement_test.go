package e2e

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/settlus/chain/x/settlement/types"
)

const (
	extChainId    = "1"
	extNftAddress = "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d"
	extNftId      = "0x0"
	extNftOwner   = "0xf7801b8115f3fe46ac55f8c0fdb5243726bdb66a"

	internalNftId    = "0x0"
	internalNftOwner = "0xa7801b8115f3fe46ac55f8c0fdb5243726bdb66a"
)

func (s *IntegrationTestSuite) SetupSettlementTestSuite() {
	s.T().Log("adding new key...")
	admin := "admin" + strconv.Itoa(rand.Intn(10000000000))
	s.admin = s.execKeyAdd(admin)

	s.T().Log("creating new account and transferring initial fund...")
	initialAmount := fmt.Sprintf("%d%s,%d%s", 10000000000000, asetlDenom, 10000000000000, uusdcDenom)
	s.execBankSend(treasuryAddr, s.admin, initialAmount, standardFees.String())

	s.T().Log("connecting to Ethereum JSON-RPC...")
	ethClient, err := ethclient.Dial(ethAPIEndpoint)
	s.Require().NoError(err)
	s.ethClient = ethClient

	s.T().Log("deploying NFT contract...")
	contractAddr, err := deployNFTContract(ethClient)
	s.Require().NoError(err)
	s.internalNftAddr = contractAddr
	s.T().Log("NFT:", contractAddr)

	s.T().Log("minting NFT...")
	err = mintNFT(ethClient, contractAddr, internalNftOwner)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestNativeTenant() {
	denom := asetlDenom
	revenue := sdk.NewCoin(denom, sdk.NewInt(10))
	var tenantId uint64 = 0

	pass := s.Run("create_new_tenant", func() {
		period := 8
		beforeTenants, err := queryTenants(chainAPIEndpoint)
		s.Require().NoError(err)

		s.execCreateTenant(s.admin, denom, strconv.FormatInt(int64(period), 10))

		afterTenants, err := queryTenants(chainAPIEndpoint)
		s.Require().NoError(err)
		s.Require().Equal(len(beforeTenants)+1, len(afterTenants))

		newTenant := afterTenants[len(afterTenants)-1]
		s.Require().Equal(denom, newTenant.Tenant.Denom)
		s.Require().Equal(period, int(newTenant.Tenant.PayoutPeriod))
		s.Require().Equal(types.PayoutMethod_Native, newTenant.Tenant.PayoutMethod)

		tenantId = newTenant.Tenant.Id
	})
	s.Require().True(pass)

	pass = s.Run("record_internal_nft_revenue", func() {
		requestId := NewReqId()
		s.execRecord(s.admin, tenantId, requestId, revenue.String(), chainId, s.internalNftAddr, internalNftId)

		utxr, err := queryUtxr(chainAPIEndpoint, tenantId, requestId)
		s.Require().NoError(err)
		s.Require().Equal(1, len(utxr.Recipients))
		s.Require().Equal(common.FromHex(internalNftOwner), utxr.Recipients[0].Address.Bytes())
		s.Require().Equal(revenue, utxr.Amount)
	})
	s.Require().True(pass)

	pass = s.Run("record_external_nft_revenue", func() {
		requestId := NewReqId()
		s.execRecord(s.admin, tenantId, requestId, revenue.String(), extChainId, extNftAddress, extNftId)
		utxr, err := queryUtxr(chainAPIEndpoint, tenantId, requestId)
		s.Require().NoError(err)
		s.Require().Equal(revenue, utxr.Amount)
		s.Require().Eventually(
			func() bool {
				utxr, err = queryUtxr(chainAPIEndpoint, tenantId, requestId)
				return err == nil && len(utxr.Recipients) > 0 && bytes.Equal(utxr.Recipients[0].Address.Bytes(), common.FromHex(extNftOwner))
			},
			time.Minute,
			time.Second,
		)
		s.Require().Equal(revenue, utxr.Amount)
	})
	s.Require().True(pass)

	pass = s.Run("deposit_to_treasury", func() {
		beforeBalance, err := getEthBalance(chainAPIEndpoint, internalNftOwner)
		s.Require().NoError(err)

		firstDeposit := revenue.SubAmount(sdk.NewInt(1))
		secondDeposit := sdk.NewCoin(denom, sdk.NewInt(1))
		s.execDepositToTreasury(s.admin, tenantId, firstDeposit.String())
		tenant, err := queryTenant(chainAPIEndpoint, tenantId)
		s.Require().NoError(err)
		s.Require().Equal(tenant.Treasury.Balance.Amount.Int64(), firstDeposit.Amount.Int64())

		s.execDepositToTreasury(s.admin, tenantId, secondDeposit.String())
		s.Require().Eventually(
			func() bool {
				tenant, err := queryTenant(chainAPIEndpoint, tenantId)
				s.Require().NoError(err)

				afterBalance, err := getEthBalance(chainAPIEndpoint, internalNftOwner)
				s.Require().NoError(err)
				return tenant.Treasury.Balance.Amount.IsZero() && afterBalance-beforeBalance == revenue.Amount.Uint64()
			},
			time.Minute,
			2*time.Second,
		)

		s.execDepositToTreasury(s.admin, tenantId, revenue.String())
		s.Require().Eventually(
			func() bool {
				tenant, err := queryTenant(chainAPIEndpoint, tenantId)
				s.Require().NoError(err)

				afterBalance, err := getEthBalance(chainAPIEndpoint, internalNftOwner)
				s.Require().NoError(err)

				return tenant.Treasury.Balance.Amount.IsZero() && afterBalance-beforeBalance == revenue.Amount.Uint64()
			},
			time.Minute,
			2*time.Second,
		)
	})
	s.Require().True(pass)
}

func (s *IntegrationTestSuite) TestMintableContractTenant() {
	denom := "eBLUC"
	revenue := sdk.NewCoin(denom, sdk.NewInt(10))
	tenantContractAddr := "0x0"
	var tenantId uint64 = 0

	pass := s.Run("create_new_tenant", func() {
		period := 8
		beforeTenants, err := queryTenants(chainAPIEndpoint)
		s.Require().NoError(err)

		s.execCreateMcTenant(s.admin, denom, strconv.FormatInt(int64(period), 10))
		afterTenants, err := queryTenants(chainAPIEndpoint)
		s.Require().NoError(err)

		s.Require().Equal(len(beforeTenants)+1, len(afterTenants))

		newTenant := afterTenants[len(afterTenants)-1]
		s.Require().Equal(denom, newTenant.Tenant.Denom)
		s.Require().Equal(period, int(newTenant.Tenant.PayoutPeriod))
		s.Require().Equal(types.PayoutMethod_MintContract, newTenant.Tenant.PayoutMethod)
		s.Require().NotEmpty(newTenant.Tenant.ContractAddress)

		tenantId = newTenant.Tenant.Id
		tenantContractAddr = newTenant.Tenant.ContractAddress
	})
	s.Require().True(pass)

	pass = s.Run("record_internal_nft_revenue", func() {
		requestId := NewReqId()
		s.execRecord(s.admin, tenantId, requestId, revenue.String(), chainId, s.internalNftAddr, internalNftId)

		utxr, err := queryUtxr(chainAPIEndpoint, tenantId, requestId)
		s.Require().NoError(err)
		s.Require().Equal(common.FromHex(internalNftOwner), utxr.Recipients[0].Address.Bytes())
		s.Require().Equal(revenue, utxr.Amount)

		beforeBalance, err := queryERC20Balance(s.ethClient, tenantContractAddr, internalNftOwner)
		s.Require().NoError(err)

		s.Require().Eventually(
			func() bool {
				afterBalance, err := queryERC20Balance(s.ethClient, tenantContractAddr, internalNftOwner)
				s.Require().NoError(err)
				return afterBalance-beforeBalance == revenue.Amount.Uint64()
			},
			time.Minute,
			2*time.Second,
		)
	})
	s.Require().True(pass)

	pass = s.Run("record_external_nft_revenue", func() {
		beforeBalance, err := queryERC20Balance(s.ethClient, tenantContractAddr, extNftOwner)
		s.Require().NoError(err)

		requestId := NewReqId()
		s.execRecord(s.admin, tenantId, requestId, revenue.String(), extChainId, extNftAddress, extNftId)
		utxr, err := queryUtxr(chainAPIEndpoint, tenantId, requestId)
		s.Require().NoError(err)
		s.Require().Equal(revenue, utxr.Amount)
		s.Require().Eventually(
			func() bool {
				utxr, err = queryUtxr(chainAPIEndpoint, tenantId, requestId)
				return err == nil && len(utxr.Recipients) > 0 && bytes.Equal(utxr.Recipients[0].Address.Bytes(), common.FromHex(extNftOwner))
			},
			time.Minute,
			time.Second,
		)

		s.Require().Eventually(
			func() bool {
				afterBalance, err := queryERC20Balance(s.ethClient, tenantContractAddr, extNftOwner)
				s.Require().NoError(err)
				return afterBalance-beforeBalance == revenue.Amount.Uint64()
			},
			time.Minute,
			2*time.Second,
		)
	})
	s.Require().True(pass)
}

func NewReqId() string {
	return "req-" + strconv.Itoa(rand.Intn(10000000000))
}
