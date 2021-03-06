package item

import (
	"github.com/Orlion/merak/container"
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
		if elem.Equals(it) {
			i = k
			break
		}
	}

	return i
}

func (set *Set) Delete(it *Item) bool {
	if i := set.index(it); i > -1 {
		set.elems = append(set.elems[0:i], set.elems[i+1:]...)
		return true
	}

	return false
}

func (set *Set) DeleteByIndex(i int) {
	set.elems = append(set.elems[0:i], set.elems[i+1:]...)
}

func (set *Set) Elems() []*Item {
	return set.elems
}

func (set *Set) DoAdd(it *Item) {
	set.elems = append(set.elems, it)
}

type ItemStack struct {
	stack *container.Stack
}

func NewItemStack() *ItemStack {
	return &ItemStack{
		stack: container.NewStack(),
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
