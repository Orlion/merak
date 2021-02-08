package item

import (
	"github.com/Orlion/merak/data_structure"
)

type Set struct {
	elems []*Item // TODO: redis intset
}

func NewSet() *Set {
	return &Set{
		elems: make([]*Item, 0),
	}
}

func (set *Set) Add(it *Item) bool {
	if !set.Exists(it) {
		set.elems = append(set.elems, it)
		return true
	}

	return false
}

func (set *Set) AddList(its []*Item) (addNum int) {
	for _, it := range its {
		if set.Add(it) {
			addNum++
		}
	}

	return
}

func (set *Set) Exists(it *Item) bool {
	return set.index(it) > -1
}

func (set *Set) index(it *Item) int {
	i := -1

	for k, elem := range set.elems {
		if elem != nil && elem.Equals(it) {
			i = k
			break
		}
	}

	return i
}

func (set *Set) Delete(it *Item) bool {
	if i := set.index(it); i > -1 {
		set.elems[i] = nil
		return true
	}

	return false
}

func (set *Set) Elems() []*Item {
	return set.elems
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
