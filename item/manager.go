package item

import (
	"github.com/Orlion/merak/symbol"
)

type Manager struct {
	m       map[symbol.Symbol][]*Production
	lastPid int
}

func NewManager() *Manager {
	return &Manager{
		m: make(map[symbol.Symbol][]*Production),
	}
}

func (m *Manager) Reg(result symbol.Symbol, params []symbol.Symbol, callback Callback) (err error) {
	it, err := NewProduction(m.lastPid, result, params, callback)
	if err != nil {
		return
	}
	m.lastPid++
	m.m[result] = append(m.m[result], it)

	return
}

func (m *Manager) GetProductions(result symbol.Symbol) []*Production {
	return m.m[result]
}

func (m *Manager) GetItems(result symbol.Symbol) []*Item {
	its := make([]*Item, 0, len(m.m[result]))
	for _, production := range m.m[result] {
		its = append(its, NewItem(production, 0))
	}

	return its
}
