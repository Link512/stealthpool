package stealthpool

import (
	"reflect"
	"runtime"
	"syscall"
	"unsafe"
)

const (
	memCommit  = 0x1000
	memReserve = 0x2000
	memRelease = 0x8000

	pageRW = 0x04

	kernelDll   = "kernel32.dll"
	allocFunc   = "VirtualAlloc"
	deallocFunc = "VirtualFree"

	errOK = 0
)

var (
	kernel         *syscall.DLL
	virtualAlloc   *syscall.Proc
	virtualDealloc *syscall.Proc
)

func init() {
	runtime.LockOSThread()
	kernel = syscall.MustLoadDLL(kernelDll)
	virtualAlloc = kernel.MustFindProc(allocFunc)
	virtualDealloc = kernel.MustFindProc(deallocFunc)
	runtime.UnlockOSThread()
}

func alloc(size int) ([]byte, error) {
	addr, _, err := virtualAlloc.Call(uintptr(0), uintptr(size), memCommit|memReserve, pageRW)
	errNo := err.(syscall.Errno)
	if errNo != errOK {
		return nil, err
	}

	var result []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	hdr.Data = addr
	hdr.Cap = size
	hdr.Len = size
	return result, nil
}

func dealloc(b []byte) error {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	_, _, err := virtualDealloc.Call(hdr.Data, 0, memRelease)
	errNo := err.(syscall.Errno)
	if errNo != errOK {
		return err
	}

	return nil
}
