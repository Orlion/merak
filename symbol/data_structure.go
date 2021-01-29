package symbol

type ZSet struct {
	list []Symbol
	keys map[Symbol]struct{}
}

func NewSymbolZSet() *ZSet {
	return &ZSet{
		keys: make(map[Symbol]struct{}),
	}
}

func (s *ZSet) Add(symbol Symbol) bool {
	if _, exists := s.keys[symbol]; !exists {
		s.list = append(s.list, symbol)
		s.keys[symbol] = struct{}{}

		return true
	}

	return false
}

func (s *ZSet) AddAll(symbolZSet *ZSet) {
	for _, symbol := range symbolZSet.list {
		if _, exists := s.keys[symbol]; !exists {
			s.list = append(s.list, symbol)
			s.keys[symbol] = struct{}{}
		}
	}
}

func (s *ZSet) Exists(symbol Symbol) bool {
	_, exists := s.keys[symbol]
	return exists
}

func (s *ZSet) List() []Symbol {
	return s.list
}
