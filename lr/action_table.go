package lr

import (
	"errors"

	"github.com/Orlion/merak/symbol"
)

type ActionTable struct {
	m map[int]map[symbol.Symbol]*Action
}

func NewActionTable() *ActionTable {
	return &ActionTable{
		m: make(map[int]map[symbol.Symbol]*Action),
	}
}

func (at *ActionTable) add(state int, jump map[symbol.Symbol]*Action) {
	at.m[state] = jump
}

func (at *ActionTable) Action(state int, s symbol.Symbol) (action *Action, err error) {
	jump, exists := at.m[state]
	if !exists {
		err = errors.New("Unexpected token")
		return
	}

	action, exists = jump[s]
	if !exists {
		err = errors.New("Unexpected token")
		return
	}

	return
}
