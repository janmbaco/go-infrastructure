package events

import (
	"reflect"
)

type RestoredEvent struct {
}

func (*RestoredEvent) IsParallelPropagation() bool {
	return true
}

func (*RestoredEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func() {})

}

func (*RestoredEvent) StopPropagation() bool {
	return false
}

func (*RestoredEvent) HasEventArgs() bool {
	return false
}

func (*RestoredEvent) GetEventArgs() interface{} {
	return nil
}
