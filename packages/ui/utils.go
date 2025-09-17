package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/mas2020-golang/cryptex/packages/console"
)

// ClearScreen clears the terminal on macOS/Linux/Windows.
// It prefers ANSI sequences (fast, no external process).
// If ANSI isn't supported, it falls back to the OS command.
func ClearScreen() {
	// Try ANSI first (works on macOS/Linux; on Windows if VT mode is enabled).
	if tryANSI() {
		return
	}

	// Fallback to OS command.
	if runtime.GOOS == "windows" {
		_ = exec.Command("cmd", "/c", "cls").Run()
	} else {
		_ = exec.Command("clear").Run()
	}
}

func tryANSI() bool {
	if runtime.GOOS == "windows" {
		// Enable VT mode on Windows consoles (no-op elsewhere).
		console.EnableANSI()
	}
	// CSI 2J: clear screen; CSI H: cursor home
	_, err := fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")
	return err == nil
}
