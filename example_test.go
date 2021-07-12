package stealthpool_test

import (
	"errors"
	"fmt"

	"github.com/Link512/stealthpool"
)

func ExampleNew_customBlockSize() {
	pool, err := stealthpool.New(2, stealthpool.WithBlockSize(8*1024*1024))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	block, err := pool.Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("len(block): %d cap(block): %d\n", len(block), cap(block))
	// Output:
	// len(block): 8388608 cap(block): 8388608
}

func copyStrToSlice(b []byte, s string) []byte {
	return append(b, []byte(s)...)
}

func ExamplePool_Get() {
	pool, err := stealthpool.New(1)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	block, err := pool.Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("len(block): %d cap(block): %d\n", len(block), cap(block))
	fmt.Printf("free blocks: %d total allocated: %d\n", pool.FreeCount(), pool.AllocCount())

	_, err = pool.Get()
	fmt.Printf("pool is full: %t\n", errors.Is(stealthpool.ErrPoolFull, err))

	// Output:
	// len(block): 4096 cap(block): 4096
	// free blocks: 0 total allocated: 1
	// pool is full: true
}

func ExamplePool_Return() {
	pool, err := stealthpool.New(1)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	block, err := pool.Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("free blocks: %d total allocated: %d\n", pool.FreeCount(), pool.AllocCount())

	err = pool.Return(block[1:])
	fmt.Printf("resliced block is invalid: %t\n", errors.Is(stealthpool.ErrInvalidBlock, err))

	// once block is returned, the pool will re-use it in a subsequent Get call
	err = pool.Return(block)
	if err != nil {
		panic(err)
	}
	fmt.Printf("free blocks: %d total allocated: %d\n", pool.FreeCount(), pool.AllocCount())

	block, err = pool.Get()
	if err != nil {
		panic(err)
	}

	// Output:
	// free blocks: 0 total allocated: 1
	// resliced block is invalid: true
	// free blocks: 1 total allocated: 1
}

func ExamplePool_Close() {
	pool, err := stealthpool.New(2)
	if err != nil {
		panic(err)
	}

	block1, err := pool.Get()
	if err != nil {
		panic(err)
	}

	block2, err := pool.Get()
	if err != nil {
		panic(err)
	}

	if err := pool.Close(); err != nil {
		panic(err)
	}

	// using the slices after Close() will panic
	copyStrToSlice(block1, "Hello")
	copyStrToSlice(block2, "World")
}
