package production

import "github.com/Orlion/merak/symbol"

type Callback func(params ...symbol.Value) symbol.Value

type Production struct {
	id       int
	result   symbol.Symbol
	params   []symbol.Symbol
	callback Callback
}

func NewProduction(id int, params []symbol.Symbol, result symbol.Symbol, callback Callback) *Production {
	if !result.IsTerminal() {
		panic("result must be a terminal")
	}

	if callback == nil {
		panic("callback cannot be a nil")
	}

	return &Production{
		id, result, params, callback,
	}
}

func (p *Production) DotForward() {

}

func (p *Production) IsDotEnd() bool {
	return false
}

func (p *Production) GetDotSymbol() symbol.Symbol {
	return nil
}

func (p *Production) Id() int {
	return p.id
}
