package utils

import (
	"fmt"
)

const (
	colorRed = uint8(iota + 91)
	colorGreen
)

// Red get red color
func Red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorRed, s)
}

// Green get green color
func Green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorGreen, s)
}
