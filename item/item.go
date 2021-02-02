package item

import "github.com/Orlion/merak/symbol"

type Callback func(params ...symbol.Value) symbol.Value

type Item struct {
	params    []symbol.Symbol
	callback  Callback
	result    symbol.Symbol
	lookAhead *symbol.Set
	dotPos    int
	id        int
}

func NewItem(id int, result symbol.Symbol, params []symbol.Symbol, callback Callback, dotPos int) *Item {
	if !result.IsTerminal() {
		panic("result must be a terminal")
	}

	if callback == nil {
		panic("callback cannot be a nil")
	}

	return &Item{
		params:   params,
		callback: callback,
		result:   result,
		dotPos:   dotPos,
		id:       id,
	}
}

func (it *Item) DotForward() *Item {
	newItem := NewItem(it.id, it.result, it.params, it.callback, it.dotPos+1)
	newItem.lookAhead.AddAll(it.lookAhead)

	return newItem
}

func (it *Item) IsDotEnd() bool {
	return false
}

func (it *Item) GetDotSymbol() symbol.Symbol {
	return nil
}

func (it *Item) Id() int {
	return it.id
}

func (it *Item) Clone() *Item {
	newProduction := NewItem(it.id, it.result, it.params, it.callback, it.dotPos)
	newProduction.lookAhead.AddAll(it.lookAhead)

	return newProduction
}

func (it *Item) AddLookAheadSet(lookAhead *symbol.Set) {
	it.lookAhead = lookAhead
}

func (it *Item) IsCoverUp(oldItem *Item) bool {
	if it.productionEquals(oldItem) && it.lookAheadCompare(oldItem) > 0 {
		return true
	}

	return false
}

func (it *Item) productionEquals(input *Item) bool {
	if it.result != input.result {
		return false
	}

	if it.dotPos != input.dotPos {
		return false
	}

	if len(it.params) != len(input.params) {
		return false
	}

	for k, v := range it.params {
		if v != input.params[k] {
			return false
		}
	}

	return true
}

func (it *Item) lookAheadCompare(input *Item) int {
	if len(it.lookAhead.Elems()) < len(input.lookAhead.Elems()) {
		return -1
	}

	if len(it.lookAhead.Elems()) > len(input.lookAhead.Elems()) {
		for s := range input.lookAhead.Elems() {
			if !it.lookAhead.Exists(s) {
				return -1
			}
		}
		return 1
	}

	for s := range it.lookAhead.Elems() {
		if !input.lookAhead.Exists(s) {
			return -1
		}
	}

	return 0
}

func (it *Item) CanBeReduce() bool {
	return it.dotPos >= len(it.params)
}

func (it *Item) GetLookAhead() *symbol.Set {
	return it.lookAhead
}

func (it *Item) GetCallback() Callback {
	return it.callback
}

func (it *Item) ParamsLen() int {
	return len(it.params)
}
