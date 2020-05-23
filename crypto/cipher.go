package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/janmbaco/go-infrastructure/errorhandler"
)

type Crypter struct {
	aead cipher.AEAD
}

func NewCrypter(key []byte) *Crypter {
	block, err := aes.NewCipher(key)
	errorhandler.TryPanic(err)
	aead, err := cipher.NewGCM(block)
	errorhandler.TryPanic(err)
	return &Crypter{
		aead: aead,
	}
}

func (this *Crypter) Encrypt(value []byte) []byte {
	nonce := make([]byte, this.aead.NonceSize())
	return this.aead.Seal(nonce, nonce, value, nil)
}

func (this *Crypter) Decrypt(value []byte) []byte {
	nonceSize := this.aead.NonceSize()
	nonce, cipherValue := value[:nonceSize], value[nonceSize:]
	plainValue, err := this.aead.Open(nil, nonce, cipherValue, nil)
	errorhandler.TryPanic(err)
	return plainValue
}
