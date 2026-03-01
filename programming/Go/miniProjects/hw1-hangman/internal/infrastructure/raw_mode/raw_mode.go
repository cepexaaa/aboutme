package rawmode

import (
	"fmt"
)

type RawMode struct {
	terminal *termios
}

func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

func MoveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func (rm *RawMode) RunRawMode() {
	rm.terminal = setRawMode()
	fmt.Print("\033[?1049h")
	hideCursor()
}

func (rm *RawMode) StopRawMode() {
	defer restoreMode(rm.terminal)
	defer showCursor()
	defer fmt.Print("\033[?1049l")
}
