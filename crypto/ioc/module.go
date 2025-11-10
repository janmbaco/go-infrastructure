package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/crypto"
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
)

// CryptoModule implements Module for cryptography services
type CryptoModule struct{}

// NewCryptoModule creates a new crypto module
func NewCryptoModule() *CryptoModule {
	return &CryptoModule{}
}

// RegisterServices registers all cryptography services
func (m *CryptoModule) RegisterServices(register dependencyinjection.Register) error {
	dependencyinjection.RegisterSingletonWithParams[crypto.Cipher](
		register,
		crypto.NewCipher,
		map[int]string{0: "key"},
	)

	return nil
}
