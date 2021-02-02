package symbol

type ProductionState int8

const (
	ProductionStateNotRun ProductionState = iota + 1
	ProductionStateRunning
	ProductionStateDone
)

type Production struct {
	value       Symbol
	productions [][]Symbol
	firstSet    *Set
	nullable    bool
	state       ProductionState
}

func NewProduction(symbol Symbol, nullable bool, production []Symbol) *Production {
	return &Production{
		value:       symbol,
		nullable:    nullable,
		productions: [][]Symbol{production},
		firstSet:    NewSymbolSet(),
	}
}

func (sp *Production) AddProduction(production []Symbol) {
	sp.productions = append(sp.productions, production)
}

func (sp *Production) AddFirstSet(s Symbol) bool {
	return sp.firstSet.Add(s)
}

func (sp *Production) AddAllFirstZSet(set *Set) {
	sp.firstSet.AddAll(set)
}

func (sp *Production) GetValue() Symbol {
	return sp.value
}

func (sp *Production) GetProductions() [][]Symbol {
	return sp.productions
}

func (sp *Production) GetFirstSet() *Set {
	return sp.firstSet
}

func (sp *Production) IsNullable() bool {
	return sp.nullable
}

func (sp *Production) GetState() ProductionState {
	return sp.state
}

func (sp *Production) SetState(state ProductionState) {
	sp.state = state
}
