package merak

import (
	"errors"
	"fmt"

	"github.com/Orlion/merak/container"
	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/lexer"
	"github.com/Orlion/merak/lr"
	"github.com/Orlion/merak/symbol"
)

var (
	SyntaxErr = errors.New("syntax error")
)

func newSyntaxErr(unexpected, expecting lexer.Token) error {
	if expecting == nil {
		return fmt.Errorf("%s:%d:%d: %w: unexpected %s", unexpected.Filename(), unexpected.Line(), unexpected.Col(), SyntaxErr, unexpected.ToString())
	} else {
		return fmt.Errorf("%s:%d:%d: %w: unexpected %s, ecpecting %s", unexpected.Filename(), unexpected.Line(), unexpected.Col(), SyntaxErr, unexpected.ToString(), expecting.ToString())
	}
}

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
func (parser *Parser) RegProduction(result symbol.Symbol, params []symbol.Symbol, callback item.Callback) (err error) {
	if result.IsTerminal() {
		err = errors.New("result cannot be a terminator")
		return
	}

	parser.itm.Reg(result, params, callback)
	parser.fsb.Reg(result, params)

	return
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
		state              int
		action             *lr.Action
		args               []symbol.Value
		currentSymbolValue symbol.Value
		currentSymbol      symbol.Symbol
	)

	err = parser.buildActionTable(goal)
	if err != nil {
		return
	}

	lexerDelegator := lexer.NewLexerDelegator(l)

	token, err := lexerDelegator.Next()
	if err != nil {
		return
	}

	currentSymbolValue = token.ToSymbol()
	currentSymbol = currentSymbolValue.Symbol()
	stateStack := container.NewStack()
	valueStack := container.NewStack()
	symbolStack := container.NewStack()
	stateStack.Push(0)

	for {
		state = stateStack.Top().(int)

		action, err = parser.at.Action(state, currentSymbol)
		if err != nil {
			err = newSyntaxErr(token, nil)
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
			currentSymbol = result.Symbol()
			symbolStack.Push(currentSymbol)
			valueStack.Push(result)
		case lr.ActionShift:
			stateStack.Push(action.State())
			symbolStack.Push(currentSymbol)
			if currentSymbol.IsTerminal() {
				valueStack.Push(currentSymbolValue)
				if token, err = lexerDelegator.Next(); err != nil {
					break
				}
			}

			currentSymbolValue = token.ToSymbol()
			currentSymbol = currentSymbolValue.Symbol()
		case lr.ActionAccept:
			break
		}
	}

	return
}
