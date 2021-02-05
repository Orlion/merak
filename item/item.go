package item

import (
	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/symbol"
)

type Item struct {
	production *Production
	lookAhead  *symbol.Set
	dotPos     int
}

func NewItem(production *Production, dotPos int) *Item {
	return &Item{
		production: production,
		dotPos:     dotPos,
		lookAhead:  symbol.NewSymbolSet(),
	}
}

func (it *Item) DotForward() *Item {
	newItem := NewItem(it.production, it.dotPos+1)
	newItem.lookAhead.AddAll(it.lookAhead)

	return newItem
}

func (it *Item) DotEnd() bool {
	return it.dotPos >= it.production.ParamsLen()
}

func (it *Item) GetDotSymbol() symbol.Symbol {
	return it.production.GetParam(it.dotPos)
}

func (it *Item) CloneWithLookAhead(lookAhead *symbol.Set) *Item {
	newItem := NewItem(it.production, it.dotPos)
	newItem.lookAhead.AddAll(lookAhead)

	return newItem
}

func (it *Item) IsCoverUp(oldItem *Item) bool {
	if it.productionEquals(oldItem) && it.lookAheadCompare(oldItem) > 0 {
		return true
	}

	return false
}

func (it *Item) productionEquals(input *Item) bool {
	return it.production.GetId() == input.production.GetId()
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
	return it.dotPos >= it.production.ParamsLen()
}

func (it *Item) GetLookAhead() *symbol.Set {
	return it.lookAhead
}

func (it *Item) GetProduction() *Production {
	return it.production
}

func (it *Item) ComputeFirstSetOfBetaAndC(fs *first_set.FirstSet) (firstSet *symbol.Set) {
	firstSet = symbol.NewSymbolSet()

	if len(it.production.params) > it.dotPos+1 {
		firstSet.AddAll(fs.Get(it.production.params[it.dotPos+1]))
	} else {
		firstSet.AddAll(it.lookAhead)
	}

	return
}
