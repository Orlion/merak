package container

import "container/list"

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
