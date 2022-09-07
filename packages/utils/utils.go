package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/security"
	"github.com/mas2020-golang/goutils/output"
	"golang.org/x/term"
	"google.golang.org/protobuf/proto"
)

var (
	Version, GitCommit string
)

func init() {
	Version = "0.2.0-dev"
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
		output.Error("", errorMsg)
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
func AskForPassword(text string, twice bool, minLength int8) (key string, err error) {
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
		if len(key) < int(minLength) {
			return "", fmt.Errorf("the password is too short, use at least a %d chars length", minLength)
		}
	}

	return key, nil
}

// getFolderBox returns the box folder
func GetFolderBox() (string, error) {
	// check the folder .cryptex
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// read the file in the home dir
	if len(os.Getenv("CRYPTEX_FOLDER")) > 0 {
		return os.Getenv("CRYPTEX_FOLDER"), nil
	} else {
		// read the file in the home dir
		return path.Join(home, ".cryptex", "boxes"), nil
	}
}

// OpenBox opens a box
func OpenBox(boxName string) (string, string, *protos.Box, error) {
	var boxPath string
	// search the CRYPTEX_BOX env if name is empty
	if len(boxName) == 0 {
		boxName = os.Getenv("CRYPTEX_BOX")
		if len(boxName) == 0 {
			return "", "", nil, fmt.Errorf("--box args is not given and the env var CRYPTEX_BOX is empty")
		}
	}
	// get the folder box
	boxFolder, err := GetFolderBox()
	if err != nil {
		return "", "", nil, fmt.Errorf("problem to determine th folder box: %v", err)
	}

	// read the box
	boxPath = path.Join(boxFolder, boxName)
	in, err := ioutil.ReadFile(boxPath)
	if err != nil {
		return "", "", nil, fmt.Errorf("reading the file box in %s: %v", boxPath, err)
	}

	// ask for the password
	key, err := AskForPassword("Box password: ", false, 0)
	if err != nil {
		return "", "", nil, err
	}
	// encrypt the box
	decIn, err := security.DecryptBox(in, key)
	if err != nil {
		return "", "", nil, fmt.Errorf("decrypting the file box in %s: %v", boxPath, err)
	}

	box := &protos.Box{}
	err = proto.Unmarshal(decIn, box)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to read the box: %v. Maybe an incorrect pwd?", err)
	}
	return boxPath, key, box, nil
}

func SaveBox(path, key string, box *protos.Box) error {
	out, err := proto.Marshal(box)
	if err != nil {
		return fmt.Errorf("failed to encode the box: %v", err)
	}
	// encrypt the box
	encOut, err := security.EncryptBox(out, key)
	if err := ioutil.WriteFile(path, encOut, 0644); err != nil {
		return fmt.Errorf("failed to encrypt the box: %v", err)
	}
	if err := ioutil.WriteFile(path, encOut, 0644); err != nil {
		return fmt.Errorf("failed to write the box: %v", err)
	}
	return nil
}
