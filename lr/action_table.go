package lr

import (
	"errors"
	"fmt"

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

func (at *ActionTable) Print() {
	for state, cm := range at.m {
		fmt.Println("====================================")
		fmt.Printf("%d  =>  ", state)
		for s, action := range cm {
			fmt.Printf("%s / ", s)
			switch action.Type() {
			case ActionAccept:
				fmt.Printf("accept\n")
			case ActionReduce:
				fmt.Printf("reduce\n")
			case ActionShift:
				fmt.Printf("shift-%d\n", action.State())
			}
		}
	}
}
