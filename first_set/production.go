package first_set

import "github.com/Orlion/merak/symbol"

type Production struct {
	result     symbol.Symbol
	paramsList [][]symbol.Symbol
	firstSet   *symbol.Set
}

func NewProduction(result symbol.Symbol, params []symbol.Symbol) *Production {
	return &Production{
		result:     result,
		paramsList: [][]symbol.Symbol{params},
		firstSet:   symbol.NewSymbolSet(),
	}
}

func (p *Production) AddParams(params []symbol.Symbol) {
	p.paramsList = append(p.paramsList, params)
}

func (p *Production) AddFirstSet(s symbol.Symbol) bool {
	return p.firstSet.Add(s)
}

func (p *Production) AddAllFirstSet(set *symbol.Set) int {
	return p.firstSet.AddAll(set)
}

func (p *Production) GetResult() symbol.Symbol {
	return p.result
}

func (p *Production) GetParamsList() [][]symbol.Symbol {
	return p.paramsList
}

func (p *Production) GetFirstSet() *symbol.Set {
	return p.firstSet
}
