package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "settlement"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_settlement"
)

var (
	UTXRPrefix          = []byte{0x00}
	UTXRRequestIdPrefix = []byte{0x01}
	TenantPrefix        = []byte{0x02}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func UTXRStoreByTenantKey(tenantId uint64) []byte {
	return append(UTXRPrefix, sdk.Uint64ToBigEndian(tenantId)...)
}

func UTXRStoreKey(tenantId, utxrId uint64) []byte {
	return append(UTXRStoreByTenantKey(tenantId), sdk.Uint64ToBigEndian(utxrId)...)
}

func UTXRStoreByRequestIdKey(tenantId uint64, requestId string) []byte {
	utxrTenantRequestIdPrefix := append(UTXRRequestIdPrefix, sdk.Uint64ToBigEndian(tenantId)...)
	return append(utxrTenantRequestIdPrefix, KeyPrefix(requestId)...)
}

func TenantStoreKey(tenantId uint64) []byte {
	return append(TenantPrefix, sdk.Uint64ToBigEndian(tenantId)...)
}
