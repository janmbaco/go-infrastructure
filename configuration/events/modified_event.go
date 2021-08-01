package events

import (
	"reflect"
)

type ModifiedEvent struct {
}

func (*ModifiedEvent) IsParallelPropagation() bool {
	return true
}

func (*ModifiedEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func() {})

}

func (*ModifiedEvent) StopPropagation() bool {
	return false
}

func (*ModifiedEvent) HasEventArgs() bool {
	return false
}

func (*ModifiedEvent) GetEventArgs() interface{} {
	return nil
}
