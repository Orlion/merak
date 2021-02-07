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

func (state *GrammarState) createTransition() {
	if state.transitionDone {
		return
	}
	state.transitionDone = true

	state.makeClosure()

	state.makePartition()

	state.makeTransition()

	state.extendTransition()
}

func (state *GrammarState) makeClosure() {
	state.closureSet = item.NewItemSet()
	state.closureSet.AddList(state.its)

	for _, it := range state.its {
		state.atb.logger.Println(it.DotPos())
	}

	itemStack := item.NewItemStack()

	for _, it := range state.its {
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

		state.atb.logger.Println(s, it.DotPos())
		closures := state.atb.itm.GetItems(s)

		lookAhead := it.ComputeFirstSetOfBetaAndC(state.atb.fs)

		for _, oldItem := range closures {
			newItem := oldItem.CloneWithLookAhead(lookAhead)
			if !state.closureSet.Exists(newItem) {
				state.closureSet.Add(newItem)

				itemStack.Push(newItem)

				state.removeRedundantProduction(newItem)
			}
		}
	}
}

func (state *GrammarState) removeRedundantProduction(newItem *item.Item) {
	for it := range state.closureSet.Elems() {
		if newItem.IsCoverUp(it) {
			state.closureSet.Delete(it)
		}
	}
}

func (state *GrammarState) makePartition() {
	state.partition = make(map[symbol.Symbol][]*item.Item)
	for it := range state.closureSet.Elems() {
		if !it.DotEnd() {
			state.partition[it.GetDotSymbol()] = append(state.partition[it.GetDotSymbol()], it)
		}
	}
}

func (state *GrammarState) makeTransition() {
	state.transition = make(map[symbol.Symbol]*GrammarState)

	for symbol, its := range state.partition {
		newStateItems := make([]*item.Item, 0, len(its))
		for _, it := range its {
			newStateItems = append(newStateItems, it.DotForward())
		}

		newState := state.atb.newGrammarState(newStateItems)
		state.transition[symbol] = newState
	}
}

func (state *GrammarState) extendTransition() {
	for _, childState := range state.transition {
		if !childState.transitionDone {
			childState.createTransition()
		}
	}
}

func (state *GrammarState) makeReduce() map[symbol.Symbol]*Action {
	m := make(map[symbol.Symbol]*Action)
	for _, it := range state.its {
		if it.CanBeReduce() {
			for s := range it.GetLookAhead().Elems() {
				m[s] = NewReduceAction(it.GetProduction().GetCallback(), it.GetProduction().ParamsLen())
			}
		}
	}

	return m
}
