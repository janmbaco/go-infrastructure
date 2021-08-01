package fileconfig

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"reflect"
)

// FileConfigHandlerError is the struct of an error occurs in FileConfigHandler object
type FileConfigHandlerError struct {
	errors.CustomError
	ErrorType FileConfigHandlerErrorType
}

func newFileConfigHandlerError(errorType FileConfigHandlerErrorType, message string) *FileConfigHandlerError {
	return &FileConfigHandlerError{
		ErrorType: errorType,
		CustomError: errors.CustomError{
			Message:       message,
			InternalError: nil,
		}}
}

type FileConfigHandlerErrorType uint8

const (
	UnexpectedError FileConfigHandlerErrorType = iota
	OldConfigNilError
)

type fileConfigHandleErrorPipe struct{}

func (*fileConfigHandleErrorPipe) Pipe(err error) error {
	resultError := err

	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); errType != reflect.TypeOf(&FileConfigHandlerError{}) {
		errorType := UnexpectedError
		resultError = &FileConfigHandlerError{
			CustomError: errors.CustomError{
				Message:       err.Error(),
				InternalError: err,
			},
			ErrorType: errorType,
		}
	}
	return resultError
}
