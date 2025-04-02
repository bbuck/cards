package cards

import "fmt"

// Value represents the face value of the card. The face value of a card is
// simply the numeric value of the card and then 11 for Jack, 12 for Queen,
// 13 for King, and 14 for Ace.
type Value int8

// Define all the values for a standard deck of cards.
const (
	Value2 Value = iota + 2
	Value3
	Value4
	Value5
	Value6
	Value7
	Value8
	Value9
	Value10
	ValueJack
	ValueQueen
	ValueKing
	ValueAce
)

var (
	valueAsciiMap = map[Value]string{
		Value2:     " 2",
		Value3:     " 3",
		Value4:     " 4",
		Value5:     " 5",
		Value6:     " 6",
		Value7:     " 7",
		Value8:     " 8",
		Value9:     " 9",
		Value10:    "10",
		ValueJack:  " J",
		ValueQueen: " Q",
		ValueKing:  " K",
		ValueAce:   " A",
	}
)

// String returns a string representation of a card's value.
func (v Value) String() string {
	if str, ok := valueAsciiMap[v]; ok {
		return str
	}

	return fmt.Sprintf("InvalidValue{%d}", v)
}

// Display returns a string representation of the card's value.
func (v Value) Display() string {
	if str, ok := valueAsciiMap[v]; ok {
		return str
	}

	return " ?"
}

// IsFaceCard determines if the value is between a Jack or Ace.
func (v Value) IsFaceCard() bool {
	return v >= ValueJack && v <= ValueAce
}
