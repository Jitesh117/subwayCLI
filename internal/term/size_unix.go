//go:build darwin || linux

package term

import (
	"os"
	"syscall"
	"unsafe"
)

func platformSize() (int, int, error) {
	ws := &struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}{}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return 0, 0, errno
	}

	return int(ws.Col), int(ws.Row), nil
}
