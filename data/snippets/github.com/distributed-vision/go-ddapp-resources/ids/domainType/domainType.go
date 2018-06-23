package domaintype

import (
	"errors"
	"strings"

	"github.com/distributed-vision/go-resources/ids"
)

const (
	SCOPE ids.DomainType = iota
	IDENTITY
	SIGNATURE
	SEQUENCE
)

func Parse(domainTypeName string) (ids.DomainType, error) {
	switch strings.ToUpper(domainTypeName) {
	case "SCOPE":
		return SCOPE, nil
	case "INDENTITY":
		return IDENTITY, nil
	case "SIGNATURE":
		return SIGNATURE, nil
	case "SEQUENCE":
		return SEQUENCE, nil
	}

	return -1, errors.New("Unknown domain type: " + domainTypeName)
}
