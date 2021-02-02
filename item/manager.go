package item

import (
	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/symbol"
)

type Manager struct {
	m       map[symbol.Symbol][]*Item
	itemNum int
}

func NewManager() *Manager {
	return &Manager{
		m: make(map[symbol.Symbol][]*Item),
	}
}

func (m *Manager) RegItem(result symbol.Symbol, params []symbol.Symbol, callback Callback) {
	it := NewItem(m.itemNum, result, params, callback, 0)
	m.itemNum++
	if _, exists := m.m[result]; exists {
		m.m[result] = append(m.m[result], it)
	} else {
		m.m[result] = []*Item{it}
	}
}

func (m *Manager) GetItems(result symbol.Symbol) []*Item {
	return m.m[result]
}

func (m *Manager) ComputeFirstSetOfBetaAndC(it *Item, fs *first_set.FirstSet) (firstSet *symbol.Set) {
	set := symbol.NewSymbolSet()

	for i := it.dotPos + 1; i < len(it.params); i++ {
		set.Add(it.params[i])
	}

	firstSet = symbol.NewSymbolSet()

	if len(set.Elems()) > 0 {
		for s := range set.Elems() {
			firstSet.AddAll(fs.Get(s))
			break
		}
	} else {
		firstSet.AddAll(it.lookAhead)
	}

	return
}
