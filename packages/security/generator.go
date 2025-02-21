package security

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GeneratePassword(length int, includeNumbers, includeLetters, includeSpecial bool) (string, error) {
	const (
		numbers      = "0123456789"
		letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialChars = "&*{}'\"<>!$@[]"
	)

	if length > 50 {
		return "", fmt.Errorf("doesn't make any sense to have such a big password")
	}
	var charset string
	if includeNumbers {
		charset += numbers
	}
	if includeLetters {
		charset += letters
	}
	if includeSpecial {
		charset += specialChars
	}

	if len(charset) == 0 {
		return "", fmt.Errorf("at least one character set must be selected")
	}

	password := make([]byte, length)
	for i := range password {
		char, err := randomChar(charset)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	return string(password), nil
}

func randomChar(charset string) (byte, error) {
	max := big.NewInt(int64(len(charset)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}
