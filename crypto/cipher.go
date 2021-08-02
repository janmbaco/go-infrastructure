package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/logs"
)

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
func NewCipher(key []byte, logger logs.Logger, thrower errors.ErrorThrower) Cipher {
	errorschecker.CheckNilParameter(map[string]interface{}{"thrower": thrower})
	block, err := aes.NewCipher(key)
	errorCatcher := errors.NewErrorCatcher(logger)
	errorschecker.TryPanic(err)
	aead, err := cipher.NewGCM(block)
	errorschecker.TryPanic(err)
	return &cipherImp{
		aead:         aead,
		errorCatcher: errorCatcher,
		errorHandler: errors.NewErrorDefer(thrower, &cipherErrorPipe{}),
	}
}

// Encrypt cipher the value
func (c *cipherImp) Encrypt(value []byte) []byte {
	defer c.errorHandler.TryThrowError()
	nonce := make([]byte, c.aead.NonceSize())
	return c.aead.Seal(nonce, nonce, value, nil)
}

// Decrypt deciphers the value
func (c *cipherImp) Decrypt(value []byte) []byte {
	defer c.errorHandler.TryThrowError()
	nonceSize := c.aead.NonceSize()
	nonce, cipherValue := value[:nonceSize], value[nonceSize:]
	plainValue, err := c.aead.Open(nil, nonce, cipherValue, nil)
	errorschecker.TryPanic(err)
	return plainValue
}
