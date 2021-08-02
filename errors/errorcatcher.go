package errors

import (
	"errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/logs"
)

// ErrorCatcher defines a object responsible to catch the errors
type ErrorCatcher interface {
	CatchError(err error, errorfn func(error))
	CatchErrorAndFinally(err error, errorfn func(error), finallyfn func())
	OnErrorContinue(tryfn func())
	EvenErrorFinally(err error, finallyfn func())
	TryCatchError(tryfn func(), errorfn func(error))
	TryFinally(tryfn func(), finallyfn func())
	TryCatchErrorAndFinally(tryfn func(), errorfn func(error), finallyfn func())
}

type errorCatcher struct {
	logger logs.ErrorLogger
}

// NewErrorCatcher returns a ErrorCatcher object
func NewErrorCatcher(logger logs.ErrorLogger) ErrorCatcher {
	return &errorCatcher{logger: logger}
}

// OnErrorContinue continues even if an error occurs
func (e *errorCatcher) OnErrorContinue(tryfn func()) {
	defer e.deferTryError(true, nil, nil)
	tryfn()
}

// TryCatchError catches the error and executes the error function
func (e *errorCatcher) TryCatchError(tryfn func(), errorfn func(error)) {
	defer e.deferTryError(false, errorfn, nil)
	tryfn()
}

// TryFinally execute allways the finally function
func (e *errorCatcher) TryFinally(tryfn func(), finallyfn func()) {
	defer e.deferTryError(false, nil, finallyfn)
	tryfn()
}

// TryCatchErrorAndFinally execute allways the finally function and then catches the error and executes the error function
func (e *errorCatcher) TryCatchErrorAndFinally(tryfn func(), errorfn func(error), finallyfn func()) {
	defer e.deferTryError(false, errorfn, finallyfn)
	tryfn()
}

// CatchError if err is different from nil executes the error function
func (e *errorCatcher) CatchError(err error, errorfn func(error)) {
	defer e.deferTryError(false, errorfn, nil)
	errorschecker.TryPanic(err)
}

// EvenErrorFinally always execute the function finally, even if the error is different from nil
func (e *errorCatcher) EvenErrorFinally(err error, finallyfn func()) {
	defer e.deferTryError(false, nil, finallyfn)
	errorschecker.TryPanic(err)
}

// CatchErrorAndFinally execute allways the finally function and then if the error is different from nil  executes the error function
func (e *errorCatcher) CatchErrorAndFinally(err error, errorfn func(error), finallyfn func()) {
	defer e.deferTryError(false, errorfn, finallyfn)
	errorschecker.TryPanic(err)
}

func (e *errorCatcher) deferTryError(shouldContinue bool, errorfn func(error), finallyfn func()) {
	if finallyfn != nil {
		finallyfn()
	}
	if re := recover(); re != nil {
		err := errors.New("unexpected error")
		switch re.(type) {
		case string:
			err = errors.New(re.(string))
		case error:
			err = re.(error)
		}

		if errorfn != nil {
			errorfn(err)
		} else if !shouldContinue {
			e.logger.PrintError(logs.Error, err)
			panic(err)
		}

	}
}
