package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

// encryptBox encrypts the in []byte and return the encrypted
// out []byte or an error
func EncryptBox(in []byte, key string) ([]byte, error) {
	return encrypt(in, key)
}

// decryptBox decrypts the in []byte and return the encrypted
// out []byte or an error
// func DecryptBox(in []byte, key string) ([]byte, error) {
// 	var msgLen int64
// 	// create the block and the iv factor
// 	block, err := getCypher(key)
// 	if err != nil {
// 		return nil, fmt.Errorf("the Cyther block has not been created: %v", err)
// 	}
// 	iv := make([]byte, block.BlockSize())

// 	// get the iv factor from the input file
// 	iv = in[len(in)-len(iv) : len(in)]

// 	// buffer size must be multiple of 16 bytes
// 	b := make([]byte, 1024)
// 	stream := cipher.NewCTR(block, iv)
// 	bufIn := bytes.NewBuffer(in)
// 	bufOut := &bytes.Buffer{}
// 	msgLen = int64(len(in)) - int64(len(iv))
// 	for {
// 		n, err := bufIn.Read(b)
// 		if n > 0 {
// 			// for decryption only
// 			if n > int(msgLen) {
// 				n = int(msgLen)
// 			}
// 			msgLen -= int64(n)
// 		}
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return nil, fmt.Errorf("read %d bytes: %v", n, err)
// 		}
// 		// cypher the bytes
// 		stream.XORKeyStream(b, b[:n])
// 		// Write into the buffer
// 		_, err = bufOut.Write(b[:n])
// 		if err != nil {
// 			return nil, fmt.Errorf("error writing to the output buffer: %v", err)
// 		}
// 	}

// 	return bufOut.Bytes(), nil
// }

func DecryptBox(in []byte, key string) ([]byte, error) {
	return decrypt(in , key)
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
