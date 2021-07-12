package stealthpool

import (
	"reflect"
	"sync"
	"unsafe"
)

type poolOpts struct {
	blockSize int
	preAlloc  int
}

var (
	defaultPoolOpts = poolOpts{
		blockSize: 4 * 1024,
	}
)

// PoolOpt is a configuration option for a stealthpool
type PoolOpt func(*poolOpts)

// WithPreAlloc specifies how many blocks the pool should preallocate on initialization. Default is 0.
func WithPreAlloc(prealloc int) PoolOpt {
	return func(opts *poolOpts) {
		opts.preAlloc = prealloc
	}
}

// WithBlockSize specifies the block size that will be returned. It is highly advised that this block size be a multiple of 4KB or whatever value
// `os.Getpagesize()`, since the mmap syscall returns page aligned memory
func WithBlockSize(blockSize int) PoolOpt {
	return func(opts *poolOpts) {
		opts.blockSize = blockSize
	}
}

type Pool struct {
	free   [][]byte
	freeMu *sync.Mutex

	allocated   map[*byte]struct{}
	allocatedMu *sync.Mutex

	initOpts  poolOpts
	maxBlocks int
}

// New returns a new stealthpool with the given capacity. The configuration options can be used to change how many blocks are preallocated or block size.
// If preallocation fails (out of memory, etc), a cleanup of all previously preallocated will be attempted
func New(maxBlocks int, opts ...PoolOpt) (*Pool, error) {
	o := defaultPoolOpts
	for _, opt := range opts {
		opt(&o)
	}
	p := &Pool{
		initOpts:    o,
		free:        make([][]byte, 0, maxBlocks),
		freeMu:      &sync.Mutex{},
		allocated:   make(map[*byte]struct{}, maxBlocks),
		allocatedMu: &sync.Mutex{},
		maxBlocks:   maxBlocks,
	}
	if o.preAlloc > 0 {
		if err := p.prealloc(o.preAlloc); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// Get returns a memory block. It will first try and retrieve a previously allocated block and if that's not possible, will allocate a new block.
// If there were maxBlocks blocks already allocated, returns ErrPoolFull
func (p *Pool) Get() ([]byte, error) {
	if b, ok := p.tryPop(); ok {
		return b, nil
	}

	p.allocatedMu.Lock()
	defer p.allocatedMu.Unlock()

	if len(p.allocated) == p.maxBlocks {
		return nil, ErrPoolFull
	}
	result, err := alloc(p.initOpts.blockSize)
	if err != nil {
		return nil, err
	}
	k := &result[0]
	p.allocated[k] = struct{}{}
	return result, nil
}

// Return gives back a block retrieved from Get and stores it for future re-use.
// The block has to be exactly the same slice object returned from Get(), otherwise ErrInvalidBlock will be returned.
func (p *Pool) Return(b []byte) error {
	if err := p.checkValidBlock(b); err != nil {
		return err
	}
	p.freeMu.Lock()
	defer p.freeMu.Unlock()
	p.free = append(p.free, b)
	return nil
}

// FreeCount returns the number of free blocks that can be reused
func (p *Pool) FreeCount() int {
	p.freeMu.Lock()
	defer p.freeMu.Unlock()
	return len(p.free)
}

// AllocCount returns the total number of allocated blocks so far
func (p *Pool) AllocCount() int {
	p.allocatedMu.Lock()
	defer p.allocatedMu.Unlock()
	return len(p.allocated)
}

// Close will cleanup the memory pool and deallocate ALL previously allocated blocks.
// Using any of the blocks returned from Get() after a call to Close() will result in a panic
func (p *Pool) Close() error {
	return p.cleanup()
}

func (p *Pool) tryPop() ([]byte, bool) {
	p.freeMu.Lock()
	defer p.freeMu.Unlock()

	if len(p.free) == 0 {
		return nil, false
	}
	n := len(p.free) - 1
	result := p.free[n]
	p.free[n] = nil
	p.free = p.free[:n]
	return result, true
}

func (p *Pool) checkValidBlock(block []byte) error {
	if len(block) == 0 || len(block) != cap(block) {
		return ErrInvalidBlock
	}

	k := &block[0]
	p.allocatedMu.Lock()
	_, found := p.allocated[k]
	p.allocatedMu.Unlock()

	if !found || len(block) != p.initOpts.blockSize {
		return ErrInvalidBlock
	}
	return nil
}

func (p *Pool) prealloc(n int) error {
	if n < 0 || n > p.maxBlocks {
		return ErrPreallocOutOfBounds
	}

	for i := 0; i < n; i++ {
		block, err := alloc(p.initOpts.blockSize)
		if err != nil {
			_ = p.cleanup()
			return err
		}
		k := &block[0]
		p.allocated[k] = struct{}{}
		p.free = append(p.free, block)
	}
	return nil
}

func (p *Pool) cleanup() error {
	p.allocatedMu.Lock()
	multiErr := NewMultiErr()
	for arrayPtr := range p.allocated {
		var block []byte
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&block))
		hdr.Cap = p.initOpts.blockSize
		hdr.Len = p.initOpts.blockSize
		hdr.Data = uintptr(unsafe.Pointer(arrayPtr))
		if err := dealloc(block); err != nil {
			multiErr.Add(err)
		}
	}
	p.allocated = nil
	p.allocatedMu.Unlock()

	p.freeMu.Lock()
	p.free = nil
	p.freeMu.Unlock()
	return multiErr.Return()
}
