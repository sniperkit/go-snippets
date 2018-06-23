// Utilities for messing with time.
package timeutil

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/stringutil"
)

const rxExtendedDurations = `(?i)((?P<number>[\d\.]+)(?P<suffix>[^\d]+))`
const durationMaxMatches = 16

// Provides an API-compatible version of time.ParseDuration that accepts additional
// formats for parsing durations:
//
//   1y, 3year, 5years: Expands to (n*24*365) hours
//   1w, 3weeks, 5wks:  Expands to (n*24*7) hours
//   1d, 1day, 5days:   Expands to (n*24) hours
//
func ParseDuration(in string) (time.Duration, error) {
	var i int
	var totalHours int

	in = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}

		return r
	}, in)

	for {
		i++

		if match := rxutil.Match(rxExtendedDurations, in); match != nil {
			if num, err := stringutil.ConvertToInteger(match.Group(`number`)); err == nil {
				var hours int

				switch strings.ToLower(match.Group(`suffix`)) {
				case `year`, `years`, `y`:
					hours = 8760
				case `week`, `weeks`, `wk`, `wks`, `w`:
					hours = 168
				case `day`, `days`, `d`:
					hours = 24
				case `hour`, `hours`, `hr`, `hrs`, `h`:
					hours = 1
				case `minute`, `minutes`, `min`:
					in = match.ReplaceGroup(`suffix`, `m`)
				default:
					break
				}

				if hours > 0 {
					in = match.ReplaceGroup(1, ``)
				}

				totalHours += int(num) * hours
			} else {
				return 0, fmt.Errorf("Invalid number: %v", err)
			}
		} else {
			break
		}

		if i >= durationMaxMatches {
			break
		}
	}

	if totalHours > 0 {
		in = fmt.Sprintf("%dh%s", totalHours, in)
	}

	return time.ParseDuration(in)
}
