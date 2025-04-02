package ui

import (
	"fmt"
	"strings"

	"bbuck.dev/cards/ansi"
	"bbuck.dev/cards/cards"
	"github.com/charmbracelet/lipgloss"
)

var (
	// PlayingCardStyle is the base style for playing cards
	PlayingCardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Width(11).
		Height(5)

	// RedCardStyle is used for hearts and diamonds
	RedCardStyle = PlayingCardStyle.Copy().
			Foreground(lipgloss.Color("#FF0000"))

	// BlackCardStyle is used for clubs and spades - using gray instead of pure black
	BlackCardStyle = PlayingCardStyle.Copy().
			Foreground(lipgloss.Color("#AAAAAA"))

	// CardBackStyle is used for facedown cards
	CardBackStyle = PlayingCardStyle.Copy().
			Background(lipgloss.Color("#0000FF")).
			Foreground(lipgloss.Color("#FFFFFF"))
)

// PlayingCard is a UI component that renders a standard playing card
type PlayingCard struct {
	Card     *cards.Card
	FaceDown bool
}

// NewPlayingCard creates a new PlayingCard component
func NewPlayingCard(card *cards.Card) PlayingCard {
	return PlayingCard{
		Card:     card,
		FaceDown: false,
	}
}

// NewFaceDownCard creates a new face down PlayingCard
func NewFaceDownCard() PlayingCard {
	return PlayingCard{
		FaceDown: true,
	}
}

// Render returns the rendered playing card
func (c PlayingCard) Render() string {
	if c.FaceDown {
		return renderCardBack()
	}

	return c.renderFaceUp()
}

func (c PlayingCard) renderFaceUp() string {
	if c.Card == nil {
		return PlayingCardStyle.Render("Empty")
	}

	style := BlackCardStyle
	if c.Card.Color() == cards.ColorRed {
		style = RedCardStyle
	}

	// Get card details
	valueDisplay := c.Card.Value.Display()
	suitDisplay := c.Card.Suit.Display()
	symbol := fmt.Sprintf("%s%s", valueDisplay, suitDisplay)
	padding := strings.Repeat(" ", 9-len(ansi.StripANSI(symbol)))

	// Build all 5 lines for the card
	var lines []string

	// Line 1: Top corner with value and suit
	lines = append(lines, symbol + padding)

	// Lines 2-4: Card pattern based on card type
	var patternLines []string
	if c.Card.IsFaceCard() {
		patternLines = c.renderFaceCardPattern(valueDisplay, suitDisplay)
	} else {
		patternLines = c.renderNumberCardPattern()
	}
	lines = append(lines, patternLines...)

	// Line 5: Bottom corner with value and suit (reversed)
	lines = append(lines, padding + symbol)

	return style.Render(strings.Join(lines, "\n"))
}

// renderFaceCardPattern returns the pattern lines for face cards (J, Q, K, A)
func (c PlayingCard) renderFaceCardPattern(value, suit string) []string {
	switch c.Card.Value {
	case cards.ValueAce:
		return []string{
			"       ",
			"   " + suit + "   ",
			"       ",
		}
	case cards.ValueKing:
		return []string{
			"       ",
			"   ♚   ",
			"       ",
		}
	case cards.ValueQueen:
		return []string{
			"       ",
			"   ♛   ",
			"       ",
		}
	case cards.ValueJack:
		return []string{
			"       ",
			"   ♝   ",
			"       ",
		}
	default:
		// Shouldn't happen, but just in case
		return []string{
			"       ",
			"   " + suit + "   ",
			"       ",
		}
	}
}

// renderNumberCardPattern returns the pattern lines for number cards (2-10)
func (c PlayingCard) renderNumberCardPattern() []string {
	suit := c.Card.Suit.Display()
	cardValue := int(c.Card.Value) // Value2 is already 2, Value3 is 3, etc.

	switch cardValue {
	case 2:
		return []string{
			"   " + suit + "   ",
			"       ",
			"   " + suit + "   ",
		}
	case 3:
		return []string{
			"   " + suit + "   ",
			"   " + suit + "   ",
			"   " + suit + "   ",
		}
	case 4:
		return []string{
			" " + suit + "   " + suit + " ",
			"       ",
			" " + suit + "   " + suit + " ",
		}
	case 5:
		return []string{
			" " + suit + "   " + suit + " ",
			"   " + suit + "   ",
			" " + suit + "   " + suit + " ",
		}
	case 6:
		return []string{
			" " + suit + "   " + suit + " ",
			" " + suit + "   " + suit + " ",
			" " + suit + "   " + suit + " ",
		}
	case 7:
		return []string{
			" " + suit + "   " + suit + " ",
			" " + suit + " " + suit + " " + suit + " ",
			" " + suit + "   " + suit + " ",
		}
	case 8:
		return []string{
			" " + suit + " " + suit + " " + suit + " ",
			" " + suit + "   " + suit + " ",
			" " + suit + " " + suit + " " + suit + " ",
		}
	case 9:
		return []string{
			" " + suit + " " + suit + " " + suit + " ",
			" " + suit + " " + suit + " " + suit + " ",
			" " + suit + " " + suit + " " + suit + " ",
		}
	case 10:
		return []string{
			" " + suit + " " + suit + " " + suit + " ",
			suit + " " + suit + " " + suit + " " + suit,
			" " + suit + " " + suit + " " + suit + " ",
		}
	default:
		return []string{
			"       ",
			"   " + suit + "   ",
			"       ",
		}
	}
}

func renderCardBack() string {
	lines := []string{
		"♠♥♣♦♠♥♣♦♠",
		"▒▒▒▒▒▒▒▒▒",
		"♠♣♥♦♠♣♥♦♠",
		"▒▒▒▒▒▒▒▒▒",
		"♠♥♣♦♠♥♣♦♠",
	}

	return CardBackStyle.Render(strings.Join(lines, "\n"))
}
