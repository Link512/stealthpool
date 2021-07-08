package main

import (
	"errors"
	"fmt"

	"github.com/Link512/stealthpool"
)

func copyStrToSlice(b []byte, s string) []byte {
	return append(b, []byte(s)...)
}

func main() {
	// initialize pool with a capacity of only 1 block
	pool, err := stealthpool.New(1)
	if err != nil {
		panic(err)
	}

	block, err := pool.Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("len: %d cap:%d\n", len(block), cap(block))
	val := copyStrToSlice(block[:0], "hello")
	// val and block point to the same underlying array
	fmt.Printf("val: %s len(val): %d, cap(val): %d, len(block): %d, cap(block): %d addr(val): %p addr(block): %p \n",
		string(val),
		len(val), cap(val),
		len(block), cap(block),
		val, block,
	)

	// since block wasn't returned to the pool, there's no more free blocks left
	_, err = pool.Get()
	if !errors.Is(err, stealthpool.ErrPoolFull) {
		panic("uh oh, pool was supposed to be full")
	}

	// after we're done with block, return it to the pool to be with its friends
	pool.Return(block)

	// we now have 1 free block
	newBlock, err := pool.Get()
	if err != nil {
		panic(err)
	}

	val = copyStrToSlice(newBlock[:0], "world")
	// val and newBlock point to the same underlying array
	fmt.Printf("val: %s len(val): %d, cap(val): %d, len(newBlock): %d, cap(newBlock): %d addr(val): %p addr(newBlock): %p \n",
		string(val),
		len(val), cap(val),
		len(newBlock), cap(newBlock),
		val, newBlock,
	)

	// newBlock is used as a scratchspace
	val = copyStrToSlice(newBlock[:0], "foo")
	fmt.Printf("val: %s\n", val)
	fmt.Printf("newBlock: %s\n", newBlock)

	// reslicing newBlock and trying to return it fails
	bkp := newBlock
	newBlock = newBlock[1:]
	err = pool.Return(newBlock)
	if !errors.Is(err, stealthpool.ErrInvalidBlock) {
		panic("uh oh, this shouldn't have worked")
	}

	// returning a different block fails
	err = pool.Return([]byte{})
	if !errors.Is(err, stealthpool.ErrInvalidBlock) {
		panic("uh oh, this shouldn't have worked")
	}

	// returning block exactly as it was gotten
	err = pool.Return(bkp)
	if err != nil {
		panic(err)
	}

	// unless you're very fond of memory leaks, never forget to close the pool
	err = pool.Close()
	if err != nil {
		panic(err)
	}
}
