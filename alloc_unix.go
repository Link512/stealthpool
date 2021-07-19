// +build !windows,!appengine

package stealthpool

import "golang.org/x/sys/unix"

func alloc(size int) ([]byte, error) {
	return unix.Mmap(
		-1,                                  // required by MAP_ANONYMOUS
		0,                                   // offset from file descriptor start, required by MAP_ANONYMOUS
		size,                                // how much memory
		unix.PROT_READ|unix.PROT_WRITE,      // protection on memory
		unix.MAP_ANONYMOUS|unix.MAP_PRIVATE, // private so other processes don't see the changes, anonymous so that nothing gets synced to the file
	)
}

func dealloc(b []byte) error {
	return unix.Munmap(b)
}
