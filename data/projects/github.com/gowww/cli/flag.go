package cli

import "flag"

// A Flag represents the state of a flag.
type Flag struct {
	flag.Flag
}

// NameWithDefValue returns the name of the flag with its default value.
func (f *Flag) NameWithDefValue() string {
	if f.DefValue == "" {
		return f.Name
	}
	return f.Name + "=" + f.DefValue
}
