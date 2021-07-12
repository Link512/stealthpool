package stealthpool

import "errors"

var (
	// ErrPoolFull is returned when the maximum number of allocated blocks has been reached
	ErrPoolFull = errors.New("pool is full")
	// ErrPreallocOutOfBounds is returned when whe number of preallocated blocks requested is either negative or above maxBlocks
	ErrPreallocOutOfBounds = errors.New("prealloc value out of bounds")
	// ErrInvalidBlock is returned when an invalid slice is passed to Return()
	ErrInvalidBlock = errors.New("trying to return invalid block")
)
