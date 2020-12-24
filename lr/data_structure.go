package lr

import (
	"container/list"
	"github.com/Orlion/merak/symbol"
)

type ProductionZSet struct {
	list []*Production
	keys map[string]struct{}
}

func newProductionZSet() *ProductionZSet {
	return &ProductionZSet{
		keys: make(map[string]struct{}),
	}
}

func (s *ProductionZSet) add(production *Production) {
	code := production.getCode()
	if _, exists := s.keys[code]; !exists {
		s.list = append(s.list, production)
		s.keys[code] = struct{}{}
	}
}

func (s *ProductionZSet) addAll(productions []*Production) {
	for _, production := range productions {
		code := production.getCode()
		if _, exists := s.keys[code]; !exists {
			s.list = append(s.list, production)
			s.keys[code] = struct{}{}
		}
	}
}

func (s *ProductionZSet) exists(production *Production) bool {
	_, exists := s.keys[production.getCode()]
	return exists
}

type ProductionStack struct {
	stack *Stack
}

func newProductionStack() *ProductionStack {
	return &ProductionStack{
		stack: NewStack(),
	}
}

func (s *ProductionStack) Push(p *Production) {
	s.stack.Push(p)
}

func (s *ProductionStack) Pop() *Production {
	return s.stack.Pop().(*Production)
}

func (s *ProductionStack) Top() *Production {
	return s.stack.Top().(*Production)
}

func (s *ProductionStack) Empty() bool {
	return s.stack.Empty()
}

type Stack struct {
	list *list.List
}

func NewStack() *Stack {
	return &Stack{
		list: list.New(),
	}
}

func (s *Stack) Push(e interface{}) {
	s.list.PushBack(e)
}

func (s *Stack) Pop() interface{} {
	v := s.list.Back()
	s.list.Remove(v)
	return v.Value
}

func (s *Stack) Top() interface{} {
	v := s.list.Back()
	return v.Value
}

func (s *Stack) Empty() bool {
	return s.list.Len() == 0
}


type Promise struct {
	set *symbol.ZSet
	symbol symbol.Symbol
	isPromise bool
}

func NewPromise(set *symbol.ZSet, symbol symbol.Symbol, isPromise bool) *Promise {
	return &Promise{set, symbol, isPromise}
}