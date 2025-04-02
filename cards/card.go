package cards

import (
	"fmt"

	"bbuck.dev/cards/ansi"
)

func init() {
	suits := []Suit{SuitSpades, SuitHearts, SuitDiamonds, SuitClubs}
	values := []Value{Value2, Value3, Value4, Value5, Value6, Value7, Value8, Value9, Value10, ValueJack, ValueQueen, ValueKing, ValueAce}

	var codepoint int32 = 0x1f0a1
	for i, suit := range suits {
		var suitCodepoint int32 = 0x10 * int32(i)
		unicodeCardMap[suit] = make(map[Value]rune)
		for j, value := range values {
			unicodeCardMap[suit][value] = rune(codepoint + suitCodepoint + int32(j))
		}
	}
}

var (
	Suits  = []Suit{SuitSpades, SuitHearts, SuitDiamonds, SuitClubs}
	Values = []Value{Value2, Value3, Value4, Value5, Value6, Value7, Value8, Value9, Value10, ValueJack, ValueQueen, ValueKing, ValueAce}
)

var unicodeCardMap = map[Suit]map[Value]rune{}

// CardBack returns a string representation for the back of a card.
func CardBack() string {
	return ansi.FGLine(ansi.BGLine("[x]", ansi.ColorBlue), ansi.ColorWhite)
}

// Card represents a physical card in a standard deck of cards.
type Card struct {
	Suit  Suit
	Value Value
}

// Color returns the color of the card based on it's suit. Spades and Clubs are
// black cards and Hearts and Diamonds are red cards.
func (c *Card) Color() Color {
	if c.Suit == SuitSpades || c.Suit == SuitClubs {
		return ColorBlack
	}

	return ColorRed
}

// String returns a string representation of the card's suit and value.
func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.Value, c.Suit)
}

// Display returns a game representation of the card.
func (c *Card) Display() string {
	cardText := fmt.Sprintf("%s%s", c.Value.Display(), c.Suit.Display())
	if c.Color() == ColorRed {
		return ansi.ResetFG(ansi.FG(cardText, ansi.ColorRed))
	}

	return ansi.ResetFG(ansi.FG(cardText, ansi.ColorBlack))
}

// IsFaceCard determines if the card is a face card. Face cards are Jacks,
// Queens, Kings, and Aces.
func (c *Card) IsFaceCard() bool {
	return c.Value.IsFaceCard()
}
