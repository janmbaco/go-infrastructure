package errors

import (
	"reflect"
	"sync"

	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
)

type (
	// ErrorCallbacks is the definition of a object responsible to get the callbacks registereb by a ErrorManager object
	ErrorCallbacks interface {
		GetCallback(err error) func(err error)
	}
	// ErrorManager is the definition of a object responsible to set callbacks to error definitions
	ErrorManager interface {
		On(err interface{}, callback func(err error))
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
	var fn func(err error)
	e.errorCallbacks.Range(func(key, value interface{}) bool {
		t :=reflect.TypeOf(err)
		if t.Implements(key.(reflect.Type)){
			fn = value.(reflect.Value).Interface().(func(err error))
			return false
		}
		return true
	})
	return fn
}

// On register a callback to an error
func (e *errorManager) On(err interface{}, callback func(err error)) {
	errorschecker.CheckNilParameter(map[string]interface{}{"err": err, "callback": callback})
	t := reflect.Indirect(reflect.ValueOf(err)).Type()
	if !t.Implements(reflect.TypeOf((*error)(nil)).Elem()){
		panic("The parameter 'err' not implement's error interface")
	}
	e.errorCallbacks.LoadOrStore(t, reflect.ValueOf(callback))
}
