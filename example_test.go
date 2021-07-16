package stealthpool_test

import (
	"errors"
	"fmt"

	"github.com/Link512/stealthpool"
)

func ExampleNew() {
	pool, err := stealthpool.New(2)
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
	// len(block): 4096 cap(block): 4096
}

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

func ExampleNew_preallocation() {
	pool, err := stealthpool.New(2, stealthpool.WithPreAlloc(2))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	fmt.Printf("free blocks: %d total allocated: %d\n", pool.FreeCount(), pool.AllocCount())

	block, err := pool.Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("len(block): %d cap(block): %d\n", len(block), cap(block))

	// Output:
	// free blocks: 2 total allocated: 2
	// len(block): 4096 cap(block): 4096
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

	newBlock, err := pool.Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("block was reused: %t\n", &block[0] == &newBlock[0])

	// Output:
	// free blocks: 0 total allocated: 1
	// resliced block is invalid: true
	// free blocks: 1 total allocated: 1
	// block was reused: true
}

func ExamplePool_AllocCount() {
	pool, err := stealthpool.New(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("total allocated: %d\n", pool.AllocCount())

	_, _ = pool.Get()
	fmt.Printf("total allocated: %d\n", pool.AllocCount())

	pool.Close()

	newPool, err := stealthpool.New(3, stealthpool.WithPreAlloc(2))
	if err != nil {
		panic(err)
	}
	fmt.Printf("total allocated: %d\n", newPool.AllocCount())
	newPool.Close()

	// Output:
	// total allocated: 0
	// total allocated: 1
	// total allocated: 2
}

func ExamplePool_FreeCount() {
	pool, err := stealthpool.New(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("free blocks: %d\n", pool.FreeCount())

	block, _ := pool.Get()
	fmt.Printf("free blocks: %d\n", pool.FreeCount())

	_ = pool.Return(block)
	fmt.Printf("free blocks: %d\n", pool.FreeCount())
	pool.Close()

	newPool, err := stealthpool.New(3, stealthpool.WithPreAlloc(2))
	if err != nil {
		panic(err)
	}
	fmt.Printf("free blocks: %d\n", newPool.FreeCount())
	newPool.Close()

	// Output:
	// free blocks: 0
	// free blocks: 0
	// free blocks: 1
	// free blocks: 2
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
	fmt.Println(block1[:5])
	fmt.Println(block2[3:])
}
