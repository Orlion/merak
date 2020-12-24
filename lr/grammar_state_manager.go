package lr

import (
	"sort"
	"strings"

	"github.com/Orlion/merak/symbol"
)

type GrammarStateManager struct {
	stateNumCount int
	states        map[string]*GrammarState
	gs            *GrammarState
	pm            *ProductionManager
}

func NewGrammarStateManager(pm *ProductionManager) *GrammarStateManager {
	gsm := &GrammarStateManager{
		stateNumCount: 0,
		states:        make(map[string]*GrammarState),
		pm:            pm,
	}

	return gsm
}

func (gsm *GrammarStateManager) GenLrActionTable(goal symbol.Symbol, eoi symbol.Symbol) map[int]map[symbol.Symbol]*Action {
	gsm.pm.fsb.genFirstSets()

	ps := gsm.pm.getProductions(goal)
	if len(ps) != 1 {
		panic("The goal symbol's production can only be one")
	}

	ps[0].lookAhead.Add(eoi)
	gs := newGrammarState(gsm, gsm.stateNumCount, gsm.pm.getProductions(goal))

	gsm.states[gsm.key(ps)] = gs

	gs.createTransition()

	// lr action table
	m := make(map[int]map[symbol.Symbol]*Action)

	for _, gs := range gsm.states {
		jump := make(map[symbol.Symbol]*Action)
		for s, childGs := range gs.transition {
			if _, exists := jump[s]; exists {
				panic("shift conflict")
			}
			jump[s] = newShiftAction(childGs.stateNum)
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

	return m
}

func (gsm *GrammarStateManager) key(ps []*Production) string {
	keyList := make([]string, 0)
	for _, p := range ps {
		keyList = append(keyList, p.getCode())
	}

	sort.Strings(keyList)
	key := strings.Join(keyList, " | ")

	return key
}

func (gsm *GrammarStateManager) getGrammarState(productions []*Production) *GrammarState {
	key := gsm.key(productions)

	if s, exists := gsm.states[key]; exists {
		return s
	} else {
		gsm.stateNumCount++
		gs := newGrammarState(gsm, gsm.stateNumCount, productions)
		gsm.states[key] = gs
		return gs
	}
}
