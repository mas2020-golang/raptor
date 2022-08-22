package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

var (
	Version, GitCommit string
)

func init() {
	Version = "0.1.0"
}

// GetBytesFromPipe reads from the pipe and return the buffer of bytes of the given argument
func GetBytesFromPipe() *os.File {
	//var bs []byte
	//buf := bytes.NewBuffer(bs)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		//scanner := bufio.NewScanner(os.Stdin)
		//
		//for scanner.Scan() {
		//	buf.Write(scanner.Bytes())
		//	fmt.Print(scanner.Text())
		//}
		//
		//if err := scanner.Err(); err != nil {
		//	log.Fatal(err)
		//}
		return os.Stdin
	}
	//fmt.Printf("number of bytes from the pipe are %d\n", len(buf.Bytes()))
	return nil
}

// ReadPassword reads the standard input in hidden mode
func ReadPassword(text string) (string, error) {
	fmt.Print(text)
	buf, err := term.ReadPassword(int(os.Stdin.Fd()))
	return string(buf), err
}

// Check checks if an error and exit
func Check(err error, message string) {
	var errorMsg string
	if err != nil {
		if len(message) > 0 {
			errorMsg = fmt.Sprintf("%s caused by %v", message, err)
		} else {
			errorMsg = fmt.Sprintf("%v", err)
		}
		Error(errorMsg)
		os.Exit(1)
	}
}

// GetText returns a text read from a bufio.Reader interface object
func GetText(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	output := strings.Replace(text, "\n", "", -1)
	return strings.Replace(output, "\r", "", -1)
}

// GetTextWithEsc returns a text read from a bufio.Reader interface object.
// The delimiter is the char sequence >>
func GetTextWithEsc(reader *bufio.Reader) string {
	buf := bytes.Buffer{}
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "ERROR!"
		} else {
			buf.Write([]byte{b})
			if buf.Len() >= 2 {
				bytesBuf := buf.Bytes()
				if bytesBuf[len(bytesBuf)-1] == 62 &&
					bytesBuf[len(bytesBuf)-2] == 62 {
					return string(bytesBuf)[0 : len(bytesBuf)-2]
				}
			}
		}
	}
}

// askForPassword asks for a password once or twice. You can change
// the default requested text. Returns the key to use
func AskForPassword(text string, twice bool) (key string, err error) {
	// only for debugging
	if os.Getenv("CRYPTEX_DBGPWD") != "" {
		key = os.Getenv("CRYPTEX_DBGPWD")
	} else {
		// ask for password
		key, err = ReadPassword(text)
		if err != nil {
			return "", err
		}
		fmt.Println("")
		if twice {
			key2, err := ReadPassword("Repeat the pwd:")
			if err != nil {
				return "", err
			}
			fmt.Println("")
			if key != key2 {
				return "", fmt.Errorf("the passwords need to be the same")
			}
		}
		if len(key) < 6 {
			return "", fmt.Errorf("the password is too short, use at least a 6 chars length")
		}
	}

	return key, nil
}
