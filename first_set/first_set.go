package first_set

import "github.com/Orlion/merak/symbol"

type FirstSet struct {
	m map[symbol.Symbol]*Production
}

func NewFirstSet(m map[symbol.Symbol]*Production) *FirstSet {
	return &FirstSet{
		m: m,
	}
}

func (fs *FirstSet) Get(s symbol.Symbol) *symbol.ZSet {
	return fs.m[s].GetFirstZSet()
}
