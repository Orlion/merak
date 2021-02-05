package merak

import (
	"errors"
	"fmt"

	"github.com/Orlion/merak/data_structure"
	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/lexer"
	"github.com/Orlion/merak/lr"
	"github.com/Orlion/merak/symbol"
)

var (
	SyntaxErr = errors.New("Parse error: syntax error")
)

type Parser struct {
	at  *lr.ActionTable
	itm *item.Manager
	fsb *first_set.Builder
}

func NewParser() *Parser {
	return &Parser{
		itm: item.NewManager(),
		fsb: first_set.NewBuilder(),
	}
}

// Register production to this parser
func (parser *Parser) RegProduction(result symbol.Symbol, params []symbol.Symbol, callback item.Callback) {
	parser.itm.RegProduction(result, params, callback)
	parser.fsb.Reg(result, params)
}

func (parser *Parser) buildActionTable(goal symbol.Symbol) (err error) {
	if parser.at != nil {
		return
	}

	// build first set
	fs, err := parser.fsb.Build()
	if err != nil {
		return
	}

	// build action table
	parser.at, err = lr.NewActionTableBuilder(parser.itm, fs).Build(goal)
	if err != nil {
		return
	}

	return
}

// Parse Input
func (parser *Parser) Parse(goal symbol.Symbol, l lexer.Lexer) (result symbol.Value, err error) {
	var (
		state            int
		action           *lr.Action
		args             []symbol.Value
		tokenSymbolValue symbol.Value
	)

	err = parser.buildActionTable(goal)
	if err != nil {
		return
	}

	lexerDelegator := lexer.NewLexerDelegator(l)

	token, lexerErr := lexerDelegator.Next()
	if lexerErr != nil {
		err = fmt.Errorf("Lexer error: [%w]", lexerErr)
		return
	}

	stateStack := data_structure.NewStack()
	valueStack := data_structure.NewStack()
	symbolStack := data_structure.NewStack()
	stateStack.Push(0)

	for {
		state = stateStack.Top().(int)

		action, err = parser.at.Action(state, token.ToSymbol().Symbol())
		if err != nil {
			err = fmt.Errorf("%w [%s]", SyntaxErr, err.Error())
			break
		}

		switch action.Type() {
		case lr.ActionReduce:
			args = []symbol.Value{}

			paramsNum := action.ParamsNum()
			for i := paramsNum; i > 0; i-- {
				stateStack.Pop()
				symbolStack.Pop()
				args = append(args, valueStack.Pop().(symbol.Value))
			}

			for j := 0; j < paramsNum/2; j++ {
				temp := args[j]
				args[j] = args[paramsNum-1-j]
				args[paramsNum-1-j] = temp
			}

			result = action.Reduce(args...)
			symbolStack.Push(result.Symbol())
			valueStack.Push(result)

			state = action.State()

			stateStack.Push(state)
		case lr.ActionShift:
			tokenSymbolValue = token.ToSymbol()
			symbolStack.Push(tokenSymbolValue.Symbol())
			valueStack.Push(tokenSymbolValue)
			stateStack.Push(action.State())
			if token, err = lexerDelegator.Next(); err != nil {
				err = fmt.Errorf("Lexer error: [%w]", lexerErr)
				break
			}

		case lr.ActionAccept:
			break
		}
	}

	return
}
