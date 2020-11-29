package redux

import "reflect"

type store struct {
	businessObjects map[uintptr]*businessObject
}

func NewStore(bos ...*businessObject) *store {
	newStore := &store{
		businessObjects: make(map[uintptr]*businessObject),
	}

	for _, bo := range bos {
		if _, ko := newStore.businessObjects[bo.actions.idx]; ko {
			panic("Cannot add multiple BusinessObjects with the same ActionsContainer!")
		}
		newStore.businessObjects[bo.actions.idx] = bo
	}

	return newStore
}

func (s *store) Dispatch(action *Action) {
	for _, bo := range s.businessObjects {
		if bo.actions.contains(action) {
			bo.stateManager.SetState(bo.reducer(bo.stateManager.GetState(), action))
			break
		}
	}
}

func (s *store) Subscribe(actions interface{}, fn func(newState interface{})) {

	idx := reflect.ValueOf(actions).Pointer()
	if _, ok := s.businessObjects[idx]; !ok {
		panic("There is no BusinessObject for that ActionsContainer!")
	}

	s.businessObjects[idx].subscriptions = append(s.businessObjects[idx].subscriptions, fn)
}
