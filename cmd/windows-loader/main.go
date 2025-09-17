//go:build windows
package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	exe, err := os.Executable()
	if err != nil {
		// Best effort: just try raptor-real.exe in CWD
		_ = chain("raptor-real.exe", os.Args[1:]...)
		return
	}
	dir := filepath.Dir(exe)
	target := filepath.Join(dir, "raptor-real.exe")
	_ = chain(target, os.Args[1:]...)
}

func chain(target string, args ...string) error {
	cmd := exec.Command(target, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// No extra environment tweaks needed; inherits env and console
	return cmd.Run()
}