package voteprocessor

import sdk "github.com/cosmos/cosmos-sdk/types"

type DataWithVoter[T any] struct {
	Data  T
	Voter sdk.ValAddress
}

type DataWithWeight[T any] struct {
	Data   T
	Weight int64
}

type ConsensusHook[Source comparable, Data comparable] func(sdk.Context, map[Source]Data)
type DataConverter[Source comparable, Data comparable] func(string) (Source, Data, error)
