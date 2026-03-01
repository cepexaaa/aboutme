//go:build darwin

package rawmode

import (
	"os"
	"syscall"
	"unsafe"
)

type termios struct {
	Iflag  uint64
	Oflag  uint64
	Cflag  uint64
	Lflag  uint64
	Cc     [20]uint8
	Ispeed uint64
	Ospeed uint64
}

func setRawMode() *termios {
	var old termios

	// macOS - TIOCGETA / TIOCSETA
	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(os.Stdin.Fd()),
		uintptr(syscall.TIOCGETA),
		uintptr(unsafe.Pointer(&old)),
	); err != 0 {
		return nil
	}

	newTermios := old
	newTermios.Lflag &^= syscall.ICANON | syscall.ECHO

	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(os.Stdin.Fd()),
		uintptr(syscall.TIOCSETA),
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
			uintptr(syscall.TIOCSETA),
			uintptr(unsafe.Pointer(old)),
		)
	}
}
