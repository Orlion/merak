package symbol

type Symbol interface {
	IsTerminal() bool
	ToString() string
}
