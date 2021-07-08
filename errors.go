package stealthpool

import "errors"

var (
	ErrPoolFull            = errors.New("pool is full")
	ErrPreallocOutOfBounds = errors.New("prealloc value out of bounds")
	ErrInvalidBlock        = errors.New("trying to return invalid block")
)
