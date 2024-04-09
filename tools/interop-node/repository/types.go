package repository

import "github.com/settlus/chain/tools/interop-node/types"

type BlockData struct {
	Number []byte
	Hash   []byte
}

type Repository interface {
	PutBlockData(hash []byte, number []byte, timestamp uint64, ownershipData []*types.OwnershipTransferEvent) error

	GetOldestBlock(timestamp uint64) (BlockData, error)
	GetBlockNumber(blockHash string) (string, error)
	GetBlockHash(blockNumber string) (string, error)
	GetNftOwnership(nftAddessHex string, tokenIdHex string, blockNumberHex string) (string, error)
}
