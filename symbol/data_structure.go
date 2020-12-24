package symbol

type ZSet struct {
	List []Symbol
	keys map[Symbol]struct{}
}

func NewSymbolZSet() *ZSet {
	return &ZSet{
		keys: make(map[Symbol]struct{}),
	}
}

func (s *ZSet) Add(symbol Symbol) bool {
	if _, exists := s.keys[symbol]; !exists {
		s.List = append(s.List, symbol)
		s.keys[symbol] = struct{}{}

		return true
	}

	return false
}

func (s *ZSet) AddAll(symbolZSet *ZSet) {
	for _, symbol := range symbolZSet.List {
		if _, exists := s.keys[symbol]; !exists {
			s.List = append(s.List, symbol)
			s.keys[symbol] = struct{}{}
		}
	}
}

func (s *ZSet) Exists(symbol Symbol) bool {
	_, exists := s.keys[symbol]
	return exists
}
