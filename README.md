# Stealthpool

[![Build Status](https://cloud.drone.io/api/badges/Link512/stealthpool/status.svg)](https://cloud.drone.io/Link512/stealthpool) [![Go Report Card](https://goreportcard.com/badge/github.com/Link512/stealthpool)](https://goreportcard.com/report/github.com/Link512/stealthpool) [![Go Reference](https://pkg.go.dev/badge/github.com/Link512/stealthpool.svg)](https://pkg.go.dev/github.com/Link512/stealthpool)

stealthpool provides a memory pool that allocates blocks off-heap that will NOT be tracked by the garbage collector.
The name stealthpool comes from the fact that the memory being allocated by the pool is `stealthy` and will not be garbage collected ever
These blocks should be used in situations where you want to keep garbage collection to a bare minimum.
Needless to say, since the GC will not track any of the memory, extreme care must be had in order to avoid memory leaks

## Installation

```bash
go get -u github.com/Link512/stealthpool
```

## Getting started

```golang

// initialize a pool which will allocate a maximum of 100 blocks
pool, err := stealthpool.New(100)
defer pool.Close() // ALWAYS close the pool unless you're very fond of memory leaks

// initialize a pool with custom block size and preallocated blocks
poolCustom, err := stealthpool.New(100, stealthpool.WithBlockSize(8*1024), stealthpool.WithPreAlloc(100))
defer poolCustom.Close() // ALWAYS close the pool unless you're very fond of memory leaks


block, err := poolCustom.Get()
// do some work with block
// then return it exactly as-is to the pool
err = poolCustom.Return(block)
```

## Docs

Go docs together with examples can be found [here](https://pkg.go.dev/github.com/Link512/stealthpool)
