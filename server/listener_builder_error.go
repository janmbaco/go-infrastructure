package server

import (
	"github.com/janmbaco/go-infrastructure/errors"
)

// ListenerBuilderError is the errors of ListenerBuilder
type ListenerBuilderError interface {
	errors.CustomError
	GetErrorType() ListenerBuilderErrorType
}

type listenerBuilderError struct {
	errors.CustomizableError
	ErrorType ListenerBuilderErrorType
}

func newListenerBuilderError(errorType ListenerBuilderErrorType, message string, internalError error) ListenerBuilderError {
	return &listenerBuilderError{
		CustomizableError: errors.CustomizableError{
			Message:       message,
			InternalError: internalError,
		},
		ErrorType: errorType,
	}
}

func (e *listenerBuilderError) GetErrorType() ListenerBuilderErrorType {
	return e.ErrorType
}

type ListenerBuilderErrorType uint8

const (
	UnexpectedBuilderError ListenerBuilderErrorType = iota
	NilBootstraperError
	NilGrpcDefinitionsError
)

