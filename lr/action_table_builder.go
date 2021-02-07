package lr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/log"
	"github.com/Orlion/merak/symbol"
)

type ActionTableBuilder struct {
	gsList      []*GrammarState
	itm         *item.Manager
	fs          *first_set.FirstSet
	logger      log.Logger
	lastStateId int
}

func NewActionTableBuilder(itm *item.Manager, fs *first_set.FirstSet, logger log.Logger) *ActionTableBuilder {
	return &ActionTableBuilder{
		logger: logger,
		itm:    itm,
		fs:     fs,
		gsList: make([]*GrammarState, 0),
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

func (atb *ActionTableBuilder) print() {
	for _, gs := range atb.gsList {
		atb.logger.Println(fmt.Sprintf("%d:", gs.id))
		for _, it := range gs.its {
			itPrintBuilder := new(strings.Builder)
			itPrintBuilder.WriteString(fmt.Sprintf("%s ->   ", it.GetProduction().GetResult()))
			for k, v := range it.GetProduction().GetParams() {
				if it.DotPos() == k {
					itPrintBuilder.WriteString(".   ")
				}
				itPrintBuilder.WriteString(v.ToString())
				itPrintBuilder.WriteString("   ")
			}

			itPrintBuilder.WriteString("(")

			for s := range it.GetLookAhead().Elems() {
				itPrintBuilder.WriteString(s.ToString())
				itPrintBuilder.WriteString(" ")
			}

			itPrintBuilder.WriteString(")")

			atb.logger.Println(itPrintBuilder.String())
		}
	}
}
