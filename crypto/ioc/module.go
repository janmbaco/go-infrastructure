package ioc
import (
	"github.com/janmbaco/go-infrastructure/crypto"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
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
		map[uint]string{0: "key"},
	)

	return nil
}
