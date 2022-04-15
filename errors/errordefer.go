package errors

import (
	"errors"
)

// ErrorDefer is the definition of a object responsible to hanles errors
type ErrorDefer interface {
	TryThrowError(fnPipeError func(err error) error)
}

type errorDefer struct {
	thrower   ErrorThrower
}

// NewErrorDefer return a ErrorDefer
func NewErrorDefer(thrower ErrorThrower) ErrorDefer {
	return &errorDefer{thrower: thrower}
}

// TryThrowError throws an error, through the handler, to panic, if error is different from nil
func (e *errorDefer) TryThrowError(fnPipeError func(err error) error)  {
	if re := recover(); re != nil {
		err := errors.New("unexpected error")
		switch re.(type) {
		case string:
			err = errors.New(re.(string))
		case error:
			err = re.(error)
		}
		if (fnPipeError != nil){
			err = fnPipeError(err)
		} 
		e.thrower.Throw(err)
	}
}
