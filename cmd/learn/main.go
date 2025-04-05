package main

import (
	"bbuck.dev/cards/cards"
	"bbuck.dev/cards/decks"
	"bbuck.dev/cards/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	prog := tea.NewProgram(newModel())
	if _, err := prog.Run(); err != nil {
		panic(err)
	}
}

type model struct {
	deck *decks.Deck
	card *cards.Card
}

func newModel() model {
	d := decks.New()
	return model{
		deck: d,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Sequence(shuffleCmd, drawCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case shuffleMsg:
		m.deck.Shuffle()

	case drawMsg:
		var err error
		if m.card, err = m.deck.Draw(); err != nil {
			return m, tea.Quit
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "d", "space", "enter":
			return m, drawCmd
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.card == nil {
		return ""
	}

	c := ui.Card{
		Suit:  m.card.Suit,
		Value: m.card.Value,
	}

	return c.Render()
}

type shuffleMsg struct{}

func shuffleCmd() tea.Msg {
	return shuffleMsg{}
}

type drawMsg struct{}

func drawCmd() tea.Msg {
	return drawMsg{}
}
