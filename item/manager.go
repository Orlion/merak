package item

import (
	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/symbol"
)

type Manager struct {
	ps map[symbol.Symbol][]*Item
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) GetItems(result symbol.Symbol) []*Item {
	return m.ps[result]
}

func (m *Manager) ComputeFirstSetOfBetaAndC(it *Item, fs *first_set.FirstSet) (firstSet *symbol.ZSet) {
	set := symbol.NewSymbolZSet()

	for i := it.dotPos + 1; i < len(it.params); i++ {
		set.Add(it.params[i])
	}

	firstSet = symbol.NewSymbolZSet()

	if len(set.List()) > 0 {
		for _, s := range set.List() {
			firstSet.AddAll(fs.Get(s))
			break
		}
	} else {
		firstSet.AddAll(it.lookAhead)
	}

	return
}
