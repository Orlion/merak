package lr

import (
	"errors"
	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/symbol"
)

type ActionTableBuilder struct {
	gsList      []*GrammarState
	itm         *item.Manager
	fs          *first_set.FirstSet
	lastStateId int
}

func NewActionTableBuilder(itm *item.Manager, fs *first_set.FirstSet) *ActionTableBuilder {
	return &ActionTableBuilder{
		itm:    itm,
		fs:     fs,
		gsList: make([]*GrammarState, 0),
	}
}

func (atb *ActionTableBuilder) newGrammarState(its []*item.Item, fromId int, fromSymbol symbol.Symbol) (gs *GrammarState) {
FindGsLoop:
	// TODO: Improve performance
	for _, gs := range atb.gsList {
		if len(gs.its) != len(its) {
			continue
		}

		for _, it := range gs.its {
			find := false
			for _, inputIt := range its {
				if it.Equals(inputIt) {
					find = true
					break
				}
			}
			if !find {
				continue FindGsLoop
			}
		}

		return gs
	}
	gs = NewGrammarState(atb.lastStateId, its, atb, fromId, fromSymbol)
	atb.lastStateId++
	atb.gsList = append(atb.gsList, gs)
	return
}

func (atb *ActionTableBuilder) Build(goal symbol.Symbol, eoi symbol.Symbol) (at *ActionTable, err error) {
	its := atb.itm.GetItems(goal)
	if len(its) > 1 {
		err = errors.New("goal has no any productions")
		return
	}

	for _, it := range its {
		if it.GetProduction().ParamsLen() != 1 {
			err = errors.New("goal's production can only be one parameter")
			return
		}
		lookAhead := symbol.NewSymbolSet()
		lookAhead.Add(eoi)
		it.SetLookAhead(lookAhead)
	}

	gs := atb.newGrammarState(its, -1, nil)
	gs.createTransition()

	at = NewActionTable()

	for _, gs := range atb.gsList {
		jump := make(map[symbol.Symbol]*Action)
		for s, childGs := range gs.transition {
			if _, exists := jump[s]; exists {
				err = errors.New("shift conflict")
				return
			}

			jump[s] = NewShiftAction(childGs.id)
		}

		reduceMap := gs.makeReduce()
		for s, action := range reduceMap {
			if _, exists := jump[s]; exists {
				err = errors.New("shift reduce conflict")
				return
			}
			jump[s] = action
		}

		at.add(gs.id, jump)
	}

	return
}
