package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/ethereum/go-ethereum/common"

	"github.com/settlus/chain/x/settlement/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetUTXRStore returns the UTXR store for the given tenantId
func (k SettlementKeeper) GetUTXRStore(ctx sdk.Context, tenantId uint64) sdk.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.UTXRStoreByTenantKey(tenantId))
}

// HasUTXRByRequestId returns whether the UTXR exists for the given tenantId and requestId
func (k SettlementKeeper) HasUTXRByRequestId(ctx sdk.Context, tenantId uint64, requestId string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.UTXRStoreByRequestIdKey(tenantId, requestId))
}

// GetLargestUTXRId returns the latest UTXR from the store by its tenantId
func (k SettlementKeeper) GetLargestUTXRId(ctx sdk.Context, tenantId uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	utxrTenantStore := prefix.NewStore(store, types.UTXRStoreByTenantKey(tenantId))
	iterator := utxrTenantStore.ReverseIterator(nil, nil)
	defer iterator.Close()

	if !iterator.Valid() {
		// if there is no UTXR, return 0
		return 0
	}

	return sdk.BigEndianToUint64(iterator.Key())
}

// CreateUTXR creates a new UTXR in the store
func (k SettlementKeeper) CreateUTXR(ctx sdk.Context, tenantId uint64, utxr *types.UTXR) (uint64, error) {
	if k.HasUTXRByRequestId(ctx, tenantId, utxr.RequestId) {
		return 0, sdkerrors.Wrapf(types.ErrDuplicateRequestId, "UTXR with [request ID: %s] [tenant ID: %d] already exists.", utxr.RequestId, tenantId)
	}

	if !common.IsHexAddress(utxr.Recipient.String()) {
		return 0, sdkerrors.Wrapf(types.ErrInvalidAccount, "Invalid recipient address: %s", utxr.Recipient)
	}

	recipient := common.HexToAddress(utxr.Recipient.String()).Bytes()
	if !k.ak.HasAccount(ctx, recipient) {
		k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, recipient))
	}

	utxrId := k.GetLargestUTXRId(ctx, tenantId) + 1

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(utxr)
	store.Set(types.UTXRStoreKey(tenantId, utxrId), bz)
	store.Set(types.UTXRStoreByRequestIdKey(tenantId, utxr.RequestId), sdk.Uint64ToBigEndian(utxrId))

	return utxrId, nil
}

// DeleteUTXR deletes UTXR and UTXR ID from the store by its tenantId and utxrId
func (k SettlementKeeper) deleteUTXR(ctx sdk.Context, tenantId, utxrId uint64) error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UTXRStoreKey(tenantId, utxrId))
	if bz == nil {
		return sdkerrors.Wrapf(types.ErrUTXRNotFound, "UTXR with [tenant ID: %d] [utxr ID: %d] not found.", tenantId, utxrId)
	}

	var utxr types.UTXR
	k.cdc.MustUnmarshal(bz, &utxr)

	store.Delete(types.UTXRStoreKey(tenantId, utxrId))
	store.Delete(types.UTXRStoreByRequestIdKey(tenantId, utxr.RequestId))
	return nil
}

// GetUTXRByRequestId returns a UTXR from the store by its tenantId and requestId
func (k SettlementKeeper) GetUTXRByRequestId(ctx sdk.Context, tenantId uint64, requestId string) *types.UTXR {
	store := ctx.KVStore(k.storeKey)
	bzUTXRId := store.Get(types.UTXRStoreByRequestIdKey(tenantId, requestId))
	if bzUTXRId == nil {
		return nil
	}
	utxrId := sdk.BigEndianToUint64(bzUTXRId)

	bzUtxr := store.Get(types.UTXRStoreKey(tenantId, utxrId))
	if bzUtxr == nil {
		return nil
	}

	var utxr types.UTXR
	k.cdc.MustUnmarshal(bzUtxr, &utxr)
	return &utxr
}

// DeleteUTXRByRequestId deletes UTXR and UTXR ID from the store by its tenantId and requestId
func (k SettlementKeeper) DeleteUTXRByRequestId(ctx sdk.Context, tenantId uint64, requestId string) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UTXRStoreByRequestIdKey(tenantId, requestId))
	if bz == nil {
		return 0, sdkerrors.Wrapf(types.ErrUTXRNotFound, "UTXR with [tenant ID: %d] [request ID: %s] not found.", tenantId, requestId)
	}

	utxrId := sdk.BigEndianToUint64(bz)

	store.Delete(types.UTXRStoreKey(tenantId, utxrId))
	store.Delete(types.UTXRStoreByRequestIdKey(tenantId, requestId))
	return utxrId, nil
}

// GetAllUTXRWithTenantAndID returns all UTXRs with tenantId and utxrId
func (k SettlementKeeper) GetAllUTXRWithTenantAndID(ctx sdk.Context) (list []types.UTXRWithTenantAndId) {
	store := ctx.KVStore(k.storeKey)
	utxrStore := prefix.NewStore(store, types.UTXRPrefix)
	iterator := utxrStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var utxr types.UTXR
		k.cdc.MustUnmarshal(iterator.Value(), &utxr)

		key := iterator.Key()
		// 8 bytes for tenantId, 8 bytes for utxrId
		tenantId := sdk.BigEndianToUint64(key[0:8])
		utxrId := sdk.BigEndianToUint64(key[8:])

		list = append(list, types.UTXRWithTenantAndId{
			TenantId: tenantId,
			Id:       utxrId,
			Utxr:     utxr,
		})
	}

	return list
}
