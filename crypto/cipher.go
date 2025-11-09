package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

// Cipher defines an object responsible to cipher and deciphers values by a key
type Cipher interface {
	Encrypt(value []byte) ([]byte, error)
	Decrypt(value []byte) ([]byte, error)
}

type cipherImp struct {
	aead cipher.AEAD
}

// NewCipher returns a Cipher object
func NewCipher(key []byte) (Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &cipherImp{
		aead: aead,
	}, nil
}

// Encrypt cipher the value
func (c *cipherImp) Encrypt(value []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	result := c.aead.Seal(nonce, nonce, value, nil)
	return result, nil
}

// Decrypt deciphers the value
func (c *cipherImp) Decrypt(value []byte) ([]byte, error) {
	nonceSize := c.aead.NonceSize()
	nonce, cipherValue := value[:nonceSize], value[nonceSize:]
	plainValue, err := c.aead.Open(nil, nonce, cipherValue, nil)
	if err != nil {
		return nil, err
	}
	return plainValue, nil
}
