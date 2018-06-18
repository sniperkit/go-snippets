// Package bashconf offers parsing of very simple bash-like configuration files
//
// Please note: this isn't a good fit for most Go-only apps!  We've built this
// to handle situations where it's a major benefit to have our configs in a
// format that works in bash (not to mention Python) natively.
package bashconf

import (
	"io/ioutil"
	"os"
	"strings"
)

// Config is a simple type for holding key/value pairs
type Config struct {
	raw          map[string]string
	allowEnvVars bool
	envPrefix    string
}

// New returns a config instance for use
func New() *Config {
	return &Config{raw: make(map[string]string)}
}

// EnvironmentOverrides turns on or off checking the environment on calls to
// c.Get(key)
func (c *Config) EnvironmentOverrides(allow bool) {
	c.allowEnvVars = allow
}

// EnvironmentPrefix sets the prefix required for environment overrides to
// work.  If EnvironmentOverrides is off, calling this will turn it on.
func (c *Config) EnvironmentPrefix(prefix string) {
	c.allowEnvVars = true
	c.envPrefix = prefix
}

// Get looks up the string in its raw datastore and returns it.  If the value
// exists in an environment variable, and this config was set up to overlay
// environment variables (this is disabled by default), that is returned instead.
func (c *Config) Get(key string) string {
	var val = c.raw[key]

	// Only override with env if (a) c has been configured to allow this, and (b)
	// the environment definitely has the given key
	if c.allowEnvVars {
		var envval, ok = os.LookupEnv(c.envPrefix + key)
		if ok {
			val = envval
		}
	}

	return val
}

// ParseFile reads config from a file and adds to the Config map.  If the file
// can't be read or parsed, an error will be returned.
//
// This function is built to read bash-like environmental variables in order to
// maximize cross-language configurability.  Very minimal processing occurs
// here, leaving that to the caller.  i.e., "FOO=1" will yield a key of "FOO"
// and a value of "1" since bash has no types.  This supports only the simplest
// of bash variable assignment: no arrays, no substitutions of other variables,
// just very basic key/value pairs.
func (c *Config) ParseFile(filename string) error {
	var content, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	c.ParseString(string(content))
	return nil
}

// ParseString reads each line in s and attempts to parse a key/value pair.
// Any line starting with a comment is ignored, and lines that don't split into
// a key/value are ignored.
func (c *Config) ParseString(s string) {
	var lines = strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line[0] == '#' {
			continue
		}

		var kvparts = strings.SplitN(line, "=", 2)
		if len(kvparts) != 2 {
			continue
		}
		var key, val = kvparts[0], kvparts[1]

		if key == "" {
			continue
		}

		// Remove surrounding quotes if any exist, but only one level of quotes
		if val[0] == '"' {
			val = strings.Trim(val, `"`)
		} else if val[0] == '\'' {
			val = strings.Trim(val, `'`)
		}
		c.raw[key] = val
	}
}
