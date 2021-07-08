package main

import (
	"fmt"

	"github.com/Link512/stealthpool"
)

func main() {
	pool, err := stealthpool.New(2, stealthpool.WithBlockSize(8*1024*1024))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	block, err := pool.Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("len(block): %d cap(block): %d \n", len(block), cap(block))
}
