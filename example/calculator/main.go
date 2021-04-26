package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Orlion/merak"
	"github.com/Orlion/merak/lexer"
	"github.com/Orlion/merak/symbol"
)

type TokenType int

const (
	TokenEoi TokenType = iota
	TokenAdd
	TokenSub
	TokenMul
	TokenDiv
	TokenLp
	TokenRp
	TokenNumber
)

type Token struct {
	Text string
	Type TokenType
	line int
	col  int
}

func NewToken(text string, tokenType TokenType, col int) *Token {
	return &Token{
		Text: text,
		Type: tokenType,
		col:  col,
	}
}

func (t *Token) ToString() string {
	return t.Text
}

func (t *Token) Filename() string {
	return "input"
}

func (t *Token) Line() int {
	return t.line
}

func (t *Token) Col() int {
	return t.col
}

func (t *Token) Symbol() symbol.Symbol {
	m := map[TokenType]Symbol{
		TokenEoi:    SymbolEoi,
		TokenAdd:    SymbolAdd,
		TokenSub:    SymbolSub,
		TokenMul:    SymbolMul,
		TokenDiv:    SymbolDiv,
		TokenLp:     SymbolLp,
		TokenRp:     SymbolRp,
		TokenNumber: SymbolNumber,
	}

	return m[t.Type]
}

func (t *Token) ToSymbol() symbol.Value {
	return t
}

type Symbol string

func (s Symbol) IsTerminal() bool {
	return string(s)[0] < 'A' || string(s)[0] > 'Z'
}

func (s Symbol) ToString() string {
	return string(s)
}

const (
	SymbolEoi    Symbol = "eoi"
	SymbolAdd    Symbol = "+"
	SymbolSub    Symbol = "-"
	SymbolMul    Symbol = "*"
	SymbolDiv    Symbol = "/"
	SymbolLp     Symbol = "("
	SymbolRp     Symbol = ")"
	SymbolNumber Symbol = "number"
	SymbolFACTOR Symbol = "FACTOR"
	SymbolTERM   Symbol = "TERM"
	SymbolEXPR   Symbol = "EXPR"
	SymbolGOAL   Symbol = "GOAL"
)

type Lexer struct {
	tokens []*Token
	pos    int
}

func (l *Lexer) Next() (lexer.Token, error) {
	if l.pos > len(l.tokens) {
		return nil, errors.New("eoi err")
	} else if l.pos == len(l.tokens) {
		return &Token{Text: "", Type: TokenEoi}, nil
	} else {
		l.pos = l.pos + 1
		return l.tokens[l.pos-1], nil
	}
}

type Goal struct {
	expr *Expr
}

func (v *Goal) Symbol() symbol.Symbol {
	return SymbolGOAL
}

type Expr struct {
	expr  *Expr
	term  *Term
	isAdd bool
}

func (v *Expr) Symbol() symbol.Symbol {
	return SymbolEXPR
}

type Term struct {
	term   *Term
	factor *Factor
	isMul  bool
}

func (v *Term) Symbol() symbol.Symbol {
	return SymbolTERM
}

type Factor struct {
	expr   *Expr
	number int64
}

func (v *Factor) Symbol() symbol.Symbol {
	return SymbolFACTOR
}

func main() {
	parser := initParser()

	l := &Lexer{
		tokens: []*Token{
			NewToken("123", TokenNumber, 1),
			NewToken("-", TokenSub, 2),
			NewToken("456", TokenNumber, 3),
			//NewToken("", TokenNumber, 4),
		},
	}

	r, err := parser.Parse(SymbolGOAL, SymbolEoi, l)
	if err != nil {
		panic(err)
	}

	goal := r.(*Goal)
	fmt.Println(interpreter(goal))
}

func initParser() *merak.Parser {
	parser := merak.NewParser()
	/*
		GOAL -> EXPR
		EXPR -> TERM
			  | EXPR + TERM
			  | EXPR - TERM
		TERM -> FACTOR
			  | TERM * FACTOR
			  | TERM / FACTOR
		FACTOR -> number
		        | '(' EXPR ')'
	*/
	parser.RegProduction(SymbolFACTOR, []symbol.Symbol{SymbolNumber}, func(params ...symbol.Value) symbol.Value {
		number, _ := strconv.ParseInt(params[0].(*Token).Text, 10, 64)
		return &Factor{
			number: number,
		}
	})

	parser.RegProduction(SymbolFACTOR, []symbol.Symbol{SymbolLp, SymbolEXPR, SymbolRp}, func(params ...symbol.Value) symbol.Value {
		return &Factor{
			expr: params[1].(*Expr),
		}
	})

	parser.RegProduction(SymbolTERM, []symbol.Symbol{SymbolFACTOR}, func(params ...symbol.Value) symbol.Value {
		return &Term{
			factor: params[0].(*Factor),
		}
	})

	parser.RegProduction(SymbolTERM, []symbol.Symbol{SymbolTERM, SymbolMul, SymbolFACTOR}, func(params ...symbol.Value) symbol.Value {
		return &Term{
			term:   params[0].(*Term),
			factor: params[2].(*Factor),
			isMul:  true,
		}
	})

	parser.RegProduction(SymbolTERM, []symbol.Symbol{SymbolTERM, SymbolDiv, SymbolFACTOR}, func(params ...symbol.Value) symbol.Value {
		return &Term{
			term:   params[0].(*Term),
			factor: params[2].(*Factor),
			isMul:  false,
		}
	})

	parser.RegProduction(SymbolEXPR, []symbol.Symbol{SymbolTERM}, func(params ...symbol.Value) symbol.Value {
		return &Expr{
			term: params[0].(*Term),
		}
	})

	parser.RegProduction(SymbolEXPR, []symbol.Symbol{SymbolEXPR, SymbolAdd, SymbolTERM}, func(params ...symbol.Value) symbol.Value {
		return &Expr{
			expr:  params[0].(*Expr),
			term:  params[2].(*Term),
			isAdd: true,
		}
	})

	parser.RegProduction(SymbolEXPR, []symbol.Symbol{SymbolEXPR, SymbolSub, SymbolTERM}, func(params ...symbol.Value) symbol.Value {
		return &Expr{
			expr:  params[0].(*Expr),
			term:  params[2].(*Term),
			isAdd: false,
		}
	})

	parser.RegProduction(SymbolGOAL, []symbol.Symbol{SymbolEXPR}, func(params ...symbol.Value) symbol.Value {
		return &Goal{
			expr: params[0].(*Expr),
		}
	})

	return parser
}

func interpreter(goal *Goal) int64 {
	return evalGoal(goal)
}

func evalGoal(goal *Goal) int64 {
	return evalExpr(goal.expr)
}

func evalExpr(expr *Expr) int64 {
	term := evalTerm(expr.term)
	if expr.expr != nil {
		if expr.isAdd {
			return evalExpr(expr.expr) + term
		} else {
			return evalExpr(expr.expr) - term
		}
	} else {
		return term
	}
}

func evalTerm(term *Term) int64 {
	factor := evalFactor(term.factor)
	if term.term != nil {
		if term.isMul {
			return evalTerm(term.term) * factor
		} else {
			return evalTerm(term.term) / factor
		}
	} else {
		return factor
	}
}

func evalFactor(factor *Factor) int64 {
	return factor.number
}
