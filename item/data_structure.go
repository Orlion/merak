package item

import "github.com/Orlion/merak/data_structure"

type ItemZSet struct {
	list []*Item
	keys map[int]struct{}
}

func NewItemZSet() *ItemZSet {
	return &ItemZSet{
		keys: make(map[int]struct{}),
	}
}

func (s *ItemZSet) Add(it *Item) {
	id := it.Id()
	if _, exists := s.keys[id]; !exists {
		s.list = append(s.list, it)
		s.keys[id] = struct{}{}
	}
}

func (s *ItemZSet) AddList(its []*Item) {
	for _, it := range its {
		id := it.Id()
		if _, exists := s.keys[id]; !exists {
			s.list = append(s.list, it)
			s.keys[id] = struct{}{}
		}
	}
}

func (s *ItemZSet) Exists(it *Item) bool {
	_, exists := s.keys[it.Id()]
	return exists
}

func (s *ItemZSet) List() []*Item {
	return s.list
}

type ItemStack struct {
	stack *data_structure.Stack
}

func NewItemStack() *ItemStack {
	return &ItemStack{
		stack: data_structure.NewStack(),
	}
}

func (s *ItemStack) Push(it *Item) {
	s.stack.Push(it)
}

func (s *ItemStack) Pop() *Item {
	return s.stack.Pop().(*Item)
}

func (s *ItemStack) Top() *Item {
	return s.stack.Top().(*Item)
}

func (s *ItemStack) Empty() bool {
	return s.stack.Empty()
}
