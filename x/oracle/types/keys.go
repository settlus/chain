package types

const (
	// ModuleName defines the module name
	ModuleName = "oracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	BlockDataKeyPrefix        = []byte{0x00}
	FeederDelegationKeyPrefix = []byte{0x01}
	MissCountKeyPrefix        = []byte{0x02}
	AggregatePrevoteKeyPrefix = []byte{0x03}
	AggregateVoteKeyPrefix    = []byte{0x04}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func BlockDataKey(chainId string) []byte {
	return append(BlockDataKeyPrefix, KeyPrefix(chainId)...)
}

func FeederDelegationKey(validatorAddress string) []byte {
	return append(FeederDelegationKeyPrefix, validatorAddress...)
}

func MissCountKey(validatorAddress string) []byte {
	return append(MissCountKeyPrefix, validatorAddress...)
}

func AggregatePrevoteKey(validatorAddress string) []byte {
	return append(AggregatePrevoteKeyPrefix, validatorAddress...)
}

func AggregateVoteKey(validatorAddress string) []byte {
	return append(AggregateVoteKeyPrefix, validatorAddress...)
}
