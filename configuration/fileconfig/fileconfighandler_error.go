package fileconfig

import (
	"github.com/janmbaco/go-infrastructure/errors"
)

// FileConfigHandlerError is the struct of an error occurs in FileConfigHandler object
type FileConfigHandlerError interface {
	errors.CustomError
	GetErrorType() FileConfigHandlerErrorType
}

type fileConfigHandlerError struct {
	errors.CustomizableError
	ErrorType FileConfigHandlerErrorType
}

func newFileConfigHandlerError(errorType FileConfigHandlerErrorType, message string, internalError error) FileConfigHandlerError {
	return &fileConfigHandlerError{
		CustomizableError: errors.CustomizableError{
			Message:       message,
			InternalError: internalError,
		},
		ErrorType: errorType,
	}
}

func (e *fileConfigHandlerError) GetErrorType() FileConfigHandlerErrorType {
	return e.ErrorType
}

// FileConfigHandlerErrorType is the type of the errors of FileConfigHandler
type FileConfigHandlerErrorType uint8

const (
	UnexpectedError FileConfigHandlerErrorType = iota
	OldConfigNilError
)