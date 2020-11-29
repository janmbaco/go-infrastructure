package redux_test

import (
	"github.com/janmbaco/go-infrastructure/redux"
	"sync"
	"testing"
)

type Actions struct {
	Sumar *redux.Action
}

func TestNewStore(t *testing.T) {

	actions := &Actions{}

	businessObjectBuilder := redux.NewBusinessObjectBuilder(0, actions)

	businessObjectBuilder.On(actions.Sumar, func(state int, payload *int) int {
		var result int
		if payload == nil {
			result = state + 1
		} else {
			result = state + *payload
		}
		return result
	})

	store := redux.NewStore(businessObjectBuilder.GetBusinessObjecct())

	wg := sync.WaitGroup{}
	pass := 1
	store.Subscribe(actions, func(newState interface{}) {
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
	wg.Wait()

}
