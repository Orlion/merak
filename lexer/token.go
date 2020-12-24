package lexer

import "github.com/Orlion/merak/symbol"

type Token interface {
	ToSymbol() symbol.Symbol
	ToString() string
}
