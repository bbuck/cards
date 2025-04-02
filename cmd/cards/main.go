package main

import (
	"fmt"
	"os"
	"strings"

	"bbuck.dev/cards/cards"
	"bbuck.dev/cards/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var gameToRun string

	// If an argument is provided, use it as the game to run
	if len(os.Args) > 1 {
		gameToRun = strings.ToLower(os.Args[1])
	} else {
		// No arguments, show the menu
		gameToRun = runGameMenu()
	}

	// If no game was selected from the menu (user quit), exit
	if gameToRun == "" {
		return
	}

	runGame(gameToRun)
}

// runGameMenu displays a menu for the user to select a game
func runGameMenu() string {
	p := tea.NewProgram(ui.NewMenuModel())
	model, err := p.Run()
	if err != nil {
		fmt.Printf("Error running menu: %v\n", err)
		os.Exit(1)
	}

	// Get the selected game from the menu model
	menuModel, ok := model.(ui.MenuModel)
	if !ok {
		fmt.Println("Could not get selected game")
		os.Exit(1)
	}

	return menuModel.Selected()
}

// runGame launches the specified game
func runGame(game string) {
	switch game {
	case "scoundrel":
		runScoundrel()
	case "cards", "deck":
		runCardTest()
	default:
		fmt.Printf("Unknown game: %s\n", game)
		os.Exit(1)
	}
}

func runScoundrel() {
	p := tea.NewProgram(ui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// runCardTest displays a test grid of all cards to demonstrate the component
func runCardTest() {
	// Create a test grid
	testGrid := ui.NewCardGrid()

	// Row 1: Face cards
	faceCardsRow := []ui.PlayingCard{
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.ValueAce}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitHearts, Value: cards.ValueKing}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitDiamonds, Value: cards.ValueQueen}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitClubs, Value: cards.ValueJack}),
	}
	testGrid.AddCards(faceCardsRow)

	// Row 2: High number cards
	highCardsRow := []ui.PlayingCard{
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.Value10}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitHearts, Value: cards.Value9}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitDiamonds, Value: cards.Value8}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitClubs, Value: cards.Value7}),
	}
	testGrid.AddCards(highCardsRow)

	// Row 3: Low number cards
	lowCardsRow := []ui.PlayingCard{
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.Value6}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitHearts, Value: cards.Value5}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitDiamonds, Value: cards.Value4}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitClubs, Value: cards.Value3}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.Value2}),
	}
	testGrid.AddCards(lowCardsRow)

	// Set spacing between rows
	testGrid.VerticalSpacing = 1
	testGrid.HorizontalSpacing = 2

	// Create a second grid for stacked/overlapping cards (like in solitaire)
	stackedGrid := ui.NewCardGrid()

	// Create a stack of cards with suits of spades
	spadesStack := []ui.PlayingCard{
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.ValueAce}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.Value2}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.Value3}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.Value4}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitSpades, Value: cards.Value5}),
	}
	stackedGrid.AddCards(spadesStack)

	// Create a stack of cards with suits of hearts
	heartsStack := []ui.PlayingCard{
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitHearts, Value: cards.ValueAce}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitHearts, Value: cards.Value2}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitHearts, Value: cards.Value3}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitHearts, Value: cards.Value4}),
		ui.NewPlayingCard(&cards.Card{Suit: cards.SuitHearts, Value: cards.Value5}),
	}
	stackedGrid.AddCards(heartsStack)

	// Create a face down row
	faceDownRow := []ui.PlayingCard{
		ui.NewFaceDownCard(),
		ui.NewFaceDownCard(),
		ui.NewFaceDownCard(),
		ui.NewFaceDownCard(),
		ui.NewFaceDownCard(),
	}
	stackedGrid.AddCards(faceDownRow)

	// Configure spacing and offset for overlapping display
	stackedGrid.VerticalSpacing = 1
	stackedGrid.RowOffset = 2
	stackedGrid.ColOffset = 4

	// Create a model to display our grids
	model := cardGridTestModel{
		regularGrid:   testGrid,
		stackedGrid:   stackedGrid,
		currentView:   "regular",
		toggleOptions: []string{"regular", "stacked", "compact"},
		toggleIndex:   0,
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// cardGridTestModel displays and allows toggling between different card grid displays
type cardGridTestModel struct {
	regularGrid   *ui.CardGrid
	stackedGrid   *ui.CardGrid
	currentView   string
	toggleOptions []string
	toggleIndex   int
}

func (m cardGridTestModel) Init() tea.Cmd {
	return nil
}

func (m cardGridTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "t", " ":
			// Toggle between different views
			m.toggleIndex = (m.toggleIndex + 1) % len(m.toggleOptions)
			m.currentView = m.toggleOptions[m.toggleIndex]
		}
	}
	return m, nil
}

func (m cardGridTestModel) View() string {
	output := "Card Grid Component Test\n\n"
	output += "Press 't' to toggle views, 'q' to quit\n\n"

	switch m.currentView {
	case "regular":
		output += "Regular Grid View:\n\n"
		output += m.regularGrid.Render()
	case "stacked":
		output += "Stacked Grid View (for games like Solitaire):\n\n"
		output += m.stackedGrid.Render()
	case "compact":
		output += "Compact Stacked View (for card piles):\n\n"
		output += m.stackedGrid.RenderCompact(2)
	}

	output += "\n\n"
	return output
}
