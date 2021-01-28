package production

import "github.com/Orlion/merak/symbol"

type Manager struct {
	ps map[symbol.Symbol][]*Production
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) GetProductions(result symbol.Symbol) []*Production {
	return m.ps[result]
}
