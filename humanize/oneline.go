package humanize

import "regexp"

var re = regexp.MustCompile(`\s*\n\s*`)

// Oneline turns a multi-line string into a single line by collapsing newlines
// surrounded by any amount of whitespace (space, tab, etc.) into a single
// ASCII space.  Due to the proliferation of one-space-after-punctuation, and
// to simplify this function, exceptions aren't made to ensure two spaces after
// periods, question-marks, etc.  No, I'm not fond of this approach, either.
func Oneline(s string) string {
	return re.ReplaceAllString(s, " ")
}
