package item

import "github.com/Orlion/merak/symbol"

type Callback func(params ...symbol.Value) symbol.Value

type Production struct {
	id       int
	result   symbol.Symbol
	params   []symbol.Symbol
	callback Callback
}

func NewProduction(id int, result symbol.Symbol, params []symbol.Symbol, callback Callback) *Production {
	if !result.IsTerminal() {
		panic("result must be a terminal")
	}

	if callback == nil {
		panic("callback cannot be a nil")
	}

	return &Production{id, result, params, callback}
}

func (p *Production) ParamsLen() int {
	return len(p.params)
}

func (p *Production) GetParam(idx int) symbol.Symbol {
	return p.params[idx]
}

func (p *Production) GetCallback() Callback {
	return p.callback
}

func (p *Production) GetId() int {
	return p.id
}
