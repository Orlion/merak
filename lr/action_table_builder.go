package lr

import (
	"errors"
	"fmt"

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
		itm:         itm,
		fs:          fs,
		gsList:      make([]*GrammarState, 0),
		lastStateId: 1,
	}
}

func (atb *ActionTableBuilder) newGrammarState(its []*item.Item) (gs *GrammarState) {
	gs = NewGrammarState(atb.lastStateId, its, atb)
	atb.lastStateId++
	atb.gsList = append(atb.gsList, gs)
	return
}

func (atb *ActionTableBuilder) Build(goal symbol.Symbol) (at *ActionTable, err error) {
	its := atb.itm.GetItems(goal)
	if len(its) < 1 {
		err = errors.New("goal has no any productions")
		return
	}

	gs := atb.newGrammarState(its)
	gs.createTransition()
	// print gs
	atb.print()

	at = NewActionTable()

	for _, gs := range atb.gsList {
		jump := make(map[symbol.Symbol]*Action)
		for s, childGs := range gs.transition {
			if _, exists := jump[s]; exists {
				err = errors.New("shift conflict")
				return
			}
			if 0 == childGs.id {
				jump[s] = NewAcceptAction()
			} else {
				jump[s] = NewShiftAction(childGs.id)
			}
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

func (atb *ActionTableBuilder) print() {
	for _, gs := range atb.gsList {
		fmt.Printf("%d:\n", gs.id)
		for _, it := range gs.its {
			fmt.Println(it.ToString())
		}
	}
}
