package errors

import (
	"fmt"

	"github.com/janmbaco/go-infrastructure/logs"
)

// ErrorCatcher defines an object responsible to catch errors without panic/recover
type ErrorCatcher interface {
	HandleError(err error, errorfn func(error)) error
	HandleErrorWithFinally(err error, errorfn func(error), finallyfn func()) error
	TryCatchError(tryfn func() error, errorfn func(error)) error
	TryFinally(tryfn func() error, finallyfn func()) error
	TryCatchErrorAndFinally(tryfn func() error, errorfn func(error), finallyfn func()) error
	OnErrorContinue(tryfn func() error) error
	CatchError(err error, errorfn func(error)) error
	CatchErrorAndFinally(err error, errorfn func(error), finallyfn func()) error
}

type errorCatcher struct {
	logger logs.ErrorLogger
}

// NewErrorCatcher returns an ErrorCatcher object
func NewErrorCatcher(logger logs.ErrorLogger) ErrorCatcher {
	if logger == nil {
		return &errorCatcher{logger: &noOpLogger{}}
	}
	return &errorCatcher{logger: logger}
}

// HandleError handles error and executes error function if err != nil
func (e *errorCatcher) HandleError(err error, errorfn func(error)) error {
	if err != nil && errorfn != nil {
		errorfn(err)
	}
	return err
}

// HandleErrorWithFinally always executes finally, then handles error
func (e *errorCatcher) HandleErrorWithFinally(err error, errorfn func(error), finallyfn func()) error {
	if finallyfn != nil {
		finallyfn()
	}
	return e.HandleError(err, errorfn)
}

// TryCatchError executes function and catches returned error
func (e *errorCatcher) TryCatchError(tryfn func() error, errorfn func(error)) error {
	if tryfn == nil {
		return fmt.Errorf("tryfn cannot be nil")
	}

	err := tryfn()
	if err != nil && errorfn != nil {
		errorfn(err)
	}
	return err
}

// TryFinally always executes finally function
func (e *errorCatcher) TryFinally(tryfn func() error, finallyfn func()) error {
	if tryfn == nil {
		return fmt.Errorf("tryfn cannot be nil")
	}

	defer func() {
		if finallyfn != nil {
			finallyfn()
		}
	}()

	return tryfn()
}

// TryCatchErrorAndFinally always executes finally, then catches error
func (e *errorCatcher) TryCatchErrorAndFinally(tryfn func() error, errorfn func(error), finallyfn func()) error {
	if tryfn == nil {
		return fmt.Errorf("tryfn cannot be nil")
	}

	defer func() {
		if finallyfn != nil {
			finallyfn()
		}
	}()

	err := tryfn()
	if err != nil && errorfn != nil {
		errorfn(err)
	}
	return err
}

// OnErrorContinue executes function and logs error if occurs, but continues execution
func (e *errorCatcher) OnErrorContinue(tryfn func() error) error {
	if tryfn == nil {
		return fmt.Errorf("tryfn cannot be nil")
	}

	err := tryfn()
	if err != nil {
		e.logger.TryError(err)
	}
	return nil
}

// CatchError handles error if present
func (e *errorCatcher) CatchError(err error, errorfn func(error)) error {
	return e.HandleError(err, errorfn)
}

// CatchErrorAndFinally always executes finally, then handles error
func (e *errorCatcher) CatchErrorAndFinally(err error, errorfn func(error), finallyfn func()) error {
	return e.HandleErrorWithFinally(err, errorfn, finallyfn)
}

// noOpLogger is a fallback logger
type noOpLogger struct{}

func (n *noOpLogger) PrintError(level logs.LogLevel, err error)    {}
func (n *noOpLogger) TryPrintError(level logs.LogLevel, err error) {}
func (n *noOpLogger) TryTrace(err error)                           {}
func (n *noOpLogger) TryInfo(err error)                            {}
func (n *noOpLogger) TryWarning(err error)                         {}
func (n *noOpLogger) TryError(err error)                           {}
func (n *noOpLogger) TryFatal(err error)                           {}
