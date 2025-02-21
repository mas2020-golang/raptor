package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"log/slog"

	wipe "github.com/0x9ef/go-wiper/wipe"
	"github.com/mas2020-golang/goutils/output"
)

var ErrInvalidFile = errors.New("invalid file type for encryption or decryption")

// encryptBox encrypts the in []byte and return the encrypted
// out []byte or an error
func EncryptBox(in []byte, key string) ([]byte, error) {
	return encrypt(in, key)
}

func DecryptBox(in []byte, key string) ([]byte, error) {
	return decrypt(in, key)
}

// getCypher return the Cipher
func getCypher(key string) (cipher.Block, error) {
	k := sha256.Sum256([]byte(key))
	return aes.NewCipher(k[:])
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	// Generate a 256-bit key from the passphrase
	key := sha256.Sum256([]byte(passphrase))

	// Create a new AES cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create a GCM cipher mode instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Generate a nonce with the required length
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Encrypt the data using AES-GCM
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	// Generate a 256-bit key from the passphrase
	key := sha256.Sum256([]byte(passphrase))

	// Create a new AES cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create a GCM cipher mode instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Ensure the ciphertext length is greater than the nonce size
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	// Extract the nonce and actual ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt the data using AES-GCM
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}

	return plaintext, nil
}

func EncryptFile(path, passphrase string) error {
	// ends with .enc
	if strings.HasSuffix(path, ".enc") {
		output.Warning("", fmt.Sprintf("file %s skipped as it is already a .enc file", path))
		return ErrInvalidFile
	}

	// Read the file contents
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Encrypt the data
	encryptedData, err := encrypt(data, passphrase)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}
	slog.Debug(fmt.Sprintf("the file %s has been encrypted", path))

	// Write the encrypted data to a new file with .enc extension
	encryptedFilePath := path + ".enc"
	err = ioutil.WriteFile(encryptedFilePath, encryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write encrypted file: %v", err)
	}

	// delete the file
	return deleteFile(path)
}

func DecryptFile(path, passphrase string) error {
	// ends with .enc
	if !strings.HasSuffix(path, ".enc") {
		output.Warning("", fmt.Sprintf("file %s skipped as it is not a .enc file", path))
		return ErrInvalidFile
	}

	// Read the encrypted file contents
	encryptedData, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %v", err)
	}

	// Decrypt the data
	decryptedData, err := decrypt(encryptedData, passphrase)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %v", err)
	}
	slog.Debug(fmt.Sprintf("the file %s has been decrypted", path))

	// Write the decrypted data to a new file without the .enc extension
	decryptedFilePath := strings.TrimSuffix(path, ".enc")
	err = ioutil.WriteFile(decryptedFilePath, decryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted file: %v", err)
	}

	// delete the .enc file
	return deleteFile(path)
}

func EncryptDirectory(dirPath, passphrase string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		err = EncryptFile(path, passphrase)
		if errors.Is(err, ErrInvalidFile) {
			return nil
		}
		return err
	})
}

func DecryptDirectory(dirPath, passphrase string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		err = DecryptFile(path, passphrase)
		if errors.Is(err, ErrInvalidFile) {
			return nil
		}
		return err
	})
}

// deleteFile securely deletes the path
func deleteFile(path string) error {
	slog.Debug("security.deleteFile()", "path", path)
	policy := &wipe.Policy{"UsDod5220_22_M", "US Department of Defense DoD 5220.22-M (3 passes)", wipe.RuleUsDod5220_22_M}
	err := wipe.Wipe(path, policy.Rule)
	if err != nil {
		return err
	}

	return os.Remove(path)
}
