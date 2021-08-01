package server

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"reflect"
)

type ListenerError struct {
	errors.CustomError
	ErrorType ListenerErrorType
}

func newListenerError(errorType ListenerErrorType, message string) *ListenerError {
	return &ListenerError{
		ErrorType: errorType,
		CustomError: errors.CustomError{
			Message:       message,
			InternalError: nil,
		}}
}

type ListenerErrorType uint8

const (
	UnexpectedError ListenerErrorType = iota
	AddressNotConfigured
)

type listenerErrorPipe struct{}

func (listenerErrorPipe *listenerErrorPipe) Pipe(err error) error {
	resultError := err

	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); errType != reflect.TypeOf(&ListenerError{}) {
		errorType := UnexpectedError
		resultError = &ListenerError{
			CustomError: errors.CustomError{
				Message:       err.Error(),
				InternalError: err,
			},
			ErrorType: errorType,
		}
	}
	return resultError
}
