package disk

import "reflect"

type fileChangedEvent struct{}

func (*fileChangedEvent) IsParallelPropagation() bool {
	return true
}

func (*fileChangedEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func() {})
}

func (*fileChangedEvent) StopPropagation() bool {
	return false
}

func (*fileChangedEvent) GetEventArgs() interface{} {
	return nil
}

func (*fileChangedEvent) HasEventArgs() bool {
	return false
}
