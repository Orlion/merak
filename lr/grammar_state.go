package lr

import (
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/symbol"
)

type GrammarState struct {
	its            []*item.Item
	partition      map[symbol.Symbol][]*item.Item
	transition     map[symbol.Symbol]*GrammarState
	closureSet     *item.ItemSet
	atb            *ActionTableBuilder
	id             int
	transitionDone bool
}

func NewGrammarState(id int, its []*item.Item, atb *ActionTableBuilder) *GrammarState {
	return &GrammarState{
		its: its,
		atb: atb,
		id:  id,
	}
}

func (gs *GrammarState) createTransition() {
	if gs.transitionDone {
		return
	}
	gs.transitionDone = true

	gs.makeClosure()
}

func (gs *GrammarState) makeClosure() {
	gs.closureSet = item.NewItemSet()
	gs.closureSet.AddList(gs.its)

	itemStack := item.NewItemStack()

	for _, it := range gs.its {
		itemStack.Push(it)
	}

	for !itemStack.Empty() {
		it := itemStack.Pop()
		if it.DotEnd() {
			continue
		}
		s := it.GetDotSymbol()
		if s.IsTerminal() {
			continue
		}

		closures := gs.atb.itm.GetItems(s)

		lookAhead := gs.atb.itm.ComputeFirstSetOfBetaAndC(it, gs.atb.fs)

		for _, oldItem := range closures {
			newItem := oldItem.CloneWithLookAhead(lookAhead)

			if !gs.closureSet.Exists(newItem) {
				gs.closureSet.Add(newItem)

				itemStack.Push(newItem)

				gs.removeRedundantProduction(newItem)
			}
		}
	}
}

func (gs *GrammarState) removeRedundantProduction(newItem *item.Item) {
	for it := range gs.closureSet.Elems() {
		if newItem.IsCoverUp(it) {
			gs.closureSet.Delete(it)
		}
	}
}

func (gs *GrammarState) makePartition() {
	gs.partition = make(map[symbol.Symbol][]*item.Item)
	for it := range gs.closureSet.Elems() {
		if !it.DotEnd() {
			gs.partition[it.GetDotSymbol()] = append(gs.partition[it.GetDotSymbol()], it)
		} else {
			gs.its = append(gs.its, it)
		}
	}
}

func (gs *GrammarState) makeTransition() {
	var newGs *GrammarState
	var newGsIts []*item.Item

	gs.transition = make(map[symbol.Symbol]*GrammarState)

	for symbol, its := range gs.partition {
		newGsIts = []*item.Item{}
		for _, it := range its {
			newGsIts = append(newGsIts, it.DotForward())
		}

		newGs = gs.atb.getGrammarState(newGsIts)
		gs.transition[symbol] = newGs
	}

	gs.transitionDone = true
}

func (gs *GrammarState) extendTransition() {
	for _, childGs := range gs.transition {
		if !childGs.transitionDone {
			childGs.createTransition()
		}
	}
}

func (gs *GrammarState) makeReduce() map[symbol.Symbol]*Action {
	m := make(map[symbol.Symbol]*Action)
	for _, it := range gs.its {
		if it.CanBeReduce() {
			for s := range it.GetLookAhead().Elems() {
				m[s] = NewReduceAction(it.GetProduction().GetCallback(), it.GetProduction().ParamsLen())
			}
		}
	}

	return m
}
