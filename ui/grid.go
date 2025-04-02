package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CardGrid represents a grid of playing cards arranged in rows and columns
type CardGrid struct {
	// Grid contents - we use PlayingCard objects directly instead of cards.Card
	Cards     [][]PlayingCard

	// Spacing between cards
	HorizontalSpacing int
	VerticalSpacing   int

	// Offset for staggered/overlapping displays (useful for solitaire)
	RowOffset int
	ColOffset int
}

// NewCardGrid creates a new empty grid
func NewCardGrid() *CardGrid {
	return &CardGrid{
		Cards:             [][]PlayingCard{},
		HorizontalSpacing: 1,
		VerticalSpacing:   0,
		RowOffset:         0,
		ColOffset:         0,
	}
}

// AddCards adds a row of cards to the grid
func (g *CardGrid) AddCards(cards []PlayingCard) {
	g.Cards = append(g.Cards, cards)
}

// Render returns the rendered grid of cards
func (g *CardGrid) Render() string {
	if len(g.Cards) == 0 {
		return ""
	}

	// If using offsets, we need to render card by card and manually position them
	if g.RowOffset > 0 || g.ColOffset > 0 {
		return g.renderWithOffsets()
	}

	// Otherwise, render row by row, handling multiline card renders properly
	var result strings.Builder

	for rowIdx, row := range g.Cards {
		if rowIdx > 0 {
			// Add vertical spacing between rows
			result.WriteString(strings.Repeat("\n", g.VerticalSpacing+1))
		}

		if len(row) == 0 {
			continue
		}

		// Render each card first
		var renderedCards [][]string
		maxLines := 0

		for _, card := range row {
			// Split the rendered card into lines
			renderedCard := card.Render()
			cardLines := strings.Split(renderedCard, "\n")
			renderedCards = append(renderedCards, cardLines)

			// Keep track of the maximum number of lines
			if len(cardLines) > maxLines {
				maxLines = len(cardLines)
			}
		}

		// Join each line of the cards horizontally
		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			if lineIdx > 0 {
				result.WriteString("\n")
			}

			for cardIdx, cardLines := range renderedCards {
				if cardIdx > 0 {
					// Add horizontal spacing between cards
					result.WriteString(strings.Repeat(" ", g.HorizontalSpacing))
				}

				// Get the current line of the card, or an empty line if it doesn't exist
				line := ""
				if lineIdx < len(cardLines) {
					line = cardLines[lineIdx]
				} else {
					// For cards with fewer lines, add empty space
					line = strings.Repeat(" ", len(cardLines[0]))
				}

				result.WriteString(line)
			}
		}
	}

	return result.String()
}

// renderWithOffsets handles rendering cards with row/column offsets
// Useful for staggered displays like in solitaire
func (g *CardGrid) renderWithOffsets() string {
	// First, determine the total width and height needed
	totalWidth := 0
	totalHeight := 0
	cardWidth := 0
	cardHeight := 0

	// Get the dimensions of a card (using a blank style to get dimensions)
	sampleCard := PlayingCardStyle.Render("")
	lines := strings.Split(sampleCard, "\n")
	cardHeight = len(lines)
	if cardHeight > 0 {
		cardWidth = lipgloss.Width(lines[0])
	}

	// Calculate total dimensions
	for rowIdx, row := range g.Cards {
		rowWidth := cardWidth + (len(row)-1)*g.ColOffset
		if rowWidth > totalWidth {
			totalWidth = rowWidth
		}
		rowHeight := cardHeight + rowIdx*g.RowOffset
		if rowHeight > totalHeight {
			totalHeight = rowHeight
		}
	}

	// Create a 2D grid of strings to represent the rendered output
	grid := make([][]string, totalHeight)
	for i := range grid {
		grid[i] = make([]string, totalWidth)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	// Place each card in the grid
	for rowIdx, row := range g.Cards {
		for colIdx, card := range row {
			renderedCard := card.Render()

			// Position of the card's top-left corner
			startRow := rowIdx * g.RowOffset
			startCol := colIdx * g.ColOffset

			// Place the card in the grid
			cardLines := strings.Split(renderedCard, "\n")
			for i, line := range cardLines {
				if startRow+i < len(grid) {
					// Insert line character by character
					for j, ch := range line {
						if startCol+j < len(grid[startRow+i]) {
							grid[startRow+i][startCol+j] = string(ch)
						}
					}
				}
			}
		}
	}

	// Convert the grid to a string
	var result strings.Builder
	for _, row := range grid {
		result.WriteString(strings.Join(row, ""))
		result.WriteString("\n")
	}

	return result.String()
}

// RenderCompact renders a grid with overlapping cards (useful for card stacks)
func (g *CardGrid) RenderCompact(overlap int) string {
	grid := *g
	grid.ColOffset = overlap
	return grid.renderWithOffsets()
}
