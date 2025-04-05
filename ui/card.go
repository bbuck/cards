package ui

import (
	"fmt"
	"strings"

	"bbuck.dev/cards/cards"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	redColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	blackColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))

	layouts = map[cards.Value]string{
		cards.ValueAce:   "V     \n       \n   S   \n       \n     V",
		cards.Value2:     "V     \n   S   \n       \n   S   \n     V",
		cards.Value3:     "V     \n   S   \n   S   \n   S   \n     V",
		cards.Value4:     "V     \n S   S \n       \n S   S \n     V",
		cards.Value5:     "V     \n S   S \n   S   \n S   S \n     V",
		cards.Value6:     "V     \n S   S \n S   S \n S   S \n     V",
		cards.Value7:     "V     \n S   S \n S S S \n S   S \n     V",
		cards.Value8:     "V     \n S S S \n S   S \n S S S \n     V",
		cards.Value9:     "V     \n S S S \n S S S \n S S S \n     V",
		cards.Value10:    "V S   \n S S S \n S   S \n S S S \n   S V",
		cards.ValueJack:  "V     \n S     \n  V   \n     S \n     V",
		cards.ValueKing:  "V     \n S     \n  V   \n     S \n     V",
		cards.ValueQueen: "V     \n S     \n  V   \n     S \n     V",
	}
)

type Card struct {
	Suit        cards.Suit
	Value       cards.Value
	BorderColor lipgloss.Color
}

func (c Card) Render() string {
	style := redColor
	if c.Suit.Color() == cards.ColorBlack {
		style = blackColor
	}

	layout := layouts[c.Value]
	card := strings.Replace(layout, "V", fmt.Sprintf("%-2s", strings.TrimSpace(c.Value.String())), 1)
	card = strings.ReplaceAll(card, "V", c.Value.String())
	card = strings.ReplaceAll(card, "S", c.Suit.Display())

	return borderStyle.BorderForeground(c.BorderColor).Render(
		style.Render(card),
	)
}
