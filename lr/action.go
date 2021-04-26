package lr

import (
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/symbol"
)

type ActionType int8

const (
	ActionReduce ActionType = iota + 1
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
	action.t = ActionReduce
	return action
}

func NewShiftAction(state int) *Action {
	action := new(Action)
	action.state = state
	action.t = ActionShift
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
