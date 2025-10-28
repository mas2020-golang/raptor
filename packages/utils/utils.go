package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mas2020-golang/cryptex/packages/security"
	"github.com/mas2020-golang/goutils/output"
	"golang.org/x/term"
	"gopkg.in/yaml.v2"
)

var (
	Version, GitCommit, BuildDate string
	BufferBox                     *Box
	BoxPath, BoxPwd               string
)

func init() {
	Version = "0.5.0-SNAPSHOT"
}

type Secret struct {
	Name        string            `yaml:"name,omitempty"`
	Id          int32             `yaml:"id,omitempty"` // Unique ID number for this secret
	Pwd         string            `yaml:"pwd,omitempty"`
	Url         string            `yaml:"url,omitempty"`
	Notes       string            `yaml:"notes,omitempty"`
	Others      map[string]string `yaml:"others,omitempty"`
	Version     string            `yaml:"version,omitempty"`
	Login       string            `yaml:"login,omitempty"`
	LastUpdated string            `yaml:"lastUpdated,omitempty"`
}

type Box struct {
	Name        string    `yaml:"name,omitempty"`
	Version     string    `yaml:"version,omitempty"`
	LastUpdated string    `yaml:"lastUpdated,omitempty"`
	Owner       string    `yaml:"owner,omitempty"`
	Secrets     []*Secret `yaml:"secrets,omitempty"`
	Size        int64     `yaml:"-"`
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

func GetComplexText() (string, error) {
	var note strings.Builder
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "EOF" {
			break
		}
		note.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}

	return note.String(), nil
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
				return "", fmt.Errorf("the passwords do not correspond")
			}
		}
		if len(key) < 6 {
			return "", fmt.Errorf("the password is too short, use at least a 6 chars length")
		}
	}

	return key, nil
}

// GetFolderBox returns the default folder where boxes are stored.
// Precedence:
// 1. CRYPTEX_FOLDER env var
// 2. OS-specific config directory + "raptor/boxes"
// 3. LOCALAPPDATA (Windows only) + "raptor/boxes"
// 4. HOME + ".raptor/boxes"
// 5. Relative "boxes" folder
func getFolderBox() string {
	if v := os.Getenv("CRYPTEX_FOLDER"); v != "" {
		return v
	}

	// 2. OS-aware config dir
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "raptor", "boxes")
	}

	// 3. Explicit Windows fallback
	if runtime.GOOS == "windows" {
		if ld := os.Getenv("LOCALAPPDATA"); ld != "" {
			return filepath.Join(ld, "raptor", "boxes")
		}
	}

	// 4. Home fallback
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".raptor", "boxes")
	}

	// 5. Last-resort relative path
	return "boxes"
}

// InitFolderBox ensures the box folder exists on the target system.
// It uses GetFolderBox() to resolve the path, then creates the directory
// (and its parents) if they don't exist. It returns the absolute path.
func InitFolderBox() (string, error) {
	path := getFolderBox()

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for box folder: %w", err)
	}

	// MkdirAll is idempotent: does nothing if the folder already exists.
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return "", fmt.Errorf("failed to create box folder %q: %w", abs, err)
	}

	return abs, nil
}

// OpenBox opens a box
func OpenBox(boxName, pwd string) (string, string, *Box, error) {
	// if the box is in the buffer you can get into it
	if BufferBox != nil {
		return BoxPath, BoxPwd, BufferBox, nil
	}

	// check if the boxName is a file, in that case BoxPath is overrided by that
	if validPath, _ := IsValidFilePath(boxName); validPath {
		BoxPath = boxName
	}

	// search the CRYPTEX_BOX env if name is empty
	if len(boxName) == 0 {
		boxName = os.Getenv("CRYPTEX_BOX")
		if len(boxName) == 0 {
			return "", "", nil, fmt.Errorf("--box args is not given and the env var CRYPTEX_BOX is empty")
		}
	}

	// get the folder box
	boxFolder, err := InitFolderBox()
	if err != nil {
		return "", "", nil, err
	}

	// read the box (if it is not assigned yet)
	if len(BoxPath) == 0 {
		BoxPath = path.Join(boxFolder, boxName)
	}

	in, err := ioutil.ReadFile(BoxPath)
	if err != nil {
		return "", "", nil, fmt.Errorf("reading the file box in %s: %v", BoxPath, err)
	}

	if len(pwd) == 0 {
		// ask for the password
		pwd, err = AskForPassword("Password: ", false)
		if err != nil {
			return "", "", nil, err
		}
	}

	// encrypt the box
	decIn, err := security.DecryptBox(in, pwd)
	if err != nil {
		return "", "", nil, fmt.Errorf("decrypting the file box in %s: %v", BoxPath, err)
	}

	box := &Box{}
	err = yaml.Unmarshal(decIn, box)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to read the box: %v. Maybe an incorrect pwd?", err)
	}
	BoxPwd = pwd
	return BoxPath, pwd, box, nil
}

func SaveBox(path, key string, box *Box) error {
	out, err := yaml.Marshal(box)
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

func IsValidFilePath(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Path does not exist
			return false, nil
		}
		// An error occurred while trying to access the path
		return false, err
	}
	// Check if the path is a file (not a directory)
	return !info.IsDir(), nil
}
