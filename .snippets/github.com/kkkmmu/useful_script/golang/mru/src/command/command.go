package command

import (
	"strings"
)

type Command struct {
	Delay  int    `json:"delay"`
	Mode   string `json:"mode"`
	CMD    string `json:"cli"`
	End    string
	Result string
}

func (cmd *Command) Dup() *Command {
	return &Command{
		Delay:  cmd.Delay,
		Mode:   cmd.Mode,
		CMD:    cmd.CMD,
		End:    cmd.End,
		Result: cmd.Result,
	}
}

func (cmd *Command) Validate() {
	cmd.CMD = strings.TrimSpace(cmd.CMD)
}
