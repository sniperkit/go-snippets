package schemevisibility

import (
	"errors"
	"strings"

	"github.com/distributed-vision/go-resources/ids"
)

const (
	UNTYPED ids.SchemeVisibility = iota
	LOCAL
	PRIVATE
	PUBLIC
)

func Parse(visibility string) (ids.SchemeVisibility, error) {
	switch strings.ToUpper(visibility) {
	case "UNTYPED":
		return UNTYPED, nil
	case "LOCAL":
		return LOCAL, nil
	case "PRIVATE":
		return PRIVATE, nil
	case "PUBLIC":
		return PUBLIC, nil
	default:
		return -1, errors.New("Unknown domain visibility: " + visibility)
	}
}
