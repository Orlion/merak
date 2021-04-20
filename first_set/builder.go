package first_set

import (
	"errors"

	"github.com/Orlion/merak/symbol"
)

type Builder struct {
	pass bool
	m    map[symbol.Symbol]*Production
}

func NewBuilder() *Builder {
	return &Builder{
		pass: true,
		m:    make(map[symbol.Symbol]*Production),
	}
}

func (b *Builder) Reg(s symbol.Symbol, params []symbol.Symbol) {
	if p, exists := b.m[s]; exists {
		p.AddParams(params)
	} else {
		b.m[s] = NewProduction(s, params)
	}

	for _, param := range params {
		if param.IsTerminal() {
			if _, exists := b.m[param]; !exists {
				b.m[param] = NewProduction(param, []symbol.Symbol{param})
			}
		}
	}
}

// build first set
func (b *Builder) Build() (fs *FirstSet, err error) {
	count := 0
	for b.pass {
		if count > 1024 {
			err = errors.New("After 1024 attempts, the FirstSet generation is still not completed.")
			return
		}

		b.pass = false
		for _, p := range b.m {
			b.genProductionFirstSet(p)
		}

		count++
	}

	fs = NewFirstSet(b.m)

	return
}

func (b *Builder) genProductionFirstSet(p *Production) {
	if p.GetResult().IsTerminal() {
		// If p.result is a terminator then it is its own FirstSet element
		if p.AddFirstSet(p.GetResult()) {
			b.pass = true
		}
	} else {
		for _, params := range p.GetParamsList() {
			if params[0].IsTerminal() {
				if p.AddFirstSet(params[0]) {
					b.pass = true
				}
			} else {
				if p.AddAllFirstSet(b.m[params[0]].GetFirstSet()) > 0 {
					b.pass = true
				}
			}
		}
	}
}
