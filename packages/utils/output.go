package utils

import (
	"fmt"
	"runtime"

	"github.com/mas2020-golang/goutils/output"
)

// TODO: remove all the code that can be replaced using goutils
const (
	Reset    = "\033[0m"
	Bold     = "\033[1m"
	Blue     = "\033[34m"
	Orange   = "\033[38;5;167m"
	Green    = "\033[32m"
	Red      = "\033[31m"
	LightRed = "\033[91m"
	Yellow   = "\033[33m"
)

// Warning returns a warning string
func Success(text string) {
	os := runtime.GOOS
	switch os {
	case "windows":
		fmt.Printf("%s%s%s\n", output.GreenS("DONE: "), text, Reset)
	case "darwin":
		fmt.Printf("%s%s%s\n", output.GreenS("👍 "), text, Reset)
	case "linux":
		fmt.Printf("%s%s%s\n", output.GreenS("👍 "), text, Reset)
	default:
		fmt.Printf("%s%s%s\n", output.GreenS("✔ "), text, Reset)
	}
}

// Note returns a note formatted text
func Note(text string) {
	os := runtime.GOOS
	switch os {
	case "windows":
		fmt.Printf("%s%s%s\n", output.GreenS("> "), text, Reset)
	case "darwin":
		fmt.Printf("%s%s%s\n", "👉 ", text, Reset)
	case "linux":
		fmt.Printf("%s%s%s\n", "✔ ", text, Reset)
	default:
		fmt.Printf("%s%s%s\n", output.GreenS("✔ "), text, Reset)
	}
}

func Verbosity(msg string, verbose bool){
	if verbose{
		output.Activity(msg)
	}
}
