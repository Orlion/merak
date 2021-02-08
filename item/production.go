package item

import (
	"errors"

	"github.com/Orlion/merak/symbol"
)

type Callback func(params ...symbol.Value) symbol.Value

type Production struct {
	id       int
	result   symbol.Symbol
	params   []symbol.Symbol
	callback Callback
}

func NewProduction(id int, result symbol.Symbol, params []symbol.Symbol, callback Callback) (production *Production, err error) {
	if callback == nil {
		err = errors.New("callback cannot be a nil")
		return
	}

	production = &Production{id, result, params, callback}
	return
}

func (p *Production) ParamsLen() int {
	return len(p.params)
}

func (p *Production) GetParam(idx int) symbol.Symbol {
	return p.params[idx]
}

func (p *Production) GetParams() []symbol.Symbol {
	return p.params
}

func (p *Production) GetCallback() Callback {
	return p.callback
}

func (p *Production) GetId() int {
	return p.id
}

func (p *Production) GetResult() symbol.Symbol {
	return p.result
}
