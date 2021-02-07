package lexer

import "github.com/Orlion/merak/symbol"

type Token interface {
	ToString() string
	Filename() string
	Line() int
	Col() int
	ToSymbol() symbol.Value
}
