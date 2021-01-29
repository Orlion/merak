package lr

import (
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/symbol"
)

type GrammarState struct {
	its            []*item.Item
	partition      map[symbol.Symbol][]*item.Item
	transition     map[symbol.Symbol]*GrammarState
	closureZSet    *item.ItemZSet
	atb            *ActionTableBuilder
	state          int
	transitionDone bool
}

func NewGrammarState(state int, its []*item.Item, atb *ActionTableBuilder) *GrammarState {
	return &GrammarState{
		its:   its,
		atb:   atb,
		state: state,
	}
}

func (gs *GrammarState) createTransition() {
	gs.makeClosure()
}

func (gs *GrammarState) makeClosure() {
	gs.closureZSet = item.NewItemZSet()

	gs.closureZSet.AddList(gs.its)

	itemStack := item.NewItemStack()

	for _, it := range gs.its {
		itemStack.Push(it)
	}

	for !itemStack.Empty() {
		it := itemStack.Pop()
		if it.IsDotEnd() {
			continue
		}
		s := it.GetDotSymbol()
		if s.IsTerminal() {
			continue
		}

		closures := gs.atb.itm.GetItems(s)

		lookAhead := gs.atb.itm.ComputeFirstSetOfBetaAndC(it, gs.atb.fs)

		for _, oldItem := range closures {
			newItem := oldItem.Clone()
			newItem.AddLookAheadSet(lookAhead)

			if !gs.closureZSet.Exists(newItem) {
				gs.closureZSet.Add(newItem)

				itemStack.Push(newItem)

				gs.removeRedundantProduction(newItem)
			}
		}
	}
}

func (gs *GrammarState) removeRedundantProduction(newItem *item.Item) {
	target := item.NewItemZSet()
	for _, it := range gs.closureZSet.List() {
		if newItem.IsCoverUp(it) {
			continue
		}

		target.Add(it)
	}

	gs.closureZSet = target
}

func (gs *GrammarState) makePartition() {
	gs.partition = make(map[symbol.Symbol][]*item.Item)
	for _, it := range gs.closureZSet.List() {
		if !it.IsDotEnd() {
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
			for _, s := range it.GetLookAhead().List() {
				m[s] = NewReduceAction(it.GetCallback(), it.ParamsLen())
			}
		}
	}

	return m
}
