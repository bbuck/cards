package ansi

import (
	"fmt"
	"strings"
)

const (
	escape = "\033"
	start  = escape + "["
	end    = "m"

	xtermForeground = "38;5;"
	xtermBackground = "48;5;"

	reset   = start + "0" + end
	resetFG = start + "39" + end
	resetBG = start + "49" + end

	EraseLine = start + "2K"
	CursorUp  = start + "1A"
)

type Color uint8

const (
	ColorBlue   = 21
	ColorGreen  = 40
	ColorRed    = 196
	ColorOrange = 214
	ColorBlack  = 245
	ColorWhite  = 252
)

// Reset will add a full attribute reset to the end of the string.
func Reset(str string) string {
	return fmt.Sprintf("%s%s", str, reset)
}

// ResetFG adds a reset foreground to the end of the string.
func ResetFG(str string) string {
	return fmt.Sprintf("%s%s", str, resetFG)
}

// ResetBG adds a reset background to the end of the string.
func ResetBG(str string) string {
	return fmt.Sprintf("%s%s", str, resetBG)
}

// FGLine adds the ANSI color code to set the foreground color and includes a
// reset code at the end of the string.
func FGLine(str string, color Color) string {
	return fmt.Sprintf("%s%s", FG(str, color), resetFG)
}

// FG adds the ANSI color code to set the foreground color.
func FG(str string, color Color) string {
	return fmt.Sprintf("%s%s%d%s%s", start, xtermForeground, color, end, str)
}

// BGLine adds the ANSI color code to set the background color and includes a
// reset code at the end of the string.
func BGLine(str string, color Color) string {
	return fmt.Sprintf("%s%s", BG(str, color), resetBG)
}

// BG adds the ANSI color code to set the background color.
func BG(str string, color Color) string {
	return fmt.Sprintf("%s%s%d%s%s", start, xtermBackground, color, end, str)
}

func fgCode(color Color) string {
	return fmt.Sprintf("%s%s%d%s", start, xtermForeground, color, end)
}

func bgCode(color Color) string {
	return fmt.Sprintf("%s%s%d%s", start, xtermBackground, color, end)
}

var colors = map[string]Color{
	"r": ColorRed,
	"l": ColorBlue,
	"b": ColorBlack,
	"o": ColorOrange,
	"w": ColorWhite,
	"g": ColorGreen,
}

func PopulateTemplateData(data map[string]any) {
	data["reset"] = reset
	data["resetFG"] = resetFG
	data["resetBG"] = resetBG
	for name, color := range colors {
		data[name] = fgCode(color)
		data[strings.ToUpper(name)] = bgCode(color)
	}
}
