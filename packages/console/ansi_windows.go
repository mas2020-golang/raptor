//go:build windows
package console

import (
	"golang.org/x/sys/windows"
)

func EnableANSI() {
	h := windows.Handle(windows.Stdout)
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err == nil {
		_ = windows.SetConsoleMode(h, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}
