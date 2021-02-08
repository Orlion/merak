package item

import (
	"testing"

	"github.com/Orlion/merak/symbol"
)

type testSymbol string

func (s testSymbol) IsTerminal() bool {
	return string(s)[0] < 'A' || string(s)[0] > 'Z'
}

func (s testSymbol) ToString() string {
	return string(s)
}

func TestAdd(t *testing.T) {
	set := NewSet()
	p1, _ := NewProduction(1, testSymbol("EXPR"), []symbol.Symbol{testSymbol("EXPR")}, func(params ...symbol.Value) symbol.Value {
		return nil
	})
	it1 := NewItem(p1, 0)
	set.Add(it1)
	if set.Add(it1) {
		t.Errorf("set.Add(it1) true")
		t.FailNow()
	}
	it2 := NewItem(p1, 0)
	if set.Add(it2) {
		t.Errorf("set.Add(it2) true")
		t.FailNow()
	}
}
