package types

const (
	// ModuleName defines the module name
	ModuleName = "oracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// TransientKey is the key to access the Oracle transient store, that is reset
	// during the Commit phase.
	TransientKey = "transient_" + ModuleName
)

var (
	BlockDataKeyPrefix        = []byte{0x00}
	FeederDelegationKeyPrefix = []byte{0x01}
	MissCountKeyPrefix        = []byte{0x02}
	AggregatePrevoteKeyPrefix = []byte{0x03}
	AggregateVoteKeyPrefix    = []byte{0x04}
	RoundKeyPrefix            = []byte{0x05}

	CurrentRoundBytesTransientKeyPrefix = []byte{0x00}
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
