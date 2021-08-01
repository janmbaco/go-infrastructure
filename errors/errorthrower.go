package errors

// ErrorThrower defines a object responsible to throws errors
type ErrorThrower interface {
	Throw(err error)
}

type errorThrower struct {
	errorCallbacks ErrorCallbacks
}

// NewErrorThrower returns a ErrorThrower
func NewErrorThrower(errorCallbacks ErrorCallbacks) ErrorThrower {
	return &errorThrower{errorCallbacks: errorCallbacks}
}

// Throw throws a error to a callback or, if callbacks is nil, to panic
func (e *errorThrower) Throw(err error) {
	if fn := e.errorCallbacks.GetCallback(err); fn != nil {
		fn(err)
	} else {
		panic(err)
	}
}
