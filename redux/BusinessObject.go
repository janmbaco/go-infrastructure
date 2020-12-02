package redux

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/event"
	"github.com/janmbaco/go-infrastructure/logs"
	"reflect"
	"strings"
)

type BusinessObject struct {
	ActionsObject ActionsObject
	Reducer       Reducer
	StateManager  StateManager
}

type businessObjectBuilder struct {
	initialState  interface{}
	stateManager  StateManager
	actionsObject ActionsObject
	blf           map[Action]reflect.Value //business logic funcionality
}

func NewBusinessObjectBuilder(initialState interface{}) *businessObjectBuilder {

	if initialState == nil {
		panic("initialState parameter can't be nil")
	}

	return &businessObjectBuilder{
		initialState: initialState,
		stateManager: &stateManager{
			publisher: event.NewEventPublisher(),
			state:     initialState,
		},
		blf: make(map[Action]reflect.Value)}
}

func (builder *businessObjectBuilder) SetActionsObject(object ActionsObject) *businessObjectBuilder {
	if object == nil {
		panic("stateManager parameter can't be nil")
	}

	builder.actionsObject = object

	return builder
}

func (builder *businessObjectBuilder) SetStateManage(stateManager StateManager) *businessObjectBuilder {
	if stateManager == nil {
		panic("stateManager parameter can't be nil")
	}
	builder.stateManager = stateManager
	return builder
}

func (builder *businessObjectBuilder) On(action Action, function interface{}) *businessObjectBuilder {
	if reflect.ValueOf(action).Pointer() == 0 {
		panic("The action can`t be nil!")
	}

	if !builder.actionsObject.Contains(action) {
		panic("This action doesn`t belong to this BusinesObject!")
	}

	if _, exists := builder.blf[action]; exists {
		panic("action already reduced!")
	}

	functionValue := reflect.ValueOf(function)
	if functionValue.Pointer() == 0 {
		panic("The function can`t be nil!")
	}

	functionType := reflect.TypeOf(function)
	if functionType.Kind() != reflect.Func {
		panic("The function must be a Func!")
	}

	if typeOfState := reflect.TypeOf(builder.initialState); functionType.NumIn() != 2 || functionType.NumOut() != 1 || functionType.In(0) != functionType.Out(0) || functionType.In(0) != typeOfState {
		panic(fmt.Sprintf("The function for action `%v` must to have the contract func(state `%v`, payload *any) `%v`", builder.actionsObject.GetNameByAction(action), typeOfState.Name(), typeOfState.Name()))
	}

	action.SetType(functionType.In(1))

	builder.blf[action] = functionValue
	return builder
}

func (builder *businessObjectBuilder) SetActionsLogicByObject(object interface{}) *businessObjectBuilder {
	if object == nil {
		panic("The object can`t be nil")
	}
	rv := reflect.ValueOf(object)
	rt := reflect.TypeOf(object)
	if rt.Kind() != reflect.Ptr && rt.Kind() != reflect.Struct {
		panic("The object must be a struct")
	}

	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mt := m.Type
		if typeOfState := reflect.TypeOf(builder.initialState); mt.NumIn() == 3 && mt.NumOut() == 1 && mt.In(1) == mt.Out(0) && mt.In(1) == typeOfState {
			if builder.actionsObject.ContainsByName(m.Name) {
				action := builder.actionsObject.GetActionByName(m.Name)
				action.SetType(mt.In(2))
				builder.blf[action] = rv.Method(i)
			} else {
				logs.Log.Warning(fmt.Sprintf("The func`%v` in the object `%v` has not a action asociated in the ActionsObject! ActionObject:`%v`", m.Name, rt.String(), builder.actionsObject.GetActionsNames()))
			}

		}
	}
	return builder
}

func (builder *businessObjectBuilder) GetBusinessObject() *BusinessObject {

	if builder.actionsObject == nil {
		panic("There isn´t any ActionsObject to load to the BusinessObject!")
	}

	panicMessage := strings.Builder{}
	for _, action := range builder.actionsObject.GetActions() {
		if _, ok := builder.blf[action]; !ok {
			panicMessage.WriteString(fmt.Sprintf("The logic for the actionsObject '%v' is not defined!\n", builder.actionsObject.GetNameByAction(action)))
		}
	}

	if panicMessage.Len() > 0 {
		panic(panicMessage.String())
	}

	reducerFunc := func(state interface{}, action Action) interface{} {

		function, exists := builder.blf[action]
		if !exists {
			panic("The action is not located in the reducer function!")
		}

		return function.Call([]reflect.Value{
			reflect.ValueOf(state),
			action.GetPayload(),
		})[0].Interface()
	}

	reducer := &reducer{
		reducer: reducerFunc,
	}

	return &BusinessObject{
		ActionsObject: builder.actionsObject,
		Reducer:       reducer,
		StateManager:  builder.stateManager,
	}
}
