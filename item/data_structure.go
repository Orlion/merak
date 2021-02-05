package item

import "github.com/Orlion/merak/data_structure"

type ItemSet struct {
	elems map[*Item]struct{}
}

func NewItemSet() *ItemSet {
	return &ItemSet{
		elems: make(map[*Item]struct{}),
	}
}

func (s *ItemSet) Add(it *Item) {
	s.elems[it] = struct{}{}
}

func (s *ItemSet) AddList(its []*Item) {
	for _, it := range its {
		s.elems[it] = struct{}{}
	}
}

func (s *ItemSet) Exists(it *Item) bool {
	_, exists := s.elems[it]
	return exists
}

func (s *ItemSet) Delete(it *Item) {
	delete(s.elems, it)
}

func (s *ItemSet) Elems() map[*Item]struct{} {
	return s.elems
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
