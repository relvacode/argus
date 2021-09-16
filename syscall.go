//go:build windows
// +build windows

package argus

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
)

var (
	procOpenFileMappingW = modkernel32.NewProc("OpenFileMappingW")
)

func getUintptrFromBool(b bool) uintptr {
	if b {
		return 1
	} else {
		return 0
	}
}

func syscallOpenFileMappingW(dwDesiredAccess uint32, bInheritHandle bool, lpName *uint16) (handle windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procOpenFileMappingW.Addr(), 3, uintptr(dwDesiredAccess), getUintptrFromBool(bInheritHandle), uintptr(unsafe.Pointer(lpName)))
	handle = windows.Handle(r0)
	if handle == 0 {
		err = e1
	}

	return
}
