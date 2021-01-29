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

func (b *Builder) Register(s symbol.Symbol, params []symbol.Symbol) {
	if p, exists := b.m[s]; exists {
		if !s.IsTerminal() {
			p.AddParams(params)
		}
	} else {
		b.m[s] = NewProduction(s, params)
	}
}

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
	if p.GetValue().IsTerminal() {
		if p.AddFirstZSet(p.GetValue()) {
			b.pass = true
		}
	} else {
		for _, params := range p.GetParamsList() {
			if params[0].IsTerminal() {
				if p.AddFirstZSet(params[0]) {
					b.pass = true
				}
			} else {
				for _, curSymbol := range params {
					curP := b.m[curSymbol]
					curSet := curP.GetFirstZSet()

					for _, s := range curSet.List() {
						if p.AddFirstZSet(s) {
							b.pass = true
						}
					}

					break
				}
			}
		}
	}
}
