package game

import (
	"fmt"
	"strings"

	"bbuck.dev/cards/ansi"
)

func Run(game Game) {
	game.Init()

	for {
		screen := game.Display()
		lines := strings.Split(screen, "\n")
		fmt.Println(screen)
		fmt.Print(game.Prompt())

		var line string
		fmt.Scanf("%s", &line)
		line = strings.TrimSpace(line)

		game.HandleInput(line)

		for range len(lines) + 1 {
			fmt.Print(ansi.EraseLine, ansi.CursorUp)
		}
	}
}
