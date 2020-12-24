package lr

import (
	"github.com/Orlion/merak/symbol"
)

type FirstSetBuilder struct {
	pass bool
	m    map[symbol.Symbol]*symbol.Production
}

func newFirstSetBuilder() *FirstSetBuilder {
	return &FirstSetBuilder{
		m: make(map[symbol.Symbol]*symbol.Production),
	}
}

func (fsb *FirstSetBuilder) register(s symbol.Symbol, productions []symbol.Symbol, nullable bool) {
	if sp, exists := fsb.m[s]; exists {
		if !s.IsTerminals() {
			sp.AddProduction(productions)
		}
	} else {
		fsb.m[s] = symbol.NewProduction(s, nullable, productions)
	}
}

func (fsb *FirstSetBuilder) getFirstZSet(s symbol.Symbol) *symbol.ZSet {
	return fsb.m[s].GetFirstZSet()
}

func (fsb *FirstSetBuilder) isSymbolNullable(s symbol.Symbol) bool {
	return false
}

func (fsb *FirstSetBuilder) genFirstSets() {
	count := 0
	fsb.pass = true
	for fsb.pass {
		if count > 100 {
			panic("After 100 attempts, the FirstSet generation is still not completed.")
		}

		fsb.pass = false
		for _, sp := range fsb.m {
			fsb.genProductionFirstSet(sp)
		}

		count++
	}
}

func (fsb *FirstSetBuilder) genProductionFirstSet(sp *symbol.Production) {
	if sp.GetValue().IsTerminals() {
		if sp.AddFirstZSet(sp.GetValue()) {
			fsb.pass = true
		}
	} else {
		for _, p := range sp.GetProductions() {
			if p[0].IsTerminals() {
				if sp.AddFirstZSet(p[0]) {
					fsb.pass = true
				}
			} else {
				for _, curSymbol := range p {
					curSp := fsb.m[curSymbol]
					curSet := curSp.GetFirstZSet()

					for _, s := range curSet.List {
						if sp.AddFirstZSet(s) {
							fsb.pass = true
						}
					}

					if !curSp.IsNullable() {
						break
					}
				}
			}
		}
	}
}
