package lr

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/symbol"
)

type ActionTableBuilder struct {
	gsList   []*GrammarState
	itm      *item.Manager
	fs       *first_set.FirstSet
	states   map[string]*GrammarState
	stateNum int
}

func NewActionTableBuilder(itm *item.Manager, fs *first_set.FirstSet) *ActionTableBuilder {
	return &ActionTableBuilder{
		itm: itm,
		fs:  fs,
	}
}

func (atb *ActionTableBuilder) Build(goal symbol.Symbol) (at *ActionTable, err error) {
	its := atb.itm.GetItems(goal)

	if len(its) < 1 {
		err = errors.New("goal has no productions")
		return
	}

	gs := NewGrammarState(len(atb.gsList), its, atb)
	gs.createTransition()

	at = NewActionTable()

	for _, gs := range atb.gsList {
		jump := make(map[symbol.Symbol]*Action)
		for s, childGs := range gs.transition {
			if _, exists := jump[s]; exists {
				panic("shift conflict")
			}
			jump[s] = NewShiftAction(childGs.state)
		}

		reduceMap := gs.makeReduce()
		for s, action := range reduceMap {
			if _, exists := jump[s]; exists {
				panic("shift reduce conflict")
			}
			jump[s] = action
		}

		at.add(gs.state, jump)
	}

	return
}

func (atb *ActionTableBuilder) getGrammarState(its []*item.Item) *GrammarState {
	key := atb.key(its)
	if s, exists := atb.states[key]; exists {
		return s
	} else {
		atb.stateNum++
		gs := NewGrammarState(atb.stateNum, its, atb)
		atb.states[key] = gs
		return gs
	}
}

func (atb *ActionTableBuilder) key(its []*item.Item) string {
	keyList := make([]string, 0)
	for _, it := range its {
		keyList = append(keyList, strconv.Itoa(it.Id()))
	}

	sort.Strings(keyList)

	return strings.Join(keyList, " | ")
}
