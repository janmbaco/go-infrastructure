package errors

import (
	"errors"
)

// ErrorPipe is the definition of a object responsible to pipe an error to another
type ErrorPipe interface {
	Pipe(err error) error
}

// ErrorDefer is the definition of a object responsible to hanles errors
type ErrorDefer interface {
	TryThrowError()
}

type errorDefer struct {
	thrower   ErrorThrower
	errorPipe ErrorPipe
}

// NewErrorDefer return a ErrorDefer
func NewErrorDefer(thrower ErrorThrower, errorPipe ErrorPipe) ErrorDefer {
	return &errorDefer{thrower: thrower, errorPipe: errorPipe}
}

// TryThrowError throws an error, through the handler, to panic, if error is different from nil
func (e *errorDefer) TryThrowError() {
	if re := recover(); re != nil {
		err := errors.New("unexpected error")
		switch re.(type) {
		case string:
			err = errors.New(re.(string))
		case error:
			err = re.(error)
		}
		panic(e.errorPipe.Pipe(err))

	}
}
