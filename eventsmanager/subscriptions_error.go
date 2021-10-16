package eventsmanager

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"reflect"
)

type SubscriptionsError interface {
	errors.CustomError
	GetErrorType() SubscriptionsErrorType
}

type subscriptionsError struct {
	errors.CustomizableError
	ErrorType SubscriptionsErrorType
}

func newSubscriptionsError(errorType SubscriptionsErrorType, message string, internalError error) SubscriptionsError {
	return &subscriptionsError{
		CustomizableError: errors.CustomizableError{
			Message:       message,
			InternalError: internalError,
		},
		ErrorType: errorType,
	}
}
func (e *subscriptionsError) GetErrorType() SubscriptionsErrorType {
	return e.ErrorType
}

type SubscriptionsErrorType uint8

const (
	Unexpected SubscriptionsErrorType = iota
	BadFunctionSignature
	FunctionNoSubscribed
)

type subscriptionsErrorPipe struct{}

func (*subscriptionsErrorPipe) Pipe(err error) error {
	resultError := err
	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); !errType.Implements(reflect.TypeOf((*SubscriptionsError)(nil)).Elem()) {
		resultError = newSubscriptionsError(Unexpected, err.Error(), err)
	}
	return resultError
}
