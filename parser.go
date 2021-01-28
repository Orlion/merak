package merak

import (
	"errors"
	"fmt"

	"github.com/Orlion/merak/lexer"
	"github.com/Orlion/merak/lr"
	"github.com/Orlion/merak/production"
	"github.com/Orlion/merak/symbol"
)

var (
	SyntaxErr = errors.New("Parse error: syntax error")
)

type Parser struct {
	pNum      int
	atBuilder *lr.ActionTableBuilder
	at        *lr.ActionTable
}

func NewParser() *Parser {
	atBuilder := lr.NewActionTableBuilder()
	return &Parser{
		atBuilder: atBuilder,
	}
}

// Register production to this parser
func (parser *Parser) RegProduction(left symbol.Symbol, rights []symbol.Symbol, nullable bool, callback production.Callback) {

}

func (parser *Parser) buildActionTable(goal symbol.Symbol) (err error) {
	if parser.at != nil {
		return
	}

	parser.at, err = parser.atBuilder.Build(nil, goal)
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
		j                uint
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

	stateStack := lr.NewStack()
	valueStack := lr.NewStack()
	symbolStack := lr.NewStack()
	stateStack.Push(0)

	for {
		state = stateStack.Top().(int)

		action, err = parser.at.Action(state, token)
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

			for j = 0; j < paramsNum/2; j++ {
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
