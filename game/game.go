package game

type Game interface {
	Init()
	Prompt() string
	HandleInput(input string) error
	Display() string
}
