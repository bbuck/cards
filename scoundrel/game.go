package scoundrel

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"bbuck.dev/cards/ansi"
	"bbuck.dev/cards/cards"
	"bbuck.dev/cards/decks"
)

var gameTemplate = template.Must(
	template.New("game").Parse(`
{{ .o }}Health: {{ .r }}{{ .hearts }} {{ .o }}({{ .health }}){{ .reset }}

{{ .o }}     {{ .select0 }} {{ .select1 }} {{ .select2 }} {{ .select3 }}{{ .reset }}
{{ .cardBack }}  {{ .room0 }} {{ .room1 }} {{ .room2 }} {{ .room3 }}  {{ .discard }}
     {{ .weapon }}
     {{ .lastMonster }}

`),
)

// Scoundrel game state.
type Scoundrel struct {
	deck        *decks.Deck
	room        []*cards.Card
	health      int
	weapon      *cards.Card
	lastMonster *cards.Card
	discard     *decks.Deck
}

// New builds a new game of Scoundrel.
func New() *Scoundrel {
	return &Scoundrel{
		deck:        nil,
		room:        nil,
		health:      20,
		weapon:      nil,
		lastMonster: nil,
		discard:     nil,
	}
}

// Init sets up the game of Scoundrel for play.
func (game *Scoundrel) Init() {
	game.deck = decks.New()
	game.room = nil
	game.health = 20
	game.weapon = nil
	game.lastMonster = nil
	game.discard = decks.Empty()

	game.deck.Purge(func(card *cards.Card) bool {
		return card.IsFaceCard() && card.Color() == cards.ColorRed
	})
	game.deck.Shuffle()
	game.room, _ = game.deck.DrawN(4)
}

const heartString = "♥ "

// templateData builds a map to use for rendering the game UI.
func (game *Scoundrel) templateData() map[string]any {
	data := make(map[string]any)
	ansi.PopulateTemplateData(data)

	data["cardBack"] = cards.CardBack()

	hearts := ansi.FGLine(strings.Repeat(heartString, game.health), ansi.ColorRed)
	data["hearts"] = hearts
	data["health"] = game.health
	data["weapon"] = "   "
	data["lastMonster"] = "   "
	data["room0"] = "   "
	data["select0"] = "   "
	data["select1"] = "   "
	data["select2"] = "   "
	data["select3"] = "   "
	data["room1"] = "   "
	data["room2"] = "   "
	data["room3"] = "   "
	data["discard"] = "   "

	if game.room != nil {
		for i, card := range game.room {
			data[fmt.Sprintf("room%d", i)] = card.Display()
			data[fmt.Sprintf("select%d", i)] = fmt.Sprintf("[%d]", i+1)
		}
	}

	if game.weapon != nil {
		data["weapon"] = game.weapon.Display()
	}

	if game.lastMonster != nil {
		data["lastMonster"] = game.lastMonster.Display()
	}

	// bottom of the deck is top of the discard (top being face down cards and
	// bottom being face up)
	if bottom, err := game.discard.Bottom(); err == nil {
		data["discard"] = bottom.Display()
	}

	return data
}

// isGameOver determines if the game has ended. A game of scoundrel is over if
// the players health hits zero or when the deck is empty before the player
// dies.
func (game *Scoundrel) isGameOver() bool {
	return game.health <= 0 || game.deck.Len() == 0
}

// score calculates the score of the player, the score is only valid once the
// game has ended.
func (game *Scoundrel) score() int {
	if game.deck.Len() == 0 && game.health > 0 {
		return game.health
	}

	score := game.health
	for _, card := range game.deck.Iter() {
		if card.Color() == cards.ColorBlack {
			damage := int(card.Value)
			score = score - damage
		}
	}

	return score
}

// Display renders the game UI for the Scoundrel game state.
func (game *Scoundrel) Display() string {
	if game.isGameOver() {
		var (
			scoreStr = ""
			score    = game.score()
		)
		if score > 0 {
			scoreStr = ansi.FG(strconv.Itoa(score), ansi.ColorGreen)
		} else {
			scoreStr = ansi.FG(strconv.Itoa(score), ansi.ColorRed)
		}

		return fmt.Sprintf("%s %s\n", ansi.FGLine("Score:", ansi.ColorOrange), ansi.Reset(scoreStr))
	}

	buf := new(bytes.Buffer)
	if err := gameTemplate.Execute(buf, game.templateData()); err != nil {
		return err.Error()
	}

	return buf.String()
}

// Prompt returns the current input prompt for the Scoundrel game.
func (game *Scoundrel) Prompt() string {
	if game.isGameOver() {
		return ""
	}

	return ansi.FGLine("Select Card [1-4] >> ", ansi.ColorOrange)
}

// HandleInput processes the input for the Scoundrel game.
func (game *Scoundrel) HandleInput(input string) error {
	if game.isGameOver() {
		return nil
	}

	index, err := strconv.Atoi(input)
	if err != nil {
		return err
	}

	if index < 1 || index > 4 {
		return errors.New("invalid room card")
	}

	index = index - 1
	if len(game.room) <= index {
		return errors.New("no card in the room at that index")
	}

	card := game.room[index]

	before := game.room[0:index]
	after := game.room[index+1:]
	var newRoom []*cards.Card
	newRoom = append(newRoom, before...)
	game.room = append(newRoom, after...)

	if len(game.room) <= 1 {
		for range 3 {
			if card, err := game.deck.Draw(); err == nil {
				game.room = append(game.room, card)
			} else {
				break
			}
		}
	}

	if card.Suit == cards.SuitDiamonds {
		if game.weapon != nil {
			if game.lastMonster != nil {
				game.discard.AddBottom([]*cards.Card{game.lastMonster})
			}
			game.discard.AddBottom([]*cards.Card{game.weapon})
		}
		game.weapon = card
		game.lastMonster = nil

		return nil
	}

	if card.Suit == cards.SuitHearts {
		game.health = game.health + int(card.Value)
		game.health = min(game.health, 20)
		game.discard.AddBottom([]*cards.Card{card})

		return nil
	}

	damage := int(card.Value)
	canUseWeapon := false
	if game.weapon != nil {
		canUseWeapon = true
		if game.lastMonster != nil {
			canUseWeapon = damage <= int(game.lastMonster.Value)
		}

		if canUseWeapon {
			weaponDamage := int(game.weapon.Value)
			damage = damage - weaponDamage
		}
	}
	if damage > 1 {
		game.health = game.health - damage
	}

	if game.weapon != nil && canUseWeapon {
		if game.lastMonster != nil {
			game.discard.AddBottom([]*cards.Card{game.lastMonster})
		}
		game.lastMonster = card
	} else {
		game.discard.AddBottom([]*cards.Card{card})
	}

	return nil
}
