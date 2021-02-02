package main

import (
	"errors"

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
}

func NewToken(text string, tokenType TokenType) *Token {
	return &Token{
		Text: text,
		Type: tokenType,
	}
}

func (t *Token) ToSymbol() symbol.Value {
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

	return NewSymbolValue(m[t.Type])
}

type SymbolValue struct {
	s symbol.Symbol
}

func NewSymbolValue(s symbol.Symbol) *SymbolValue {
	return &SymbolValue{
		s: s,
	}
}

func (sv *SymbolValue) Symbol() symbol.Symbol {
	return sv.s
}

type Symbol string

func (s Symbol) IsTerminal() bool {
	return string(s)[0] >= 'A' && string(s)[0] <= 'Z'
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

func main() {
	parser := initParser()

	l := &Lexer{
		tokens: []*Token{NewToken("123", TokenNumber), NewToken("+", TokenAdd), NewToken("456", TokenNumber)},
	}

	parser.Parse(SymbolGOAL, l)
}

func initParser() *merak.Parser {
	parser := merak.NewParser()
	/*
		GOAL -> EXPR eoi
		EXPR -> term
			  | EXPR + term
			  | EXPR - term
		TERM -> FACTOR
			  | term * FACTOR
			  | term / FACTOR
		FACTOR -> number
		        | '(' EXPR ')'
	*/
	parser.RegProduction(SymbolFACTOR, []symbol.Symbol{SymbolNumber}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	parser.RegProduction(SymbolFACTOR, []symbol.Symbol{SymbolLp, SymbolNumber, SymbolRp}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	parser.RegProduction(SymbolTERM, []symbol.Symbol{SymbolFACTOR}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	parser.RegProduction(SymbolTERM, []symbol.Symbol{SymbolTERM, SymbolMul, SymbolFACTOR}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	parser.RegProduction(SymbolTERM, []symbol.Symbol{SymbolTERM, SymbolDiv, SymbolFACTOR}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	parser.RegProduction(SymbolEXPR, []symbol.Symbol{SymbolTERM}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	parser.RegProduction(SymbolEXPR, []symbol.Symbol{SymbolEXPR, SymbolAdd, SymbolTERM}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	parser.RegProduction(SymbolEXPR, []symbol.Symbol{SymbolEXPR, SymbolSub, SymbolTERM}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	parser.RegProduction(SymbolGOAL, []symbol.Symbol{SymbolEXPR, SymbolEoi}, func(params ...symbol.Value) symbol.Value {
		return nil
	})

	return parser
}
