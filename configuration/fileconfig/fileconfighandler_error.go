package fileconfig

import (
	"github.com/janmbaco/go-infrastructure/errors"
)

// HandlerError is the struct of an error occurs in FileConfigHandler object
type HandlerError interface {
	errors.CustomError
	GetErrorType() HandlerErrorType
}

type fileConfigHandlerError struct {
	errors.CustomizableError
	ErrorType HandlerErrorType
}

func newFileConfigHandlerError(errorType HandlerErrorType, message string, internalError error) HandlerError {
	return &fileConfigHandlerError{
		CustomizableError: errors.CustomizableError{
			Message:       message,
			InternalError: internalError,
		},
		ErrorType: errorType,
	}
}

func (e *fileConfigHandlerError) GetErrorType() HandlerErrorType {
	return e.ErrorType
}

// HandlerErrorType is the type of the errors of FileConfigHandler
type HandlerErrorType uint8

const (
	UnexpectedError HandlerErrorType = iota
	OldConfigNilError
)
