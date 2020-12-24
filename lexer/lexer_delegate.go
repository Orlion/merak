package lexer

type LexerDelegate struct {
	lexer   Lexer
	current Token
}

func NewLexerDelegate(lexer Lexer) *LexerDelegate {
	return &LexerDelegate{
		lexer: lexer,
	}
}

func (ld *LexerDelegate) Next() (token Token, err error) {
	token, err = ld.lexer.Next()
	if err == nil {
		ld.current = token
	}
	return
}

func (ld *LexerDelegate) Current() Token {
	if ld.current == nil {
		panic("no current token")
	}

	return ld.current
}
