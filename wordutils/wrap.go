// Package wordutils is my very simple and naive port of one small aspect of
// the apache commons WordUtils Java library
package wordutils

import "strings"

const newline = "\n"

// Wrap takes s and wraps it at the given maximum length by inserting newlines
// ("\n").  This is intentionally not as configurable as the apache commons
// WordUtils version in order to simplify code and calling.
func Wrap(s string, wrapLen int) string {
	if wrapLen < 1 {
		wrapLen = 1
	}

	inputLen := len(s)
	offset := 0
	wrapped := ""

	for offset < inputLen {
		if s[offset] == ' ' {
			offset++
			continue
		}

		// only last line without leading spaces is left
		if inputLen-offset <= wrapLen {
			break
		}

		spaceToWrapAt := strings.LastIndex(s[:wrapLen+offset], " ")

		if spaceToWrapAt >= offset {
			// normal case
			wrapped += s[offset:spaceToWrapAt] + newline
			offset = spaceToWrapAt + 1
		} else {
			// wrap really long word one line at a time
			wrapped += s[offset:wrapLen+offset] + newline
			offset += wrapLen
		}
	}

	// Whatever is left in line is short enough to just pass through
	wrapped += s[offset:]

	return wrapped
}
