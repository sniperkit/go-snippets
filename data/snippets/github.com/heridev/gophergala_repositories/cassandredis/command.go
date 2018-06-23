package cassandredis

import (
	"errors"
	"fmt"
	"strings"
)

//go:generate stringer -type=CommandType
type CommandType uint

const (
	CommandLPUSH CommandType = iota + 1
	CommandLRANGE
	CommandUnknown
)

type Command struct {
	Name string
	Type CommandType
	Args []interface{}
}

func (c *Command) String() string {
	return fmt.Sprintf("type: %s args: %v", c.Type, c.Args)
}

func newCommandFromRespArray(array *respArrayValue) (*Command, error) {
	if len(array.values) < 1 {
		return nil, errors.New("no command defined")
	}

	commandName, ok := array.values[0].(*respBulkStringValue)
	if !ok {
		return nil, errors.New("should have received a bulk string value")
	}

	cmd := &Command{}
	cmd.Name = string(commandName.value)
	cmd.Type = commandTypeFromString(cmd.Name)
	cmd.Args = array.values[1:]

	return cmd, nil
}

func commandTypeFromString(s string) CommandType {
	switch strings.ToUpper(s) {
	case "LPUSH":
		return CommandLPUSH
	case "LRANGE":
		return CommandLRANGE
	default:
		return CommandUnknown
	}
}
