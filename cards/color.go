package cards

// Color represents the color of the card, either black or red.
type Color int8

const (
	// ColorRed is the color of hearts and diamonds
	ColorRed = iota

	// ColorBlack is the color of clubs and spades
	ColorBlack
)
