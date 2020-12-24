package merak

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Orlion/merak/ast"
	"github.com/Orlion/merak/lexer"
	"github.com/Orlion/merak/symbol"
)

type TokenType int

const (
	TokenEoi        TokenType = 0
	TokenAdd        TokenType = 1
	TokenSub                  = 2
	TokenNumber               = 3
	TokenBreak                = 4
	TokenSemicolon            = 5
	TokenIdentifier           = 6
)

type Token struct {
	Value string
	T     TokenType
}

func (t Token) ToSymbol() symbol.Symbol {
	switch t.T {
	case TokenEoi:
		return SymbolEoi
	case TokenAdd:
		return SymbolAdd
	case TokenSub:
		return SymbolSub
	case TokenNumber:
		return SymbolNumber
	case TokenBreak:
		return SymbolBreak
	case TokenSemicolon:
		return SymbolSemicolon
	case TokenIdentifier:
		return SymbolIdentifier
	default:
		panic(fmt.Sprintf("ToSymbol, token.Value: %s", t.Value))
	}
}

func (t Token) ToString() string {
	return t.Value
}

type Lexer struct {
	tokens []*Token
	pos    int
}

func (l *Lexer) Next() (lexer.Token, error) {
	if l.pos > len(l.tokens) {
		return nil, errors.New("eoi err")
	} else if l.pos == len(l.tokens) {
		return &Token{Value: "", T: TokenEoi}, nil
	} else {
		l.pos = l.pos + 1
		return l.tokens[l.pos-1], nil
	}
}

type Symbol int

func (s Symbol) IsTerminals() bool {
	switch s {
	case SymbolStmt:
		return false
	case SymbolExpr:
		return false
	case SymbolNumber:
		return true
	case SymbolAdd:
		return true
	case SymbolSub:
		return true
	case SymbolEoi:
		return true
	case SymbolBreakStmt:
		return false
	case SymbolBreak:
		return true
	case SymbolIdentifierOpt:
		return false
	case SymbolIdentifier:
		return true
	case SymbolSemicolon:
		return true
	default:
		panic("IsTerminals")
	}
}

func (s Symbol) ToString() string {
	switch s {
	case SymbolStmt:
		return "stmt"
	case SymbolExpr:
		return "expr"
	case SymbolNumber:
		return "number"
	case SymbolAdd:
		return "+"
	case SymbolSub:
		return "-"
	case SymbolEoi:
		return "eoi"
	case SymbolBreak:
		return "break"
	case SymbolBreakStmt:
		return "break_stmt"
	case SymbolIdentifier:
		return "identifier"
	case SymbolIdentifierOpt:
		return "identifer_opt"
	case SymbolSemicolon:
		return ";"
	default:
		panic("ToString")
	}
}

const (
	SymbolExpr Symbol = iota + 1
	SymbolStmt
	SymbolNumber
	SymbolAdd
	SymbolSub
	SymbolEoi
	SymbolBreak
	SymbolBreakStmt
	SymbolIdentifierOpt
	SymbolIdentifier
	SymbolSemicolon
)

/*
stmt -> expr
expr -> NUMBER
	 |  expr + NUMBER
	 |  expr - NUMBER
*/

type Stmt struct {
	Expr *Expr
}

type Expr struct {
	Number string
	IsAdd  bool
	IsSub  bool
	Expr   *Expr
}

/*
break_stmt -> BREAK identifier_opt SEMICOLON
identifier_opt -> identifier
		 |  //empty
*/

type BreakStmt struct {
	IdentifierOpt *IdentifierOpt
}

type IdentifierOpt struct {
	identifier string
}

func TestParse(t *testing.T) {
	myParser := NewParser()
	myParser.RegisterProduction(SymbolStmt, []symbol.Symbol{SymbolExpr}, false, func(args []interface{}) ast.Node {
		return &Stmt{
			Expr: args[0].(*Expr),
		}
	})

	myParser.RegisterProduction(SymbolExpr, []symbol.Symbol{SymbolNumber}, false, func(args []interface{}) ast.Node {
		return &Expr{
			Number: args[0].(*Token).Value,
		}
	})

	myParser.RegisterProduction(SymbolExpr, []symbol.Symbol{SymbolExpr, SymbolAdd, SymbolNumber}, false, func(args []interface{}) ast.Node {
		return &Expr{
			Expr:   args[0].(*Expr),
			IsAdd:  true,
			Number: args[2].(*Token).Value,
		}
	})

	myParser.RegisterProduction(SymbolExpr, []symbol.Symbol{SymbolExpr, SymbolSub, SymbolNumber}, false, func(args []interface{}) ast.Node {
		return &Expr{
			Expr:   args[0].(*Expr),
			IsSub:  true,
			Number: args[2].(*Token).Value,
		}
	})

	myLexer := &Lexer{
		tokens: []*Token{&Token{Value: "1", T: TokenNumber}, &Token{Value: "+", T: TokenAdd}, &Token{Value: "2", T: TokenNumber},
			&Token{Value: "-", T: TokenSub}, &Token{Value: "3", T: TokenNumber}, &Token{Value: "eoi", T: TokenEoi}},
	}

	myAst, err := myParser.Build(SymbolStmt, SymbolEoi).SetLexer(myLexer).Parse()
	if err != nil {
		t.Fatal(err)
	}

	if stmt, ok := myAst.(*Stmt); !ok {
		t.Fatal("myAst convert to stmt failed")
	} else {
		if stmt.Expr.Expr.Number != "2" {
			t.Fatal("stmt.Expr.Expr.Number != '2'")
		}
	}

	myParser1 := NewParser()
	myParser1.RegisterProduction(SymbolBreakStmt, []symbol.Symbol{SymbolBreak, SymbolIdentifierOpt, SymbolSemicolon}, false, func(args []interface{}) ast.Node {
		return &BreakStmt{
			IdentifierOpt: args[1].(*IdentifierOpt),
		}
	})

	myParser1.RegisterProduction(SymbolIdentifierOpt, []symbol.Symbol{SymbolIdentifier}, true, func(args []interface{}) ast.Node {
		return &IdentifierOpt{
			identifier: args[0].(*Token).Value,
		}
	})

	myParser1.RegisterProduction(SymbolIdentifierOpt, []symbol.Symbol{}, true, func(args []interface{}) ast.Node {
		return &IdentifierOpt{}
	})

	myLexer1 := &Lexer{
		tokens: []*Token{&Token{Value: "break", T: TokenBreak}, &Token{Value: ";", T: TokenSemicolon}},
	}

	myAst1, err := myParser1.Build(SymbolBreakStmt, SymbolEoi).SetLexer(myLexer1).Parse()
	if err != nil {
		t.Fatal(err)
	}

	if breakStmt, ok := myAst1.(*BreakStmt); !ok {
		t.Fatal("myAst1 convert to breakStmt failed")
	} else {
		if breakStmt.IdentifierOpt.identifier != "" {
			t.Fatal("stmt.Expr.Expr.Number != ''")
		}
	}
}
