package subscriber

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

type BlockCache struct {
	data *treemap.Map
	size int
}

type BlockData struct {
	Hash   string
	Number int64
}

func NewBlockCache(size int) *BlockCache {
	return &BlockCache{
		data: treemap.NewWith(utils.UInt64Comparator),
		size: size,
	}
}

func (c BlockCache) PutBlockData(hash string, number int64, timestamp uint64) {
	c.data.Put(timestamp, BlockData{
		Hash:   hash,
		Number: number,
	})

	if c.data.Size() > c.size {
		key, _ := c.data.Min()
		c.data.Remove(key)
	}
}

func (c BlockCache) GetOldestBlock(timetsamp uint64) (string, int64) {
	key, value := c.data.Ceiling(timetsamp)
	if key == nil {
		return "", 0
	}

	bd := value.(BlockData)
	return bd.Hash, bd.Number
}
