package args

import (
	"bytes"
	"fmt"
	proto "github.com/chremoas/chremoas/proto"
	"golang.org/x/net/context"
)

type Args struct {
	cmdName string
	argMap  map[string]*Command
	argList []string
}

type Command struct {
	Funcptr func(ctx context.Context, request *proto.ExecRequest) string
	Help    string
}

func NewArg(cmdName string) *Args {
	a := &Args{}
	a.argMap = make(map[string]*Command)
	a.cmdName = cmdName
	return a
}

func (a *Args) Add(name string, command *Command) {
	a.argList = append(a.argList, name)
	a.argMap[name] = command
}

func (a Args) Exec(ctx context.Context, req *proto.ExecRequest, rsp *proto.ExecResponse) error {
	var response string

	if len(req.Args) == 1 || req.Args[1] == "help" {
		response = a.help()
	} else {
		f, ok := a.argMap[req.Args[1]]
		if ok {
			response = f.Funcptr(ctx, req)
		} else {
			return fmt.Errorf("not a valid subcommand: %s", req.Args[1])
		}
	}

	rsp.Result = []byte(response)
	return nil
}

func (a Args) help() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Usage: !%s <subcommand> <arguments>\n", a.cmdName))
	buffer.WriteString("\nSubcommands:\n")

	for cmd := range a.argList {
		if a.argMap[a.argList[cmd]].Help != "" {
			buffer.WriteString(fmt.Sprintf("\t%s: %s\n",
				a.argList[cmd],
				a.argMap[a.argList[cmd]].Help,
			))
		}
	}

	return fmt.Sprintf("```%s```", buffer.String())
}
