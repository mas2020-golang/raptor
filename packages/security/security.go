package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

// encryptBox encrypts the in []byte and return the encrypted
// out []byte or an error
func EncryptBox(in []byte, key string) ([]byte, error) {
	// create the block and the iv factor
	block, err := getCypher(key)
	if err != nil {
		return nil, fmt.Errorf("the Cyther block has not been created: %v", err)
	}
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("error creating the iv factor: %v", err)
	}

	// buffer size must be multiple of 16 bytes
	b := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	bufIn := bytes.NewBuffer(in)
	bufOut := &bytes.Buffer{}
	for {
		n, err := bufIn.Read(b)
		if n > 0 {
			stream.XORKeyStream(b, b[:n])
			// Write into file
			_, err = bufOut.Write(b[:n])
			if err != nil {
				return nil, fmt.Errorf("error writing to the output buffer: %v", err)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("read %d bytes: %v", n, err)
		}
	}
	// Append the IV
	if _, err = bufOut.Write(iv); err != nil {
		return nil, fmt.Errorf("error writing the iv factor to the buffer: %v", err)
	}
	return bufOut.Bytes(), nil
}

// decryptBox decrypts the in []byte and return the encrypted
// out []byte or an error
func DecryptBox(in []byte, key string) ([]byte, error) {
	var msgLen int64
	// create the block and the iv factor
	block, err := getCypher(key)
	if err != nil {
		return nil, fmt.Errorf("the Cyther block has not been created: %v", err)
	}
	iv := make([]byte, block.BlockSize())

	// get the iv factor from the input file
	iv = in[len(in)-len(iv) : len(in)]

	// buffer size must be multiple of 16 bytes
	b := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	bufIn := bytes.NewBuffer(in)
	bufOut := &bytes.Buffer{}
	msgLen = int64(len(in)) - int64(len(iv))
	for {
		n, err := bufIn.Read(b)
		if n > 0 {
			// for decryption only
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read %d bytes: %v", n, err)
		}
		// cypher the bytes
		stream.XORKeyStream(b, b[:n])
		// Write into the buffer
		_, err = bufOut.Write(b[:n])
		if err != nil {
			return nil, fmt.Errorf("error writing to the output buffer: %v", err)
		}
	}

	return bufOut.Bytes(), nil
}

// getCypher return the Cipher
func getCypher(key string) (cipher.Block, error) {
	k := sha256.Sum256([]byte(key))
	return aes.NewCipher(k[:])
}
