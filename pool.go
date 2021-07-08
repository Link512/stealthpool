package stealthpool

import (
	"reflect"
	"sync"
	"unsafe"

	"go.uber.org/multierr"
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

type PoolOpt func(*poolOpts)

func WithPreAlloc(prealloc int) PoolOpt {
	return func(opts *poolOpts) {
		opts.preAlloc = prealloc
	}
}

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

func (p *Pool) Return(b []byte) error {
	if err := p.checkValidBlock(b); err != nil {
		return err
	}
	p.freeMu.Lock()
	defer p.freeMu.Unlock()
	p.free = append(p.free, b)
	return nil
}

func (p *Pool) FreeCount() int {
	p.freeMu.Lock()
	defer p.freeMu.Unlock()
	return len(p.free)
}

func (p *Pool) AllocCount() int {
	p.allocatedMu.Lock()
	defer p.allocatedMu.Unlock()
	return len(p.allocated)
}

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
	if n < 0 || n >= p.maxBlocks {
		return ErrPreallocOutOfBounds
	}

	for i := 0; i < n; i++ {
		block, err := alloc(p.initOpts.blockSize)
		if err != nil {
			p.cleanup()
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
	var err error
	for arrayPtr := range p.allocated {
		var block []byte
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&block))
		hdr.Cap = p.initOpts.blockSize
		hdr.Len = p.initOpts.blockSize
		hdr.Data = uintptr(unsafe.Pointer(arrayPtr))
		if dealocErr := dealloc(block); dealocErr != nil {
			err = multierr.Append(err, dealocErr)
		}
	}
	p.allocated = nil
	p.allocatedMu.Unlock()

	p.freeMu.Lock()
	p.free = nil
	p.freeMu.Unlock()
	return err
}
