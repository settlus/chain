package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/settlus/chain/contracts"
	"github.com/settlus/chain/x/settlement/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTenantStore returns the tenant store for the given tenantId
func (k SettlementKeeper) GetTenantStore(ctx sdk.Context) sdk.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.TenantPrefix)
}

// GetLargestTenantId returns the latest tenant from the store
func (k SettlementKeeper) GetLargestTenantId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	iterator := prefix.NewStore(store, types.TenantPrefix).ReverseIterator(nil, nil)
	defer iterator.Close()

	if !iterator.Valid() {
		// if there is no tenant, return 0
		return 0
	}

	return sdk.BigEndianToUint64(iterator.Key())
}

// GetTenant returns the tenant by its tenantId
func (k SettlementKeeper) GetTenant(ctx sdk.Context, tenantId uint64) *types.Tenant {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TenantStoreKey(tenantId))
	if bz == nil {
		return nil
	}

	var tenant types.Tenant
	k.cdc.MustUnmarshal(bz, &tenant)
	return &tenant
}

// GetTenantTreasury returns the tenant treasury account and its balance
func (k SettlementKeeper) GetTenantTreasury(ctx sdk.Context, tenantId uint64) (sdk.AccAddress, sdk.Coins) {
	tenantTreasury := types.GetTenantTreasuryAccount(tenantId)
	treasuryBalance := k.bk.SpendableCoins(ctx, tenantTreasury)
	return tenantTreasury, treasuryBalance
}

// deployTokenContract deploys a soul bound token (SBT) contract on the EVM
func (k SettlementKeeper) deployTokenContract(ctx sdk.Context, tenantId uint64, denom string) (string, error) {
	tenantAccount := types.GetTenantTreasuryAccount(tenantId)
	tenantAddress := common.BytesToAddress(tenantAccount)
	ctorArgs, err := contracts.SBTContract.ABI.Pack(
		"",
		denom,
		denom,
		tenantAddress,
	)
	if err != nil {
		return "", err
	}

	data := make([]byte, len(contracts.SBTContract.Bin)+len(ctorArgs))
	copy(data[:len(contracts.SBTContract.Bin)], contracts.SBTContract.Bin)
	copy(data[len(contracts.SBTContract.Bin):], ctorArgs)

	nonce, err := k.ak.GetSequence(ctx, tenantAccount)
	if err != nil {
		return "", err
	}

	contractAddr := crypto.CreateAddress(tenantAddress, nonce)
	_, err = k.evmk.CallEVMWithData(ctx, tenantAddress, nil, data, true)
	if err != nil {
		return "", err
	}

	contractAddressHex := contractAddr.Hex()

	err = ctx.EventManager().EmitTypedEvent(&types.EventDeployContract{
		Tenant:          tenantId,
		ContractAddress: contractAddressHex,
		TokenName:       denom,
		ContractAdmin:   tenantAddress.Hex(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to emit event: deploy contract: (%s)", err)
	}

	return contractAddressHex, nil
}

// CreateNewTenant creates a new tenant
func (k SettlementKeeper) CreateNewTenant(ctx sdk.Context, initialAdmin string, denom string, payoutPeriod uint64, payoutMethod string, contractAddr string) (tenantId uint64, err error) {
	tenantId = k.GetLargestTenantId(ctx) + 1
	tenant := &types.Tenant{
		Id:              tenantId,
		Admins:          []string{initialAdmin},
		Denom:           denom,
		PayoutPeriod:    payoutPeriod,
		PayoutMethod:    payoutMethod,
		ContractAddress: contractAddr,
	}

	k.CreateTreasuryAccount(ctx, tenantId)

	if payoutMethod == types.PayoutMethod_MintContract && contractAddr == "" {
		contractAddr, err = k.deployTokenContract(ctx, tenantId, tenant.Denom)
		if err != nil {
			return 0, err
		}
		tenant.ContractAddress = contractAddr
	}

	k.SetTenant(ctx, tenant)

	return tenantId, nil
}

// CheckAdminPermission returns true if the account has admin permission
func (k SettlementKeeper) CheckAdminPermission(ctx sdk.Context, tenantId uint64, account string) bool {
	tenant := k.GetTenant(ctx, tenantId)
	if tenant == nil {
		return false
	}

	for _, admin := range tenant.Admins {
		if admin == account {
			return true
		}
	}

	return false
}

// GetPayoutPeriod returns the payout period of the tenant
func (k SettlementKeeper) GetPayoutPeriod(ctx sdk.Context, tenantId uint64) uint64 {
	tenant := k.GetTenant(ctx, tenantId)
	if tenant == nil {
		return 0
	}

	return tenant.PayoutPeriod
}

// CheckTenantExist returns true if tenant exists
func (k SettlementKeeper) CheckTenantExist(ctx sdk.Context, tenantId uint64) bool {
	return k.GetTenant(ctx, tenantId) != nil
}

// SetTenant sets the tenant to the store
func (k SettlementKeeper) SetTenant(ctx sdk.Context, tenant *types.Tenant) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(tenant)
	store.Set(types.TenantStoreKey(tenant.Id), bz)
}

// GetAllTenants returns all tenants
func (k SettlementKeeper) GetAllTenants(ctx sdk.Context) []types.Tenant {
	store := ctx.KVStore(k.storeKey)
	iterator := prefix.NewStore(store, types.TenantPrefix).Iterator(nil, nil)
	defer iterator.Close()

	var tenants []types.Tenant
	for ; iterator.Valid(); iterator.Next() {
		var tenant types.Tenant
		k.cdc.MustUnmarshal(iterator.Value(), &tenant)
		tenants = append(tenants, tenant)
	}

	return tenants
}

func (k SettlementKeeper) CreateTreasuryAccount(ctx sdk.Context, tenantId uint64) {
	k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, types.GetTenantTreasuryAccount(tenantId)))
}
