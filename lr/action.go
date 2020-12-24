package lr

type Action struct {
	isReduce         bool        // true: reduce，false: shift
	reduceProduction *Production // reduce production
	shiftStateNum    int         // shift to stateNum
}

func (action *Action) IsReduce() bool {
	return action.isReduce
}

func (action *Action) ReduceProduction() *Production {
	return action.reduceProduction
}

func (action *Action) ShiftStateNum() int {
	return action.shiftStateNum
}

func newReduceAction(reduceProduction *Production) *Action {
	return &Action{isReduce: true, reduceProduction: reduceProduction}
}

func newShiftAction(shiftStateNum int) *Action {
	return &Action{isReduce: false, shiftStateNum: shiftStateNum}
}
