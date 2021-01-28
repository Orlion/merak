package lr

import (
	"errors"

	"github.com/Orlion/merak/production"
	"github.com/Orlion/merak/symbol"
)

type ActionTableBuilder struct {
	gsList []*GrammarState
}

func NewActionTableBuilder() *ActionTableBuilder {
	return &ActionTableBuilder{}
}

func (atb *ActionTableBuilder) Build(pm *production.Manager, goal symbol.Symbol) (at *ActionTable, err error) {
	ps := pm.GetProductions(goal)

	if len(ps) < 1 {
		err = errors.New("goal has no productions")
		return
	}

	gs := NewGrammerState(len(atb.gsList), ps)
	gs.createTransition()

	at = NewActionTable()

	for _, gs := range atb.gsList {
		jump := make(map[symbol.Symbol]*Action)
		for s, childGs := range gs.transition {
			if _, exists := jump[s]; exists {
				panic("shift conflict")
			}
			jump[s] = NewShiftAction(childGs.stateNum)
		}

		reduceMap := gs.makeReduce()
		for s, action := range reduceMap {
			if _, exists := jump[s]; exists {
				panic("shift reduce conflict")
			}
			jump[s] = action
		}

		m[gs.stateNum] = jump
	}

	return
}
