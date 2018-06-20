// Package color provide some methods to color or add attributes to a string which will
// be print to terminal.
//
// Struct Color could add many attributes to the string, but it is verbose for adding
// single attribute. In this case, you could use Dye or other Color Helpers.
package color

import (
	"fmt"
	"strconv"
	"strings"
)

// Color defines a custom color object which is defined by SGR parameters.
type Color struct {
	params []Attribute
}

// Attribute defines a single SGR Code
type Attribute int

const (
	escape      = "\x1b"
	formatStr   = "\x1b[%sm"
	unformatStr = "\x1b[0m"
)

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

// New returns a newly created color object.
func New(value ...Attribute) *Color {
	c := &Color{params: make([]Attribute, 0)}
	c.Add(value...)
	return c
}

// Add is used to chain SGR parameters. Use as many as parameters to combine
// and create custom color objects. Example: Add(color.FgRed, color.Underline).
func (c *Color) Add(value ...Attribute) *Color {
	c.params = append(c.params, value...)
	return c
}

// sequence returns a formated SGR sequence to be plugged into a "\x1b[...m"
// an example output might be: "1;36" -> bold cyan
func (c *Color) sequence() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

// Wrap wraps the s string with the colors attributes. The string is ready to
// be printed.
func (c *Color) Wrap(s string) string {
	return c.format() + s + c.unformat()
}

func (c *Color) format() string {
	return fmt.Sprintf(formatStr, c.sequence())
}

func (c *Color) unformat() string {
	return unformatStr
}

// Dye set one attribute for string, if attribute for your string is not complex,
// you should use this function.
//
// Dye(a1, s) and Dye(a2, Dye(a1, s)) is both valid and the effort of Dye(a1, Dye(a2,
// s)) is the same as that of Dye(a2, Dye(a1, s))
func Dye(a Attribute, s string) string {
	return fmt.Sprintf(formatStr, strconv.Itoa(int(a))) + s + unformatStr
}

// Black is an convenient helper function to dye string with black foreground.
func Black(s string) string { return Dye(FgBlack, s) }

// Red is an convenient helper function to dye string with red foreground.
func Red(s string) string { return Dye(FgRed, s) }

// Cyan is an convenient helper function to dye string with cyan foreground.
func Cyan(s string) string { return Dye(FgCyan, s) }

// White is an convenient helper function to dye string with white foreground.
func White(s string) string { return Dye(FgWhite, s) }

// Magenta is an convenient helper function to dye string with magenta foreground.
func Magenta(s string) string { return Dye(FgMagenta, s) }

// Blue is an convenient helper function to dye string with blue foreground.
func Blue(s string) string { return Dye(FgBlue, s) }

// Yellow is an convenient helper function to dye string with yellow foreground.
func Yellow(s string) string { return Dye(FgYellow, s) }

// Green is an convenient helper function to dye string with green foreground.
func Green(s string) string { return Dye(FgGreen, s) }
