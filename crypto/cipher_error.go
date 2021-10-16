package crypto

import "github.com/janmbaco/go-infrastructure/errors"

// CipherError is the definition of errors that can occur using a Cipher object
type CipherError interface {
	errors.CustomError
}

type cipherError struct {
	errors.CustomizableError
}

type cipherErrorPipe struct{}

func (*cipherErrorPipe) Pipe(err error) error {
	return &cipherError{CustomizableError: errors.CustomizableError{Message: err.Error(), InternalError: err}}
}
