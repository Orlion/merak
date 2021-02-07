package merak

import (
	"errors"
	"fmt"

	"github.com/Orlion/merak/data_structure"
	"github.com/Orlion/merak/first_set"
	"github.com/Orlion/merak/item"
	"github.com/Orlion/merak/lexer"
	"github.com/Orlion/merak/log"
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
	logger log.Logger
	at     *lr.ActionTable
	itm    *item.Manager
	fsb    *first_set.Builder
}

func NewParser(logger log.Logger) *Parser {
	return &Parser{
		logger: logger,
		itm:    item.NewManager(),
		fsb:    first_set.NewBuilder(),
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
	parser.at, err = lr.NewActionTableBuilder(parser.itm, fs, parser.logger).Build(goal)
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

	token, err := lexerDelegator.Next()
	if err != nil {
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
			symbolStack.Push(result.Symbol())
			valueStack.Push(result)
		case lr.ActionShift:
			tokenSymbolValue = token.ToSymbol()
			symbolStack.Push(tokenSymbolValue.Symbol())
			valueStack.Push(tokenSymbolValue)
			stateStack.Push(action.State())
			if token, err = lexerDelegator.Next(); err != nil {
				break
			}

		case lr.ActionAccept:
			break
		}
	}

	return
}
