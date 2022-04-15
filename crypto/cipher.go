package crypto

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
)

// Cipher defines an object responsible to cipher and deciphers values by a key
type Cipher interface {
	Encrypt(value []byte) []byte
	Decrypt(value []byte) []byte
}

type cipherImp struct {
	aead         cipher.AEAD
	errorCatcher errors.ErrorCatcher
	errorHandler errors.ErrorDefer
}

// NewCipher returns a Cipher object
func NewCipher(key []byte, errorCatcher errors.ErrorCatcher, errorDefer errors.ErrorDefer,) Cipher {
	errorschecker.CheckNilParameter(map[string]interface{}{"errorCatcher": errorCatcher, "errorDefer": errorDefer})
	block, err := aes.NewCipher(key)
	errorschecker.TryPanic(err)
	aead, err := cipher.NewGCM(block)
	errorschecker.TryPanic(err)
	return &cipherImp{
		aead:         aead,
		errorCatcher: errorCatcher,
		errorHandler: errorDefer,
	}
}

// Encrypt cipher the value
func (c *cipherImp) Encrypt(value []byte) []byte {
	defer c.errorHandler.TryThrowError(c.pipeError)
	nonce := make([]byte, c.aead.NonceSize())
	return c.aead.Seal(nonce, nonce, value, nil)
}

// Decrypt deciphers the value
func (c *cipherImp) Decrypt(value []byte) []byte {
	defer c.errorHandler.TryThrowError(c.pipeError)
	nonceSize := c.aead.NonceSize()
	nonce, cipherValue := value[:nonceSize], value[nonceSize:]
	plainValue, err := c.aead.Open(nil, nonce, cipherValue, nil)
	errorschecker.TryPanic(err)
	return plainValue
}

func (c *cipherImp) pipeError(err error) error {
	return &cipherError{CustomizableError: errors.CustomizableError{Message: err.Error(), InternalError: err}}
}
