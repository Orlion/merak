package main

type NonterminalSymbol string

func (s NonterminalSymbol) Name() string {
	return string(s)
}

type TerminalSymbol string

func (s TerminalSymbol) Name() string {
	return string(s)
}

type Symbol interface {
	Name() string
}

type Terminal interface {
	Symbol
}

type Nonterminal interface {
	Symbol
}

const SymbolAdd TerminalSymbol = "+"

func T(t Terminal, nt Nonterminal) {

}

func main() {
	T(SymbolAdd, SymbolAdd)
}
