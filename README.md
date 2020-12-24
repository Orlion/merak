# Merak
Merak is an LR(1) parser library for Go written in Go, But it is not a code generator.  
Merak 是一个用Go编写的Go LR(1) parser 库，但不是一个代码生成器   

# 安装
```
go get github.com/Orlion/merak
```

# Getting Started
0. 首先我们应该对编译原理有所了解，并清楚生成式的概念。我们的目标是生成一个简单计算器的parser，并返回解析输入字符的AST，它的生成式如下：
```
stmt -> expr
expr -> NUMBER
	 |  expr + NUMBER
	 |  expr - NUMBER
```
1. Merak中存在一个symbol.Symbol接口，它代表生成式中的一个符号，包括终结符与非终结符。我们首先需要提供一个symbo.Symbol的实现，并提前定义好生成式中的所有符号，如下：
```
type Symbol int

// 标识该符号是否是终结符
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
)

```
2. 然后我们需要提供一个lexer.Token接口的实现，如下
```
type TokenType int

const (
	TokenEoi    TokenType = 0
	TokenAdd    TokenType = 1
	TokenSub              = 2
	TokenNumber           = 3
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
	default:
		panic("ToSymbol")
	}
}

func (t Token) ToString() string {
	return t.Value
}
```
注意：Merak要求在输入结束处返回一个Eoi Token，所以需要提前定义一个TokenEoi

3. 接下来我们提供一个lexer.Lexer接口的实现，lexer.Lexer需要你实现Next()方法，该方法返回下一个lexer.Token，出错时返回error，如下：
```
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
```
4. 准备工作完成！开始使用Merak
```
// 创建merak.Parser实例
myParser := NewParser()
// 注册生成式
// stmt -> expr
myParser.RegisterProduction(SymbolStmt, []symbol.Symbol{SymbolExpr}, false, func(args []interface{}) ast.Node {
    // 识别到此生成式时执行下面的操作
    return &Stmt{
        Expr: args[0].(*Expr),
    }
})
// expr -> NUMBER
myParser.RegisterProduction(SymbolExpr, []symbol.Symbol{SymbolNumber}, false, func(args []interface{}) ast.Node {
    return &Expr{
        Number: args[0].(*Token).Value,
    }
})
// expr -> expr + Number
myParser.RegisterProduction(SymbolExpr, []symbol.Symbol{SymbolExpr, SymbolAdd, SymbolNumber}, false, func(args []interface{}) ast.Node {
    return &Expr{
        Expr:   args[0].(*Expr),
        IsAdd:  true,
        Number: args[2].(*Token).Value,
    }
})
// expr -> expr - Number
myParser.RegisterProduction(SymbolExpr, []symbol.Symbol{SymbolExpr, SymbolSub, SymbolNumber}, false, func(args []interface{}) ast.Node {
    return &Expr{
        Expr:   args[0].(*Expr),
        IsSub:  true,
        Number: args[2].(*Token).Value,
    }
})

// 创建Lexer实例
// 这里我们假定输入是1+2-3
myLexer := &Lexer{
    tokens: []*Token{&Token{Value: "1", T: TokenNumber}, &Token{Value: "+", T: TokenAdd}, &Token{Value: "2", T: TokenNumber},
        &Token{Value: "-", T: TokenSub}, &Token{Value: "3", T: TokenNumber}, &Token{Value: "eoi", T: TokenEoi}},
}

// 调用Build，注入目标符号和结束符号，调用SetLexer注入lexer，最后调用Parse
myAst, err := myParser.Build(SymbolStmt, SymbolEoi).SetLexer(myLexer).Parse()
```