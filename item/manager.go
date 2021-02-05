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

func (m *Manager) RegProduction(result symbol.Symbol, params []symbol.Symbol, callback Callback) {
	it := NewProduction(m.lastPid, result, params, callback)
	m.lastPid++
	m.m[result] = append(m.m[result], it)
}

func (m *Manager) GetProductions(result symbol.Symbol) []*Production {
	return m.m[result]
}
