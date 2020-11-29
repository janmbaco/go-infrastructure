package redux

import (
	"fmt"
	"reflect"
)

type Action struct {
	payload   reflect.Value
	typ       reflect.Type
	payloaded bool
}

func (action *Action) With(payload interface{}) *Action {

	if action.typ.Kind() == reflect.Invalid {
		panic("This action is not asociated to any reduce!")
	}

	if reflect.TypeOf(payload) != action.typ {
		panic(fmt.Sprintf("The type of payload must be '%v'", action.typ.String()))
	}

	action.payload = reflect.ValueOf(payload)
	action.payloaded = true
	return action
}

func (action *Action) getPayload() reflect.Value {
	var result reflect.Value
	if action.payloaded {
		result = action.payload
		action.payloaded = false
	} else {
		result = reflect.Zero(reflect.TypeOf(reflect.New(action.typ).Elem().Interface()))
	}
	return result
}

type actionName struct {
	name   string
	action *Action
}

type actionsContainer struct {
	idx        uintptr
	actionName map[uintptr]*actionName
}

func getActionsContainer(v interface{}) *actionsContainer {

	result := &actionsContainer{
		idx:        reflect.ValueOf(v).Pointer(),
		actionName: make(map[uintptr]*actionName),
	}

	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		if rv.Field(i).Type() == reflect.TypeOf(&Action{}) {
			if rv.Field(i).Pointer() == 0 {
				rv.Field(i).Set(reflect.ValueOf(&Action{}))
			}
			result.actionName[rv.Field(i).Pointer()] = &actionName{
				name:   rt.Field(i).Name,
				action: rv.Field(i).Interface().(*Action),
			}
		}
	}
	if len(result.actionName) == 0 {
		panic("There isn`t any action on the actions object!")
	}

	return result
}

func (a *actionsContainer) contains(action *Action) bool {
	_, ok := a.actionName[reflect.ValueOf(action).Pointer()]
	return ok
}
