package server
import (
	"github.com/janmbaco/go-infrastructure/errors"
)

type ListenerError interface {
	errors.CustomError
	GetErrorType() ListenerErrorType
}

type listenerError struct {
	errors.CustomizableError
	ErrorType ListenerErrorType
}

func newListenerError(errorType ListenerErrorType, message string, internalError error) ListenerError {
	return &listenerError{
		CustomizableError: errors.CustomizableError{
			Message:       message,
			InternalError: internalError,
		},
		ErrorType: errorType,
	}
}

func (e *listenerError) GetErrorType() ListenerErrorType {
	return e.ErrorType
}

type ListenerErrorType uint8

const (
	UnexpectedError ListenerErrorType = iota
	AddressNotConfigured
)

