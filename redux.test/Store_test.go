package redux_test

import (
	"github.com/janmbaco/go-infrastructure/redux"
	"sync"
	"testing"
)

type Actions struct {
	actionsObject redux.ActionsObject
	Sumar         redux.Action
	Restar        redux.Action
}

func Sumar(state int, payload *int) int {
	var result int
	if payload == nil {
		result = state + 1
	} else {
		result = state + *payload
	}
	return result
}

type RestarObject struct{}

func (ro *RestarObject) Restar(state int, payload *int) int {
	return state - *payload
}

func TestNewStore(t *testing.T) {

	var actions = &Actions{}

	actions.actionsObject = redux.NewActionObject(actions)

	businessObjectBuilder := redux.NewBusinessObjectBuilder(0)

	businessObjectBuilder.SetActionsObject(actions.actionsObject)

	businessObjectBuilder.On(actions.Sumar, Sumar)

	businessObjectBuilder.SetActionsLogicByObject(&RestarObject{})

	store := redux.NewStore(businessObjectBuilder.GetBusinessObject())

	wg := sync.WaitGroup{}
	pass := 1
	store.Subscribe(actions.actionsObject, func(newState interface{}) {
		t.Log(newState)
		var expected int
		switch pass {
		case 1:
			expected = 1
		case 2:
			expected = 6
		case 3:
			expected = 7
		case 4:
			expected = 0
		case 5:
			expected = 7
		}
		if newState.(int) != expected {
			t.Errorf("expected: `%v`, found: `%v`", expected, newState.(int))
		}
		pass++
		wg.Done()
	})
	wg.Add(1)
	store.Dispatch(actions.Sumar)
	wg.Add(1)
	a := 5
	store.Dispatch(actions.Sumar.With(&a))
	wg.Add(1)
	store.Dispatch(actions.Sumar)
	wg.Add(1)

	a = -7
	store.Dispatch(actions.Sumar.With(&a))
	wg.Add(1)
	store.Dispatch(actions.Restar.With(&a))
	wg.Wait()

}
