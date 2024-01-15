package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	PayoutMethod_Native       = "native"
	PayoutMethod_MintContract = "mintable_contract"
)

func GetTreasuryAccountName(tenantId uint64) string {
	return fmt.Sprintf("tenant-%d", tenantId)
}

// GetTenantTreasuryAccount gets the treasury account for the tenant by hashing the tenant account name
func GetTenantTreasuryAccount(tenantId uint64) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(GetTreasuryAccountName(tenantId))))
}

// Tenant is a struct that contains the information of a tenant
func (t Tenant) String() string {
	return fmt.Sprintf(`Tenant:
	TenantId: %d
	Admins: %s
	Denom: %s
	PayoutPeriod: %d
	PayoutMethod: %s
	ContractAddress: %s
`, t.Id, t.Admins, t.Denom, t.PayoutPeriod, t.PayoutMethod, t.ContractAddress)
}
