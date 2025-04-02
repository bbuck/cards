package main

import (
	"bbuck.dev/cards/game"
	"bbuck.dev/cards/scoundrel"
)

func main() {
	g := scoundrel.New()
	game.Run(g)
}
