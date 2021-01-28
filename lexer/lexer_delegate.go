package lexer

type LexerDelegator struct {
	lexer   Lexer
	current Token
}

func NewLexerDelegator(lexer Lexer) *LexerDelegator {
	return &LexerDelegator{
		lexer: lexer,
	}
}

func (ld *LexerDelegator) Next() (token Token, err error) {
	token, err = ld.lexer.Next()
	if err == nil {
		ld.current = token
	}
	return
}

func (ld *LexerDelegator) Current() Token {
	return ld.current
}
