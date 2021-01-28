package lr

import "github.com/Orlion/merak/production"

type GrammarState struct {
	state       int
	ps          []*production.Production
	closureZSet *ProductionZSet
}

func NewGrammerState(state int, ps []*production.Production) *GrammarState {
	return &GrammarState{
		state: state,
		ps:    ps,
	}
}

func (gs *GrammarState) createTransition() {
	gs.makeClosure()
}

func (gs *GrammarState) makeClosure() {
	gs.closureZSet = NewProductionZSet()

	gs.closureZSet.addAll(gs.ps)

	pStack := NewProductionStack()

	for _, p := range gs.ps {
		pStack.Push(p)
	}

	for !pStack.Empty() {
		p := pStack.Pop()
		if p.IsDotEnd() {
			continue
		}
		s := p.GetDotSymbol()
		if s.IsTerminal() {
			continue
		}

		closures := gs.gsm.pm.getProductions(s)

		lookAhead := gs.gsm.pm.computeFirstSetOfBetaAndC(p)

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
