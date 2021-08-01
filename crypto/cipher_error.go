package crypto

import "github.com/janmbaco/go-infrastructure/errors"

// CipherError is the definition of errors that can occur using a Cipher object
type CipherError struct {
	errors.CustomError
}

type cipherErrorPipe struct{}

func (*cipherErrorPipe) Pipe(err error) error {
	return &CipherError{CustomError: errors.CustomError{Message: err.Error(), InternalError: err}}
}
