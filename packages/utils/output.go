package utils

import (
	"fmt"
	"runtime"
)

const (
	Reset    = "\033[0m"
	Bold     = "\033[1m"
	Blue     = "\033[34m"
	Orange   = "\033[38;5;167m"
	Green    = "\033[32m"
	Red = "\033[31m"
	LightRed = "\033[91m"
	Yellow   = "\033[33m"
)

// BoldS returns a string bold
func BoldS(t string) string {
	return fmt.Sprintf("%s%s%s", Bold, t, Reset)
}

// BoldOut bolds the passed argument
func BoldOut(t string) {
	fmt.Printf("%s%s%s", Bold, t, Reset)
}

// BlueS returns a blue string
func BlueS(t string) string {
	return fmt.Sprintf("%s%s%s", Blue, t, Reset)
}

// BlueOut write a blue text into the given writer
func BlueOut(t string) {
	fmt.Printf("%s%s%s", Blue, t, Reset)
}

// GreenS returns a green string
func GreenS(t string) string {
	return fmt.Sprintf("%s%s%s", Green, t, Reset)
}

// GreenOut write a green text into the given writer
func GreenOut(t string) {
	fmt.Printf("%s%s%s", Green, t, Reset)
}

// RedS returns a red string
func RedS(t string) string {
	return fmt.Sprintf("%s%s%s", Red, t, Reset)
}

// RedS returns a light red string
func LightRedS(t string) string {
	return fmt.Sprintf("%s%s%s", LightRed, t, Reset)
}

// RedOut write a red text into the given writer
func RedOut(t string) {
	fmt.Printf("%s%s%s", LightRed, t, Reset)
}

// OrangeS returns an orange string
func OrangeS(t string) string {
	return fmt.Sprintf("%s%s%s", Orange, t, Reset)
}

// YellowS returns a yellow string
func YellowS(t string) string {
	return fmt.Sprintf("%s%s%s", Yellow, t, Reset)
}

// Warning returns a warning string
func Success(text string) {
	os := runtime.GOOS
	switch os {
	case "windows":
		fmt.Printf("%s%s%s\n", GreenS("DONE: "), text, Reset)
	case "darwin":
		fmt.Printf("%s%s%s\n", GreenS("👍 "), text, Reset)
	case "linux":
		fmt.Printf("%s%s%s\n", GreenS("✔ "), text, Reset)
	default:
		fmt.Printf("%s%s%s\n", GreenS("✔ "), text, Reset)
	}
}

// Warning returns a warning string
func Warning(text string) {
	fmt.Printf("%s%s%s\n", YellowS("Warning: "), text, Reset)
}

// Warning returns a warning string
func Error(text string) {
	fmt.Printf("%s%s%s\n", RedS("Error: "), text, Reset)
}
