package lr

import (
	"fmt"
	"strings"

	"github.com/Orlion/merak/ast"
	"github.com/Orlion/merak/symbol"
)

type AstNodeBuilder func([]interface{}) ast.Node

type Production struct {
	left      symbol.Symbol
	rights    []symbol.Symbol
	dotPos    int
	Builder   AstNodeBuilder
	lookAhead *symbol.ZSet
	code      string
}

func newProduction(left symbol.Symbol, rights []symbol.Symbol, dotPos int, builder AstNodeBuilder) *Production {
	if builder == nil {
		panic("AstNodeBuilder cannot be nil")
	}

	return &Production{
		left:      left,
		rights:    rights,
		dotPos:    dotPos,
		Builder:   builder,
		lookAhead: symbol.NewSymbolZSet(),
	}
}

func (p *Production) Rights() []symbol.Symbol {
	return p.rights
}

func (p *Production) Left() symbol.Symbol {
	return p.left
}

func (p *Production) dotForward() *Production {
	newProduction := newProduction(p.left, p.rights, p.dotPos+1, p.Builder)
	newProduction.lookAhead.AddAll(p.lookAhead)

	return newProduction
}

func (p *Production) isDotEnd() bool {
	return p.dotPos >= len(p.rights)
}

func (p *Production) getDotSymbol() symbol.Symbol {
	return p.rights[p.dotPos]
}

func (p *Production) isCoverUp(oldProduction *Production) bool {
	if p.productionEquals(oldProduction) && p.lookAheadCompare(oldProduction) > 0 {
		return true
	}

	return false
}

func (p *Production) productionEquals(production *Production) bool {
	if p.left != production.left {
		return false
	}

	if p.dotPos != production.dotPos {
		return false
	}

	if len(p.rights) != len(production.rights) {
		return false
	}

	for k, v := range p.rights {
		if v != production.rights[k] {
			return false
		}
	}

	return true
}

func (p *Production) lookAheadCompare(production *Production) int {
	if len(p.lookAhead.List) < len(production.lookAhead.List) {
		return -1
	}

	if len(p.lookAhead.List) > len(production.lookAhead.List) {
		for _, s := range production.lookAhead.List {
			if !p.lookAhead.Exists(s) {
				return -1
			}
		}
		return 1
	}

	for _, s := range p.lookAhead.List {
		if !p.lookAhead.Exists(s) {
			return -1
		}
	}

	return 0
}

func (p *Production) addLookAheadSet(lookAhead *symbol.ZSet) {
	p.lookAhead = lookAhead
}

func (p *Production) cloneSelf() *Production {
	newProduction := newProduction(p.left, p.rights, p.dotPos, p.Builder)
	newProduction.lookAhead.AddAll(p.lookAhead)

	return newProduction
}

func (p *Production) canBeReduce() bool {
	return p.dotPos >= len(p.rights)
}

func (p *Production) getCode() string {
	if p.code == "" {
		var codeBuilder strings.Builder
		codeBuilder.WriteString(fmt.Sprintf("%s ->   ", p.left.ToString()))
		for k, v := range p.rights {
			if p.dotPos == k {
				codeBuilder.WriteString(".   ")
			}
			codeBuilder.WriteString(v.ToString())
			codeBuilder.WriteString("   ")
		}

		codeBuilder.WriteString("(")

		list := make([]string, 0)

		for _, k := range p.lookAhead.List {
			list = append(list, k.ToString())
		}

		for _, s := range list {
			codeBuilder.WriteString(s)
			codeBuilder.WriteString(" ")
		}

		codeBuilder.WriteString(")")

		p.code = codeBuilder.String()
	}

	return p.code
}
