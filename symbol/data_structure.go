package symbol

type Set struct {
	elems map[Symbol]struct{}
}

func NewSymbolSet() *Set {
	return &Set{
		elems: make(map[Symbol]struct{}),
	}
}

func (s *Set) Add(symbol Symbol) bool {
	if _, exists := s.elems[symbol]; !exists {
		s.elems[symbol] = struct{}{}
		return true
	}

	return false
}

func (s *Set) AddAll(anotherSet *Set) (num int) {
	for symbol := range anotherSet.elems {
		if _, exists := s.elems[symbol]; !exists {
			s.elems[symbol] = struct{}{}
			num++
		}
	}

	return num
}

func (s *Set) Exists(symbol Symbol) bool {
	_, exists := s.elems[symbol]
	return exists
}

func (s *Set) Elems() map[Symbol]struct{} {
	return s.elems
}
