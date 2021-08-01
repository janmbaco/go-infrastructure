package eventsmanager

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"reflect"
)

type SubscriptionsError struct {
	errors.CustomError
	ErrorType SubscriptionsErrorType
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
	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); errType != reflect.TypeOf(&SubscriptionsError{}) {
		resultError = &SubscriptionsError{
			CustomError: errors.CustomError{
				Message:       err.Error(),
				InternalError: err,
			},
			ErrorType: Unexpected,
		}
	}
	return resultError
}
