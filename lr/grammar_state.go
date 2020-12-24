package lr

import (
	"github.com/Orlion/merak/symbol"
)

type GrammarState struct {
	gsm            *GrammarStateManager
	stateNum       int
	productions    []*Production
	transition     map[symbol.Symbol]*GrammarState
	closureZSet    *ProductionZSet
	partition      map[symbol.Symbol][]*Production
	transitionDone bool
}

func newGrammarState(gsm *GrammarStateManager, stateNum int, productions []*Production) *GrammarState {
	return &GrammarState{
		gsm:         gsm,
		stateNum:    stateNum,
		productions: productions,
		closureZSet: newProductionZSet(),
	}
}

func (gs *GrammarState) createTransition() {
	gs.makeClosure()

	gs.makePartition()

	gs.makeTransition()

	gs.extendTransition()
}

func (gs *GrammarState) makeClosure() {
	gs.closureZSet.addAll(gs.productions)

	pStack := newProductionStack()

	for _, p := range gs.productions {
		pStack.Push(p)
	}

	for !pStack.Empty() {
		production := pStack.Pop()
		if production.isDotEnd() {
			continue
		}
		s := production.getDotSymbol()
		if s.IsTerminals() {
			continue
		}

		closures := gs.gsm.pm.getProductions(s)

		lookAhead := gs.gsm.pm.computeFirstSetOfBetaAndC(production)

		for _, oldProduct := range closures {
			newProduction := oldProduct.cloneSelf()
			newProduction.addLookAheadSet(lookAhead)

			if !gs.closureZSet.exists(newProduction) {
				gs.closureZSet.add(newProduction)

				pStack.Push(newProduction)

				gs.removeRedundantProduction(newProduction)
			}
		}
	}
}

func (gs *GrammarState) removeRedundantProduction(newProduction *Production) {
	target := newProductionZSet()
	for _, item := range gs.closureZSet.list {
		if newProduction.isCoverUp(item) {
			continue
		}

		target.add(item)
	}

	gs.closureZSet = target
}

func (gs *GrammarState) makePartition() {
	gs.partition = make(map[symbol.Symbol][]*Production)
	for _, p := range gs.closureZSet.list {
		if !p.isDotEnd() {
			gs.partition[p.getDotSymbol()] = append(gs.partition[p.getDotSymbol()], p)
		} else {
			gs.productions = append(gs.productions, p)
		}
	}
}

func (gs *GrammarState) makeTransition() {
	var newGs *GrammarState
	var newGsPs []*Production

	gs.transition = make(map[symbol.Symbol]*GrammarState)

	for symbol, ps := range gs.partition {
		newGsPs = []*Production{}
		for _, p := range ps {
			newGsPs = append(newGsPs, p.dotForward())
		}

		newGs = gs.gsm.getGrammarState(newGsPs)
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
	for _, p := range gs.productions {
		if p.canBeReduce() {
			for _, s := range p.lookAhead.List {
				m[s] = newReduceAction(p)
			}
		}
	}

	return m
}
