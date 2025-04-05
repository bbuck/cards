package cards

import "fmt"

// Suit represents the specific suit the card belongs too.
type Suit int8

// Define the four suits of a standard deck, Spades, Hearts, Clubs, and
// Diamonds
const (
	SuitSpades Suit = iota
	SuitHearts
	SuitDiamonds
	SuitClubs
)

// suitUnicodeMap provides a mapping between a suit and a unicode character
// representing that suit.
var (
	suitUnicodeMap = map[Suit]string{
		SuitSpades:   "\u2660",
		SuitHearts:   "\u2665",
		SuitDiamonds: "\u2666",
		SuitClubs:    "\u2663",
	}

	suitAsciiMap = map[Suit]string{
		SuitSpades:   "<3-",
		SuitHearts:   "<3 ",
		SuitDiamonds: "<> ",
		SuitClubs:    "c3-",
	}

	suitLetterMap = map[Suit]string{
		SuitSpades:   "S",
		SuitHearts:   "H",
		SuitDiamonds: "D",
		SuitClubs:    "C",
	}
)

// String returns the name of the suit.
func (s Suit) String() string {
	switch s {
	case SuitSpades:
		return "spades"
	case SuitHearts:
		return "hearts"
	case SuitDiamonds:
		return "diamonds"
	case SuitClubs:
		return "clubs"
	default:
		return fmt.Sprintf("InvalidSuit{%d}", s)
	}
}

// Display returns a representation of the given type for the suit.
func (s Suit) Display() string {
	if str, ok := suitUnicodeMap[s]; ok {
		return str
	}

	return "?"
}

// Color returns the color this suit belongs to.
func (s Suit) Color() Color {
	if s == SuitSpades || s == SuitClubs {
		return ColorBlack
	}

	return ColorRed
}
