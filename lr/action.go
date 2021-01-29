package lr

import (
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/symbol"
)

type ActionType int8

const (
	ActionAccept ActionType = iota + 1
	ActionReduce
	ActionShift
)

type Action struct {
	ReduceAction
	ShiftAction
	t ActionType
}

func (action *Action) Type() ActionType {
	return action.t
}

func NewReduceAction(callback item.Callback, paramsNum int) *Action {
	action := new(Action)
	action.callback = callback
	action.paramsNum = paramsNum
	return action
}

func NewShiftAction(state int) *Action {
	action := new(Action)
	action.state = state
	return action
}

type ReduceAction struct {
	callback  item.Callback
	paramsNum int
}

func (action *ReduceAction) Reduce(params ...symbol.Value) symbol.Value {
	return action.callback(params...)
}

func (action *ReduceAction) ParamsNum() int {
	return action.paramsNum
}

type ShiftAction struct {
	state int
}

func (action *ShiftAction) State() int {
	return action.state
}
