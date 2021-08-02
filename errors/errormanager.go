package errors

import (
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
	"sync"
)

type (
	ErrorSetter interface {
		On(err error, callback func(err error))
	}
	ErrorCallbacks interface {
		GetCallback(err error) func(err error)
	}
	ErrorManager interface {
		ErrorSetter
		ErrorCallbacks
	}
)
type errorManager struct {
	errorCallbacks sync.Map
}

func NewErrorManager() ErrorManager {
	return &errorManager{}
}

func (e *errorManager) GetCallback(err error) func(err error) {
	errorschecker.CheckNilParameter(map[string]interface{}{"err": err})
	if fn, ok := e.errorCallbacks.Load(reflect.Indirect(reflect.ValueOf(err)).Type()); ok {
		return fn.(reflect.Value).Interface().(func(err error))
	}
	return nil
}

func (e *errorManager) On(err error, callback func(err error)) {
	errorschecker.CheckNilParameter(map[string]interface{}{"err": err, "callback": callback})
	e.errorCallbacks.LoadOrStore(reflect.Indirect(reflect.ValueOf(err)).Type(), reflect.ValueOf(callback))
}
