package redux

import (
	"fmt"
	"reflect"
)

type Action interface {
	With(interface{}) Action
	GetPayload() reflect.Value
	SetType(reflect.Type)
	GetName() string
}

type ActionsObject interface {
	GetActions() []Action
	GetActionsNames() []string
	Contains(action Action) bool
	ContainsByName(actionName string) bool
	GetActionByName(name string) Action
	GetNameByAction(action Action) string
}

type action struct {
	payload   reflect.Value
	typ       reflect.Type
	payloaded bool
	name      string
}

func (action *action) With(payload interface{}) Action {

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

func (action *action) GetPayload() reflect.Value {
	var result reflect.Value
	if action.payloaded {
		result = action.payload
		action.payloaded = false
	} else {
		result = reflect.Zero(action.typ)
	}
	return result
}

func (action *action) SetType(typ reflect.Type) {
	action.typ = typ
}

func (action *action) GetName() string {
	return action.name
}

type actionObject struct {
	actions      []Action
	actionsNames []string
	actionByName map[string]Action
	nameByAction map[Action]string
}

func NewActionObject(object interface{}) *actionObject {
	result := &actionObject{
		actions:      getActionsIn(object),
		actionsNames: make([]string, 0),
		actionByName: make(map[string]Action),
		nameByAction: make(map[Action]string),
	}
	for _, b := range result.actions {
		result.actionsNames = append(result.actionsNames, b.GetName())
		result.actionByName[b.GetName()] = b
		result.nameByAction[b] = b.GetName()
	}
	return result
}

func (actionObject *actionObject) GetActions() []Action {
	return actionObject.actions
}

func (actionObject *actionObject) Contains(action Action) bool {
	_, ok := actionObject.nameByAction[action]
	return ok
}

func (actionObject *actionObject) GetActionsNames() []string {
	return actionObject.actionsNames
}

func (actionObject *actionObject) ContainsByName(actionName string) bool {
	_, ok := actionObject.actionByName[actionName]
	return ok
}

func (actionObject *actionObject) GetActionByName(name string) Action {
	return actionObject.actionByName[name]
}

func (actionObject *actionObject) GetNameByAction(action Action) string {
	return actionObject.nameByAction[action]
}

func getActionsIn(object interface{}) []Action {
	if object == nil {
		panic("The object parameter can`t be nil!")
	}
	result := make([]Action, 0)
	rv := reflect.Indirect(reflect.ValueOf(object))
	rt := rv.Type()
	actionType := reflect.TypeOf((*Action)(nil)).Elem()
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Type.Implements(actionType) {
			if rv.Field(i).IsNil() {
				rv.Field(i).Set(reflect.ValueOf(&action{name: rt.Field(i).Name}))
			}
			result = append(result, rv.Field(i).Elem().Interface().(Action))
		}
	}
	if len(result) == 0 {
		panic("There isn`t any action on the actionsObject object!")
	}
	return result
}
