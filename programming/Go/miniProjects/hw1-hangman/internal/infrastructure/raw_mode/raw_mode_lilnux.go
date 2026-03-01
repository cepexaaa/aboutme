//go:build linux

package rawmode

import (
	"os"
	"syscall"
	"unsafe"
)

type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

func setRawMode() *termios {
	var old termios

	// current settings
	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(os.Stdin.Fd()),
		uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(&old)),
	); err != 0 {
		return nil
	}

	newTermios := old
	newTermios.Lflag &^= syscall.ICANON | syscall.ECHO

	// new settings
	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(os.Stdin.Fd()),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&newTermios)),
	); err != 0 {
		return nil
	}

	return &old
}

func restoreMode(old *termios) {
	if old != nil {
		syscall.Syscall(
			syscall.SYS_IOCTL,
			uintptr(os.Stdin.Fd()),
			uintptr(syscall.TCSETS),
			uintptr(unsafe.Pointer(old)),
		)
	}
}
