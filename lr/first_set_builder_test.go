package lr

import (
	"testing"
	"unicode"

	"github.com/Orlion/merak/symbol"
)

type TestSymbol string

const (
	Symbola TestSymbol = "a"
	Symbolb TestSymbol = "b"
	SymbolC TestSymbol = "C"
	Symbold TestSymbol = "d"
	SymbolE TestSymbol = "E"
	Symbolf TestSymbol = "f"
	SymbolG TestSymbol = "G"
	SymbolH TestSymbol = "H"
	SymbolI TestSymbol = "I"
)

func (s TestSymbol) ToString() string {
	return string(s)
}

func (s TestSymbol) IsTerminals() bool {
	r := unicode.IsUpper([]rune(s)[0])
	return r
}

func TestGenFirstSets(t *testing.T) {
	fsb := newFirstSetBuilder()
	fsb.register(Symbola, []symbol.Symbol{Symbolb}, false)
	fsb.register(Symbola, []symbol.Symbol{Symbolf}, false)
	fsb.register(Symbola, []symbol.Symbol{SymbolC}, false)
	fsb.register(Symbolb, []symbol.Symbol{Symbold}, false)
	fsb.register(Symbolb, []symbol.Symbol{SymbolE}, false)
	fsb.register(Symbolb, []symbol.Symbol{SymbolG}, false)
	fsb.register(Symbolb, []symbol.Symbol{SymbolH}, false)
	fsb.register(Symbold, []symbol.Symbol{Symbola}, false)
	fsb.register(Symbolf, []symbol.Symbol{Symbolb}, false)
	fsb.register(Symbolf, []symbol.Symbol{SymbolI}, false)
	fsb.genFirstSets()
	if !fsb.getFirstZSet(Symbola).Exists(SymbolC) {
		t.Fatal("SymbolA 's FirstSet not contains SymbolC")
	}

	if !fsb.getFirstZSet(Symbola).Exists(SymbolE) {
		t.Fatal("SymbolA 's FirstSet not contains SymbolE")
	}

	if !fsb.getFirstZSet(Symbola).Exists(SymbolG) {
		t.Fatal("SymbolA 's FirstSet not contains SymbolG")
	}

	if !fsb.getFirstZSet(Symbola).Exists(SymbolH) {
		t.Fatal("SymbolA 's FirstSet not contains SymbolH")
	}

	if !fsb.getFirstZSet(Symbola).Exists(SymbolI) {
		t.Fatal("SymbolA 's FirstSet not contains SymbolI")
	}

	if !fsb.getFirstZSet(Symbolb).Exists(SymbolE) {
		t.Fatal("Symbolb 's FirstSet not contains SymbolE")
	}

	if !fsb.getFirstZSet(Symbolb).Exists(SymbolG) {
		t.Fatal("Symbolb 's FirstSet not contains SymbolG")
	}

	if !fsb.getFirstZSet(Symbolb).Exists(SymbolH) {
		t.Fatal("Symbolb 's FirstSet not contains SymbolH")
	}

	if !fsb.getFirstZSet(Symbolb).Exists(SymbolC) {
		t.Fatal("Symbolb 's FirstSet not contains SymbolC")
	}

	if !fsb.getFirstZSet(Symbolb).Exists(SymbolI) {
		t.Fatal("Symbolb 's FirstSet not contains SymbolI")
	}
}
