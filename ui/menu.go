package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	menuTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA500")).
		Margin(1, 0, 2, 0)

	menuItemStyle        = lipgloss.NewStyle().Padding(0, 4)
	selectedMenuItemStyle = menuItemStyle.Copy().
				Background(lipgloss.Color("#FFA500")).
				Foreground(lipgloss.Color("#000000"))
)

// GameOption represents a selectable game option
type GameOption struct {
	Title       string
	Description string
	ID          string
}

// MenuModel represents the game selection menu
type MenuModel struct {
	options  []GameOption
	cursor   int
	selected string
}

// NewMenuModel creates a new menu model with the available games
func NewMenuModel() MenuModel {
	return MenuModel{
		options: []GameOption{
			{
				Title:       "Scoundrel",
				Description: "A roguelike card game",
				ID:          "scoundrel",
			},
			{
				Title:       "Card Viewer",
				Description: "View and test playing card components",
				ID:          "cards",
			},
			// Add more games here as they are implemented
		},
		cursor:   0,
		selected: "",
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.options[m.cursor].ID
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	var s string

	s += menuTitleStyle.Render("Card Games")
	s += "\nSelect a game to play:\n\n"

	for i, option := range m.options {
		item := fmt.Sprintf("%s - %s", option.Title, option.Description)

		if i == m.cursor {
			s += selectedMenuItemStyle.Render("> " + item)
		} else {
			s += menuItemStyle.Render("  " + item)
		}
		s += "\n"
	}

	s += "\n(↑/↓ to navigate, enter to select, q to quit)\n"

	return s
}

// Selected returns the ID of the selected game, or an empty string if none is selected
func (m MenuModel) Selected() string {
	return m.selected
}
