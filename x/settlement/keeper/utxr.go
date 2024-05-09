package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/ethereum/go-ethereum/common"

	"github.com/settlus/chain/x/settlement/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/settlus/chain/types"
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

// GenerateUtxrId returns new UTXR ID for each tenant
func (k SettlementKeeper) GenerateUtxrId(ctx sdk.Context, tenantId uint64) uint64 {
	var utxrId uint64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastUtxrIdStoreKey(tenantId))
	if bz != nil {
		utxrId = sdk.BigEndianToUint64(bz)
		utxrId++
	}

	store.Set(types.LastUtxrIdStoreKey(tenantId), sdk.Uint64ToBigEndian(utxrId))
	return utxrId
}

// CreateUTXR creates a new UTXR in the store
func (k SettlementKeeper) CreateUTXR(ctx sdk.Context, tenantId uint64, utxr *types.UTXR) (uint64, error) {
	if k.HasUTXRByRequestId(ctx, tenantId, utxr.RequestId) {
		return 0, sdkerrors.Wrapf(types.ErrDuplicateRequestId, "UTXR with [request ID: %s] [tenant ID: %d] already exists.", utxr.RequestId, tenantId)
	}

	for _, recipient := range utxr.Recipients {
		if !common.IsHexAddress(recipient.Address.String()) {
			return 0, sdkerrors.Wrapf(types.ErrInvalidAccount, "Invalid recipient address: %s", recipient.Address)
		}

		accAddr := sdk.AccAddress(recipient.Address.Bytes())
		if !k.ak.HasAccount(ctx, accAddr) {
			k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, accAddr))
		}
	}

	utxrId := k.GenerateUtxrId(ctx, tenantId)

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

func (k SettlementKeeper) GetAllUniqueNftToVerify(ctx sdk.Context, until uint64) (list []ctypes.Nft) {
	nfts := make(map[ctypes.Nft]struct{})

	store := ctx.KVStore(k.storeKey)
	utxrStore := prefix.NewStore(store, types.UTXRPrefix)
	iterator := utxrStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var utxr types.UTXR
		k.cdc.MustUnmarshal(iterator.Value(), &utxr)
		if len(utxr.Recipients) == 0 && until >= utxr.CreatedAt {
			nfts[*utxr.Nft] = struct{}{}
		}
	}

	for nft := range nfts {
		list = append(list, nft)
	}

	return list
}

func (k SettlementKeeper) SetRecipients(ctx sdk.Context, nfts map[ctypes.Nft]ctypes.HexAddressString, until uint64) {
	store := ctx.KVStore(k.storeKey)
	utxrStore := prefix.NewStore(store, types.UTXRPrefix)
	iterator := utxrStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var utxr types.UTXR
		k.cdc.MustUnmarshal(iterator.Value(), &utxr)
		if len(utxr.Recipients) == 0 && until >= utxr.CreatedAt {
			owner, ok := nfts[*utxr.Nft]
			if !ok {
				continue
			}
			utxr.Recipients = []*types.Recipient{{
				Address: owner,
			}}
			bz := k.cdc.MustMarshal(&utxr)
			key := iterator.Key()
			utxrStore.Set(key, bz)

			err := ctx.EventManager().EmitTypedEvents(&types.EventSetRecipients{
				Tenant:     sdk.BigEndianToUint64(key[0:8]),
				UtxrId:     sdk.BigEndianToUint64(key[8:]),
				Recipients: utxr.Recipients,
			})
			if err != nil {
				k.Logger(ctx).Error("failed to emit EventSetRecipients", "error", err)
			}
		}
	}
}
