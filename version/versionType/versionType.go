package versiontype

import (
	"errors"
	"strings"
)

const (
	UNVERSIONED VersionType = iota
	NUMERIC
	SEMANTIC
)

type VersionType int

func (vt VersionType) IsNumeric() bool {
	return vt == NUMERIC
}

func (vt VersionType) IsSemantic() bool {
	return vt == SEMANTIC
}

func Parse(versionTypeName string) (VersionType, error) {
	switch strings.ToUpper(versionTypeName) {
	case "UNVERSIONED":
		return UNVERSIONED, nil
	case "NUMERIC":
		return NUMERIC, nil
	case "SEMANTIC":
		return SEMANTIC, nil
	}
	return -1, errors.New("Unknown version type: " + versionTypeName)
}
