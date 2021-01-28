package symbol

type Symbol interface {
	Name() string
	IsTerminal() bool
}
