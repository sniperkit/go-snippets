package schemeformat

import (
	"errors"
	"strings"

	"github.com/distributed-vision/go-resources/ids"
)

const (
	FIXED ids.SchemeFormat = iota
	LV
)

func Parse(format string) (ids.SchemeFormat, error) {
	switch strings.ToUpper(format) {
	case "FIXED":
		return FIXED, nil
	case "LV":
		return LV, nil
	default:
		return -1, errors.New("Unknown domain format: " + format)
	}
}
