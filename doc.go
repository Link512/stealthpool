// Package stealthpool provides a memory pool that allocates blocks off-heap that will NOT be tracked by the garbage collector.
// The name stealthpool comes from the fact that the memory being allocated by the pool is `stealthy` and will not be garbage collected ever
//
// These blocks should be used in situations where you want to keep garbage collection to a bare minimum.
// Needless to say, since the GC will not track any of the memory, extreme care must be had in order to avoid memory leaks:
//
//	pool, _ := stealthpool.New(1)
//	// ...
//	defer pool.Close() // always call Close to avoid memory leaks
//
// For now, the library works only on unix type OSes that support the `mmap` sycall
package stealthpool
