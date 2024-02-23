package repository

import "github.com/settlus/chain/tools/interop-node/types"

type BlockData struct {
	Number []byte
	Hash   []byte
}

type Repository interface {
	PutBlockData(hash []byte, number []byte, timestamp uint64, ownershipData []*types.OwnershipTransferEvent) error

	GetRecentBlock(timestamp uint64) (BlockData, error)
	GetBlockNumber(blockNumber string) (string, error)
	GetBlockHash(blockHash string) (string, error)
	GetNftOwnership(nftAddessHex string, tokenIdHex string, blockNumberHex string) (string, error)
}
