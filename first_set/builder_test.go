package first_set

import (
	"testing"

	"github.com/Orlion/merak/symbol"
)

type Symbol string

func (s Symbol) IsTerminal() bool {
	return string(s)[0] < 'A' || string(s)[0] > 'Z'
}

const (
	SymbolEoi    Symbol = "eoi"
	SymbolAdd    Symbol = "+"
	SymbolSub    Symbol = "-"
	SymbolMul    Symbol = "*"
	SymbolDiv    Symbol = "/"
	SymbolLp     Symbol = "("
	SymbolRp     Symbol = ")"
	SymbolNumber Symbol = "number"
	SymbolFACTOR Symbol = "FACTOR"
	SymbolTERM   Symbol = "TERM"
	SymbolEXPR   Symbol = "EXPR"
	SymbolGOAL   Symbol = "GOAL"
)

func TestBuild(t *testing.T) {
	fsb := NewBuilder()
	/*
		GOAL -> EXPR eoi
		EXPR -> TERM
			  | EXPR + TERM
			  | EXPR - TERM
		TERM -> FACTOR
			  | term * FACTOR
			  | term / FACTOR
		FACTOR -> number
		        | '(' EXPR ')'
	*/
	fsb.Reg(SymbolFACTOR, []symbol.Symbol{SymbolNumber})
	fsb.Reg(SymbolFACTOR, []symbol.Symbol{SymbolLp, SymbolNumber, SymbolRp})
	fsb.Reg(SymbolTERM, []symbol.Symbol{SymbolFACTOR})
	fsb.Reg(SymbolTERM, []symbol.Symbol{SymbolTERM, SymbolMul, SymbolFACTOR})
	fsb.Reg(SymbolTERM, []symbol.Symbol{SymbolTERM, SymbolDiv, SymbolFACTOR})
	fsb.Reg(SymbolEXPR, []symbol.Symbol{SymbolTERM})
	fsb.Reg(SymbolEXPR, []symbol.Symbol{SymbolEXPR, SymbolAdd, SymbolTERM})
	fsb.Reg(SymbolEXPR, []symbol.Symbol{SymbolEXPR, SymbolSub, SymbolTERM})
	fsb.Reg(SymbolGOAL, []symbol.Symbol{SymbolEXPR, SymbolEoi})
	fs, err := fsb.Build()
	if err != nil {
		t.Errorf("fsb.Build err: %s", err.Error())
		t.FailNow()
	}

	if _, exists := fs.Get(SymbolGOAL).Elems()[SymbolNumber]; !exists {
		t.Errorf("SymbolGOAL's FirstSet do not have SymbolNumber")
		t.FailNow()
	}

	if _, exists := fs.Get(SymbolGOAL).Elems()[SymbolLp]; !exists {
		t.Errorf("SymbolGOAL's FirstSet do not have SymbolLp")
		t.FailNow()
	}

	if _, exists := fs.Get(SymbolFACTOR).Elems()[SymbolLp]; !exists {
		t.Errorf("SymbolFACTOR's FirstSet do not have SymbolLp")
		t.FailNow()
	}
}
