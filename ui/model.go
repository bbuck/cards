package ui

import (
	"fmt"
	"strings"

	"bbuck.dev/cards/scoundrel"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFA500"))

	healthStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Width(8)

	selectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFA500")).
				Padding(0, 1).
				Width(8)

	weaponStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FF00")).
			Padding(0, 1).
			Width(8)

	monsterStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF0000")).
			Padding(0, 1).
			Width(8)

	deckStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Width(8)
)

// Card represents a single card in the UI
type Card struct {
	Number  int
	Content string
	Style   lipgloss.Style
}

// Render returns the rendered card
func (c Card) Render() string {
	return c.Style.Render(fmt.Sprintf("[%d]\n%s", c.Number, c.Content))
}

// CardRow represents a row of cards
type CardRow struct {
	Cards []Card
}

// Render returns the rendered row of cards
func (r CardRow) Render() string {
	var sb strings.Builder
	for i, card := range r.Cards {
		sb.WriteString(card.Render())
		if i < len(r.Cards)-1 {
			sb.WriteString("  ")
		}
	}
	return sb.String()
}

// GameState represents the current state of the game UI
type GameState struct {
	Title    string
	Health   int
	Deck     string
	Room     CardRow
	Discard  string
	Weapon   *Card
	Monster  *Card
	Error    error
	Selected int
}

// Render returns the rendered game state
func (s GameState) Render() string {
	var sb strings.Builder

	// Title
	sb.WriteString(titleStyle.Render(s.Title))
	sb.WriteString("\n\n")

	// Health
	hearts := strings.Repeat("♥ ", s.Health)
	sb.WriteString(fmt.Sprintf("%s (%d)\n\n",
		healthStyle.Render(hearts),
		s.Health))

	// Deck and Room cards
	sb.WriteString(deckStyle.Render(s.Deck))
	sb.WriteString("  ")
	sb.WriteString(s.Room.Render())
	sb.WriteString("  ")
	sb.WriteString(deckStyle.Render(s.Discard))
	sb.WriteString("\n\n")

	// Weapon and Monster
	if s.Weapon != nil {
		sb.WriteString(s.Weapon.Render())
		sb.WriteString("  ")
	}
	if s.Monster != nil {
		sb.WriteString(s.Monster.Render())
	}
	sb.WriteString("\n")

	// Error message
	if s.Error != nil {
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(s.Error.Error()))
	}

	return sb.String()
}

type model struct {
	game     *scoundrel.Scoundrel
	selected int
	err      error
}

func NewModel() model {
	return model{
		game:     scoundrel.New(),
		selected: -1,
	}
}

func (m model) Init() tea.Cmd {
	m.game.Init()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4":
			index := int(msg.String()[0] - '1')
			if index >= 0 && index < len(m.game.Room()) {
				m.selected = index
				if err := m.game.HandleInput(fmt.Sprintf("%d", index+1)); err != nil {
					m.err = err
				} else {
					m.err = nil
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.game.IsGameOver() {
		score := m.game.Score()
		scoreColor := "#00FF00"
		if score <= 0 {
			scoreColor = "#FF0000"
		}
		return fmt.Sprintf("%s %s\n",
			titleStyle.Render("Score:"),
			lipgloss.NewStyle().Foreground(lipgloss.Color(scoreColor)).Render(fmt.Sprintf("%d", score)))
	}

	// Build room cards
	roomCards := make([]Card, len(m.game.Room()))
	for i, card := range m.game.Room() {
		style := cardStyle
		if i == m.selected {
			style = selectedCardStyle
		}
		roomCards[i] = Card{
			Number:  i + 1,
			Content: card.Display(),
			Style:   style,
		}
	}

	// Build weapon card if exists
	var weaponCard *Card
	if weapon := m.game.Weapon(); weapon != nil {
		weaponCard = &Card{
			Content: weapon.Display(),
			Style:   weaponStyle,
		}
	}

	// Build monster card if exists
	var monsterCard *Card
	if monster := m.game.LastMonster(); monster != nil {
		monsterCard = &Card{
			Content: monster.Display(),
			Style:   monsterStyle,
		}
	}

	// Create game state
	state := GameState{
		Title:    "Scoundrel",
		Health:   m.game.Health(),
		Deck:     "Deck",
		Room:     CardRow{Cards: roomCards},
		Discard:  "Discard",
		Weapon:   weaponCard,
		Monster:  monsterCard,
		Error:    m.err,
		Selected: m.selected,
	}

	return state.Render()
}
