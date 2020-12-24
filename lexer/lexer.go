package lexer

type Lexer interface {
	Next() (Token, error)
}
