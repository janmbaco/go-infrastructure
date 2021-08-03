package errors

import (
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
	"sync"
)

type (
	// ErrorCallbacks is the definition of a object responsible to get the callbacks registereb by a ErrorManager object
	ErrorCallbacks interface {
		GetCallback(err error) func(err error)
	}
	// ErrorManager is the definition of a object responsible to set callbacks to error definitions
	ErrorManager interface {
		On(err error, callback func(err error))
	}
)
type errorManager struct {
	errorCallbacks sync.Map
}

// NewErrorManager returns a ErrorManager
func NewErrorManager() ErrorManager {
	return &errorManager{}
}

// GetCallback gets the callback from a error
func (e *errorManager) GetCallback(err error) func(err error) {
	errorschecker.CheckNilParameter(map[string]interface{}{"err": err})
	if fn, ok := e.errorCallbacks.Load(reflect.Indirect(reflect.ValueOf(err)).Type()); ok {
		return fn.(reflect.Value).Interface().(func(err error))
	}
	return nil
}

// On register a callback to an error
func (e *errorManager) On(err error, callback func(err error)) {
	errorschecker.CheckNilParameter(map[string]interface{}{"err": err, "callback": callback})
	e.errorCallbacks.LoadOrStore(reflect.Indirect(reflect.ValueOf(err)).Type(), reflect.ValueOf(callback))
}
