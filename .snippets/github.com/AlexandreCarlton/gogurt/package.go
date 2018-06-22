package gogurt

// Package ...
type Package interface {
	Name() string
	URL(version string) string

	// These execute in the checked out directory.
	// Should these take in a config object? Would make things much easier and flexible.
	Build(config Config) error
	Install(config Config) error
	// Clean() int

	Dependencies() []Package
}
