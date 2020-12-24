package symbol

type Symbol interface {
	IsTerminals() bool
	ToString() string
}
