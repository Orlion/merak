package lr

import (
	"errors"

	"github.com/Orlion/merak/lexer"
)

type ActionTable struct {
	m map[int]map[lexer.Token]*Action
}

func NewActionTable() *ActionTable {
	return &ActionTable{
		m: make(map[int]map[lexer.Token]*Action),
	}
}

func (at *ActionTable) Action(state int, token lexer.Token) (action *Action, err error) {
	jump, exists := at.m[state]
	if !exists {
		err = errors.New("Unexpected token")
		return
	}

	action, exists = jump[token]
	if !exists {
		err = errors.New("Unexpected token")
		return
	}

	return
}
