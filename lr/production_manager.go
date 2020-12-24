package lr

import "github.com/Orlion/merak/symbol"

type ProductionManager struct {
	m   map[symbol.Symbol][]*Production
	fsb *FirstSetBuilder
}

func NewProductionManager() *ProductionManager {
	return &ProductionManager{
		m:   make(map[symbol.Symbol][]*Production),
		fsb: newFirstSetBuilder(),
	}
}

func (pm *ProductionManager) Register(left symbol.Symbol, rights []symbol.Symbol, nullable bool, builder AstNodeBuilder) {
	production := newProduction(left, rights, 0, builder)
	if _, exists := pm.m[left]; exists {
		pm.m[left] = append(pm.m[left], production)
	} else {
		pm.m[left] = []*Production{production}
	}

	if len(rights) > 0 {
		pm.fsb.register(left, rights, nullable)
	}
	for _, s := range rights {
		if s.IsTerminals() {
			pm.fsb.register(s, []symbol.Symbol{s}, false)
		}
	}
}

func (pm *ProductionManager) getProductions(left symbol.Symbol) []*Production {
	return pm.m[left]
}

func (pm *ProductionManager) computeFirstSetOfBetaAndC(production *Production) *symbol.ZSet {
	set := symbol.NewSymbolZSet()

	for i := production.dotPos + 1; i < len(production.rights); i++ {
		set.Add(production.rights[i])
	}

	firstSet := symbol.NewSymbolZSet()

	if len(set.List) > 0 {
		for _, s := range set.List {
			lookAhead := pm.fsb.getFirstZSet(s)
			firstSet.AddAll(lookAhead)

			if !pm.fsb.isSymbolNullable(s) {
				break
			}
		}
	} else {
		firstSet.AddAll(production.lookAhead)
	}

	return firstSet
}
