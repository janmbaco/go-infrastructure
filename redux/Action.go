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
}

type ActionsContainer interface {
	GetActionsObject() ActionsObject
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
	actions []Action
}

func NewActionObject(object interface{}) *actionObject {
	return &actionObject{getActionsIn(object)}
}

func (actionObject *actionObject) GetActions() []Action {
	return actionObject.actions
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
		panic("There isn`t any action on the actionsContainer object!")
	}
	return result
}

type actionsContainer struct {
	actionsObject ActionsObject
	actionsNames  []string
	actionByName  map[string]Action
	nameByAction  map[Action]string
}

func NewActionsContainer(actionsObject ActionsObject) ActionsContainer {

	result := &actionsContainer{
		actionsObject: actionsObject,
		actionsNames:  make([]string, 0),
		actionByName:  make(map[string]Action),
		nameByAction:  make(map[Action]string),
	}
	for _, b := range actionsObject.GetActions() {
		result.actionsNames = append(result.actionsNames, b.GetName())
		result.actionByName[b.GetName()] = b
		result.nameByAction[b] = b.GetName()
	}
	return result
}

func (aContainer *actionsContainer) GetActionsObject() ActionsObject {
	return aContainer.actionsObject
}

func (aContainer *actionsContainer) Contains(action Action) bool {
	_, ok := aContainer.nameByAction[action]
	return ok
}

func (aContainer *actionsContainer) GetActionsNames() []string {
	return aContainer.actionsNames
}

func (aContainer *actionsContainer) ContainsByName(actionName string) bool {
	_, ok := aContainer.actionByName[actionName]
	return ok
}

func (aContainer *actionsContainer) GetActionByName(name string) Action {
	return aContainer.actionByName[name]
}

func (aContainer *actionsContainer) GetNameByAction(action Action) string {
	return aContainer.nameByAction[action]
}
