// Package cli wraps the standard flag package for a cleaner command line interface.
package cli

import (
	"flag"
	"time"
)

// A CommandUnit is a CLI command with name, description and flags.
type CommandUnit struct {
	f           func()
	name        string
	description string
	flagSet     *flag.FlagSet
}

// Command adds a new command to the CLI.
func Command(name string, f func(), description string) *CommandUnit {
	cmd := &CommandUnit{
		name:        name,
		description: description,
		f:           f,
		flagSet:     flag.NewFlagSet(name, flag.ExitOnError),
	}
	cmd.flagSet.Usage = cmd.usage // [command] [subcommand] -h
	commands = append(commands, cmd)
	return cmd
}

func (c *CommandUnit) usage() {
	printUsage(c)
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (c *CommandUnit) Bool(p *bool, name string, value bool, usage string) *CommandUnit {
	c.flagSet.BoolVar(p, name, value, usage)
	return c
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
// The flag accepts a value acceptable to time.ParseDuration.
func (c *CommandUnit) Duration(p *time.Duration, name string, value time.Duration, usage string) *CommandUnit {
	c.flagSet.DurationVar(p, name, value, usage)
	return c
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func (c *CommandUnit) Float64(p *float64, name string, value float64, usage string) *CommandUnit {
	c.flagSet.Float64Var(p, name, value, usage)
	return c
}

// Int defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func (c *CommandUnit) Int(p *int, name string, value int, usage string) *CommandUnit {
	c.flagSet.IntVar(p, name, value, usage)
	return c
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func (c *CommandUnit) Int64(p *int64, name string, value int64, usage string) *CommandUnit {
	c.flagSet.Int64Var(p, name, value, usage)
	return c
}

// String defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (c *CommandUnit) String(p *string, name string, value string, usage string) *CommandUnit {
	c.flagSet.StringVar(p, name, value, usage)
	return c
}

// Uint defines a uint flag with specified name, default value, and usage string.
// The argument p points to a uint variable in which to store the value of the flag.
func (c *CommandUnit) Uint(p *uint, name string, value uint, usage string) *CommandUnit {
	c.flagSet.UintVar(p, name, value, usage)
	return c
}

// Uint64 defines a uint64 flag with specified name, default value, and usage string.
// The argument p points to a uint64 variable in which to store the value of the flag.
func (c *CommandUnit) Uint64(p *uint64, name string, value uint64, usage string) *CommandUnit {
	c.flagSet.Uint64Var(p, name, value, usage)
	return c
}
