package merak

import (
	"errors"
	"fmt"

	"github.com/Orlion/merak/ast"
	"github.com/Orlion/merak/lexer"
	"github.com/Orlion/merak/lr"
	"github.com/Orlion/merak/symbol"
)

var (
	SyntaxErr = errors.New("Parse error: syntax error")
)

type Parser struct {
	goal          symbol.Symbol
	lexer         *lexer.LexerDelegate
	pm            *lr.ProductionManager
	lrActionTable map[int]map[symbol.Symbol]*lr.Action
}

func NewParser() *Parser {
	return &Parser{
		pm: lr.NewProductionManager(),
	}
}

// Add production to this parser
func (parser *Parser) RegisterProduction(left symbol.Symbol, rights []symbol.Symbol, nullable bool, builder lr.AstNodeBuilder) {
	parser.pm.Register(left, rights, nullable, builder)
}

// Build LR action table
func (parser *Parser) Build(goal symbol.Symbol, eoi symbol.Symbol) *Parser {
	gsm := lr.NewGrammarStateManager(parser.pm)
	parser.lrActionTable = gsm.GenLrActionTable(goal, eoi)
	parser.goal = goal
	return parser
}

func (parser *Parser) SetLexer(coreLexer lexer.Lexer) *Parser {
	parser.lexer = lexer.NewLexerDelegate(coreLexer)
	return parser
}

// Parse Input
func (parser *Parser) Parse() (ast ast.Node, err error) {
	var (
		currentState int
		action       *lr.Action
		args         []interface{}
	)

	token, lexerErr := parser.lexer.Next()
	if lexerErr != nil {
		err = fmt.Errorf("Lexer error: [%w]", lexerErr)
		return
	}

	valueToken := token

	currentSymbol := token.ToSymbol()

	stateStack := lr.NewStack()
	valueStack := lr.NewStack()
	symbolStack := lr.NewStack()
	stateStack.Push(0)
	valueStack.Push(token)

	for {
		currentState = stateStack.Top().(int)

		action, err = parser.getAction(currentState, currentSymbol)
		if err != nil {
			break
		}

		if action.IsReduce() { // reduce
			args = []interface{}{}

			for i := len(action.ReduceProduction().Rights()); i > 0; i-- {
				symbolStack.Pop()
				stateStack.Pop()
				args = append(args, valueStack.Pop())
			}

			for i := 0; i < len(args)/2; i++ {
				temp := args[i]
				args[i] = args[len(args)-1-i]
				args[len(args)-1-i] = temp
			}

			if action.ReduceProduction().Left() == parser.goal {
				ast = action.ReduceProduction().Builder(args)
				return
			} else {
				symbolStack.Push(action.ReduceProduction().Left())

				currentSymbol = action.ReduceProduction().Left()

				valueStack.Push(action.ReduceProduction().Builder(args))
			}
		} else { // shift
			stateStack.Push(action.ShiftStateNum())

			symbolStack.Push(currentSymbol)

			if currentSymbol.IsTerminals() {
				valueStack.Push(valueToken)
				token, lexerErr = parser.lexer.Next()
				if lexerErr != nil {
					err = fmt.Errorf("Lexer error: [%w]", lexerErr)
					return
				}
			} else {
				token = parser.lexer.Current()
			}

			currentSymbol = token.ToSymbol()
			valueToken = token
		}
	}

	return
}

func (parser *Parser) getAction(currentState int, s symbol.Symbol) (action *lr.Action, err error) {
	jump, exists := parser.lrActionTable[currentState]
	if !exists {
		err = fmt.Errorf("%w, unexpected '%s'", SyntaxErr, parser.lexer.Current().ToString())
		return
	}

	action, exists = jump[s]
	if !exists {
		err = fmt.Errorf("%w, unexpected '%s'", SyntaxErr, parser.lexer.Current().ToString())
		return
	}

	return
}
