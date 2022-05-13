package utils

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
)

var (
	Version, GitCommit string
)

func init() {
	Version = "0.1.0-dev"
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
	buf, err := term.ReadPassword(0)
	return string(buf), err
}

// Check checks if an error and exit
func Check(err error, message string) {
	var errorMsg string
	if err != nil {
		if len(message) > 0{
			errorMsg = fmt.Sprintf("%s caused by %v", message, err)
		}else{
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

// askForPassword asks for a password once or twice. Returns the key to use
func AskForPassword(twice bool) (key string, err error) {
	// ask for password
	key, err = ReadPassword("Password: ")
	if err != nil {
		return "", err
	}
	fmt.Println("")
	if twice {
		key2, err := ReadPassword("Repeat the password:")
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
	return key, nil
}