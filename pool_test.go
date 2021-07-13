package stealthpool_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Link512/stealthpool"
)

func TestPool_New(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	pool, err := stealthpool.New(3, stealthpool.WithBlockSize(8*1024))
	require.NoError(err)

	defer pool.Close()

	assert.Equal(0, pool.FreeCount())
	assert.Equal(0, pool.AllocCount())

	block, err := pool.Get()
	require.NoError(err)
	assert.Len(block, 8*1024)
}

func TestPool_New_prealloc(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	_, err := stealthpool.New(3, stealthpool.WithBlockSize(8*1024), stealthpool.WithPreAlloc(4))
	require.ErrorIs(err, stealthpool.ErrPreallocOutOfBounds)

	pool, err := stealthpool.New(3, stealthpool.WithBlockSize(8*1024), stealthpool.WithPreAlloc(2))
	require.NoError(err)

	defer pool.Close()

	assert.Equal(2, pool.FreeCount())
	assert.Equal(2, pool.AllocCount())

	block, err := pool.Get()
	require.NoError(err)
	assert.Len(block, 8*1024)

	assert.Equal(1, pool.FreeCount())
	assert.Equal(2, pool.AllocCount())
}

func TestPool_Get(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	pool, err := stealthpool.New(3, stealthpool.WithBlockSize(8*1024), stealthpool.WithPreAlloc(2))
	require.NoError(err)

	defer pool.Close()

	assert.Equal(2, pool.FreeCount())
	assert.Equal(2, pool.AllocCount())

	block, err := pool.Get()
	require.NoError(err)
	assert.Len(block, 8*1024)

	assert.Equal(1, pool.FreeCount())
	assert.Equal(2, pool.AllocCount())

	block, err = pool.Get()
	require.NoError(err)
	assert.Len(block, 8*1024)

	assert.Equal(0, pool.FreeCount())
	assert.Equal(2, pool.AllocCount())

	block, err = pool.Get()
	require.NoError(err)
	assert.Len(block, 8*1024)

	assert.Equal(0, pool.FreeCount())
	assert.Equal(3, pool.AllocCount())

	_, err = pool.Get()
	require.ErrorIs(err, stealthpool.ErrPoolFull)
}

func TestPool_Return(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	pool, err := stealthpool.New(2, stealthpool.WithBlockSize(8*1024), stealthpool.WithPreAlloc(2))
	require.NoError(err)

	defer pool.Close()

	assert.Equal(2, pool.FreeCount())
	assert.Equal(2, pool.AllocCount())

	block, err := pool.Get()
	require.NoError(err)
	assert.Len(block, 8*1024)

	block2, err := pool.Get()
	require.NoError(err)
	assert.Len(block2, 8*1024)

	assert.Equal(0, pool.FreeCount())
	assert.Equal(2, pool.AllocCount())

	err = pool.Return(block[2:])
	assert.ErrorIs(err, stealthpool.ErrInvalidBlock)

	err = pool.Return(block2)
	require.NoError(err)
	assert.Equal(1, pool.FreeCount())
	assert.Equal(2, pool.AllocCount())

	block3, err := pool.Get()
	require.NoError(err)
	assert.Len(block2, 8*1024)
	assert.Equal(block2, block3)
}

func TestPool_CanAccessMemory(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	pool, err := stealthpool.New(1, stealthpool.WithBlockSize(256))
	require.NoError(err)
	defer pool.Close()

	block, err := pool.Get()
	require.NoError(err)

	payload := make([]byte, 256)
	for i := range payload {
		payload[i] = byte(i)
	}
	copy(block, payload)
	assert.Equal(payload, block)
}
