package registry

import (
	"errors"
)

var (
	// ErrNotARegistryAddress is returned when a registry address is not valid.
	ErrNotARegistryAddress = errors.New("not a registry address")

	// ErrNotAValidRepositoryName is returned when a repository name is not valid.
	ErrNotAValidRepositoryName = errors.New("not a valid repository name")
)
