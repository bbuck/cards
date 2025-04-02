package decks

import (
	"errors"
	"iter"
	"math/rand/v2"
	"slices"
	"sort"

	"bbuck.dev/cards/cards"
)

var EmptyDeckErr = errors.New("no more cards in deck")

// Deck represents a set of cards, typically "stacked" and in some kind of
// order.
type Deck struct {
	cards []*cards.Card
}

// New constructs a brand new deck of 52 standard playing cards.
func New() *Deck {
	deck := &Deck{cards: make([]*cards.Card, 0, 52)}
	for _, suit := range cards.Suits {
		for _, value := range cards.Values {
			card := &cards.Card{Suit: suit, Value: value}
			deck.cards = append(deck.cards, card)
		}
	}

	return deck
}

// Empty creates a deck with no cards in it.
func Empty() *Deck {
	return &Deck{cards: make([]*cards.Card, 0)}
}

// Shuffle the cards in the deck into a random order.
func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Draw a single card from the deck.
func (d *Deck) Draw() (*cards.Card, error) {
	if drawn, err := d.DrawN(1); err != nil {
		return nil, err
	} else {
		return drawn[0], nil
	}
}

// Len returns the number of cards remaining in the deck.
func (d *Deck) Len() int {
	return len(d.cards)
}

// AddBottom adds the cards to the bottom of the deck.
func (d *Deck) AddBottom(newCards []*cards.Card) {
	d.cards = append(d.cards, newCards...)
}

// Bottom returns the card at the bottom of the deck.
func (d *Deck) Bottom() (*cards.Card, error) {
	if len(d.cards) == 0 {
		return nil, EmptyDeckErr
	}

	return d.cards[len(d.cards)-1], nil
}

// Draw n cards from the deck and return them as a slice.
func (d *Deck) DrawN(n int) ([]*cards.Card, error) {
	if n < 1 {
		return nil, errors.New("cannot draw fewer than 1 card")
	}

	if len(d.cards) <= n {
		return nil, EmptyDeckErr
	}

	drawn := d.cards[:n]
	d.cards = d.cards[n:]

	return drawn, nil
}

// Purge removes all cards from the deck that match the predicate.
func (d *Deck) Purge(predicate func(*cards.Card) bool) {
	var remaining []*cards.Card
	for _, card := range d.cards {
		if !predicate(card) {
			remaining = append(remaining, card)
		}
	}
	d.cards = remaining
}

// Sort allows sorting the cards in the deck in a specific order using the
// provided function to determine if the first card should occur before (higher
// in the deck) than the second card.
func (d *Deck) Sort(less func(a, b *cards.Card) bool) {
	sort.Slice(d.cards, func(i, j int) bool {
		return less(d.cards[i], d.cards[j])
	})
}

// Iter returns a new iterator over the cards remaining in the deck.
func (d *Deck) Iter() iter.Seq2[int, *cards.Card] {
	return slices.All(d.cards)
}
