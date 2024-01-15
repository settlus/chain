package repository

import "github.com/settlus/chain/tools/interop-node/types"

type BlockData struct {
	Number []byte
	Hash   []byte
}

type Repository interface {
	PutBlockData(hash []byte, number []byte, ownershipData []*types.OwnershipTransferEvent) error

	GetRecentBlock() (BlockData, error)
	GetBlockNumber(blockNumber string) (string, error)
	GetBlockHash(blockHash string) (string, error)
	GetNftOwnership(nftAddessHex string, tokenIdHex string, blockNumberHex string) (string, error)
}
