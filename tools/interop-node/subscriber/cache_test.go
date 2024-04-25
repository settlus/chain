package subscriber

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockCache_PutBlockData(t *testing.T) {
	c := NewBlockCache(2)

	// Add some initial data to the cache
	c.PutBlockData("hash1", 1, 123456789)
	c.PutBlockData("hash2", 2, 123456790)

	// Add a new block data
	c.PutBlockData("hash3", 3, 123456791)

	require.Equal(t, 2, c.data.Size())

	_, ok := c.data.Get(uint64(123456789))
	require.False(t, ok)

	_, ok = c.data.Get(uint64(123456790))
	require.True(t, ok)

	_, ok = c.data.Get(uint64(123456791))
	require.True(t, ok)
}

func TestBlockCache_GetOldestBlock(t *testing.T) {
	c := NewBlockCache(2)

	// Add some initial data to the cache
	c.PutBlockData("hash1", 1, 123456789)
	c.PutBlockData("hash2", 2, 123456790)

	hash, _ := c.GetOldestBlock(123456788)
	require.Equal(t, "hash1", hash)

	hash, _ = c.GetOldestBlock(123456790)
	require.Equal(t, "hash2", hash)

	hash, _ = c.GetOldestBlock(123456791)
	require.Equal(t, "", hash)
}
