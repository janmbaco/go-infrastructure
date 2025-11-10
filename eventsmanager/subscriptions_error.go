package eventsmanager

import (
	"github.com/janmbaco/go-infrastructure/v2/errors"
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
