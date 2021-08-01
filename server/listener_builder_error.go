package server

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"reflect"
)

type ListenerBuilderError struct {
	errors.CustomError
	ErrorType ListenerBuilderErrorType
}

func newListenerBuilderError(errorType ListenerBuilderErrorType, message string) *ListenerBuilderError {
	return &ListenerBuilderError{
		ErrorType: errorType,
		CustomError: errors.CustomError{
			Message:       message,
			InternalError: nil,
		}}
}

type ListenerBuilderErrorType uint8

const (
	UnexpectedBuilderError ListenerBuilderErrorType = iota
	NilBootstraperError
	NilGrpcDefinitionsError
)

type listenerBuilderErrorPipe struct{}

func (listenerErrorPipe *listenerBuilderErrorPipe) Pipe(err error) error {
	resultError := err
	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); errType != reflect.TypeOf(&ListenerBuilderError{}) {
		resultError = &ListenerBuilderError{
			CustomError: errors.CustomError{
				Message:       err.Error(),
				InternalError: err,
			},
			ErrorType: UnexpectedBuilderError,
		}
	}
	return resultError
}
