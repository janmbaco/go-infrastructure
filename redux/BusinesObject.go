package redux

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/event"
	"reflect"
	"strings"
)

type businessObject struct {
	actions       *actionsContainer
	stateManager  *stateManager
	reducer       Reducer
	subscriptions []func(state interface{})
}

func (bo *businessObject) onNewState() {
	newState := bo.stateManager.GetState()
	for _, fn := range bo.subscriptions {
		errorhandler.OnErrorContinue(func() { fn(newState) })
	}
}

func (bo *businessObject) subscribe(fn func(state interface{})) {
	bo.subscriptions = append(bo.subscriptions, fn)
}

type businessObjectBuilder struct {
	initialState interface{}
	actions      *actionsContainer
	blf          map[uintptr]reflect.Value //business logic funcionality
}

func NewBusinessObjectBuilder(initialState interface{}, actions interface{}) *businessObjectBuilder {

	if initialState == nil {
		panic("initialState parameter can't be nil")
	}

	if actions == nil {
		panic("actions parameter can't be nil")
	}

	return &businessObjectBuilder{initialState: initialState, actions: getActionsContainer(actions), blf: make(map[uintptr]reflect.Value)}
}

func (builder *businessObjectBuilder) On(action *Action, function interface{}) *businessObjectBuilder {
	idx := reflect.ValueOf(action).Pointer()
	if idx == 0 {
		panic("The action can`t be nil!")
	}

	actionName, ok := builder.actions.actionName[idx]
	if !ok {
		panic("This action doesn`t belong to this BusinesObject!")
	}

	if _, exists := builder.blf[idx]; exists {
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
		panic(fmt.Sprintf("The function for action `%v` must to have the contract func(state `%v`, payload *any) `%v`", builder.actions.actionName[idx].name, typeOfState.Name(), typeOfState.Name()))
	}

	actionName.action.typ = functionType.In(1)

	builder.blf[idx] = functionValue
	return builder
}

func (builder *businessObjectBuilder) GetBusinessObjecct() *businessObject {
	stateManager := &stateManager{
		publisher: event.NewEventPublisher(),
		state:     builder.initialState,
	}

	panicMessage := strings.Builder{}
	for idx, actionName := range builder.actions.actionName {
		if _, ok := builder.blf[idx]; !ok {
			panicMessage.WriteString(fmt.Sprintf("The logic for the actions '%v' is not defined!\n", actionName.name))
		}
	}

	if panicMessage.Len() > 0 {
		panic(panicMessage.String())
	}

	reducer := func(state interface{}, action *Action) interface{} {

		function, exists := builder.blf[reflect.ValueOf(action).Pointer()]
		if !exists {
			panic("The action is not located in the reducer function!")
		}

		return function.Call([]reflect.Value{
			reflect.ValueOf(state),
			action.getPayload(),
		})[0].Interface()
	}

	bo := &businessObject{
		actions:      builder.actions,
		stateManager: stateManager,
		reducer:      reducer,
	}
	bo.stateManager.Subscribe(bo.onNewState)
	return bo
}
