package first_set

import "github.com/Orlion/merak/symbol"

type ProductionState int8

const (
	ProductionStateNotRun ProductionState = iota + 1
	ProductionStateRunning
	ProductionStateDone
)

type Production struct {
	value      symbol.Symbol
	paramsList [][]symbol.Symbol
	firstZSet  *symbol.ZSet
	nullable   bool
	state      ProductionState
}

func NewProduction(s symbol.Symbol, params []symbol.Symbol) *Production {
	return &Production{
		value:      s,
		paramsList: [][]symbol.Symbol{params},
		firstZSet:  symbol.NewSymbolZSet(),
	}
}

func (p *Production) AddParams(params []symbol.Symbol) {
	p.paramsList = append(p.paramsList, params)
}

func (p *Production) AddFirstZSet(s symbol.Symbol) bool {
	return p.firstZSet.Add(s)
}

func (p *Production) AddAllFirstZSet(set *symbol.ZSet) {
	p.firstZSet.AddAll(set)
}

func (p *Production) GetValue() symbol.Symbol {
	return p.value
}

func (p *Production) GetParamsList() [][]symbol.Symbol {
	return p.paramsList
}

func (p *Production) GetFirstZSet() *symbol.ZSet {
	return p.firstZSet
}

func (p *Production) GetState() ProductionState {
	return p.state
}

func (p *Production) SetState(state ProductionState) {
	p.state = state
}
