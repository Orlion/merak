# Merak
Merak is an LR(1) parser library for Go written in Go, But it is not a code generator.  
Merak 是一个用Go编写的Go LR(1) parser 库，但不是一个代码生成器   

# Install
```
go get github.com/Orlion/merak
```

# Getting Started
1. New A Parser
```
parser := merak.NewParser()
```
2. Register Production
```
parser.RegProduction(...)
parser.RegProduction(...)
...
```
3. Parse Input
```
r, err := parser.Parse(SymbolGOAL, SymbolEoi, lexer)
```
# Example
[calculator](https://github.com/Orlion/merak/tree/main/example/calculator)