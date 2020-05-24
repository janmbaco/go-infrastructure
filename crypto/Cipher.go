package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/janmbaco/go-infrastructure/errorhandler"
)

type Cipher struct {
	aead cipher.AEAD
}

func NewCipher(key []byte) *Cipher {
	block, err := aes.NewCipher(key)
	errorhandler.TryPanic(err)
	aead, err := cipher.NewGCM(block)
	errorhandler.TryPanic(err)
	return &Cipher{
		aead: aead,
	}
}

func (this *Cipher) Encrypt(value []byte) []byte {
	nonce := make([]byte, this.aead.NonceSize())
	return this.aead.Seal(nonce, nonce, value, nil)
}

func (this *Cipher) Decrypt(value []byte) []byte {
	nonceSize := this.aead.NonceSize()
	nonce, cipherValue := value[:nonceSize], value[nonceSize:]
	plainValue, err := this.aead.Open(nil, nonce, cipherValue, nil)
	errorhandler.TryPanic(err)
	return plainValue
}
