package plugin

import (
	"bytes"
	"fmt"
	"github.com/gizak/gopp/plugin/ptype"
	"os/exec"
	"text/template"
)

type cmdHelp struct {
	tmpl    string
	Subcmds map[string]ptype.SubcmdRunner
}

var pCmdHelp = cmdHelp{tmpl: `Go is a tool for managing Go source code.

This is the gopp version of go command with some additional features.
Vanilla go commands are expecting to work as usual.

Usage:

	gopp command [arguments]

The original go commands are:

	build       compile packages and dependencies
	clean       remove object files
	doc         show documentation for package or symbol
	env         print Go environment information
	fix         run go tool fix on packages
	fmt         run gofmt on package sources
	generate    generate Go files by processing source
	get         download and install packages and dependencies
	install     compile and install packages and dependencies
	list        list packages
	plugin      manage plugins
	rm          remove a installed pacakge
	run         compile and run Go program
	test        test packages
	tool        run specified go tool
	version     print Go version
	vet         run go tool vet on packages

The additional gopp commands are:

{{ range $i, $p := .Subcmds -}}
{{ printf "\t%-12s%s\n" $p.Subcmd $p.Descrip }}
{{- end }}
Use "go help [command]" for more information about a command.

Additional help topics:

	c           calling between Go and C
	buildmode   description of build modes
	filetype    file types
	gopath      GOPATH environment variable
	environment environment variables
	importpath  import path syntax
	packages    description of package lists
	testflag    description of testing flags
	testfunc    description of testing functions

Use "gopp help [topic]" for more information about that topic.
`}

func (p cmdHelp) Subcmd() string {
	return "help"
}

func (p cmdHelp) Descrip() string {
	return "provide commands/plugins help info"
}

func (p cmdHelp) Usage() string {
	p.Subcmds = SubcmdRunners
	t, err := template.New("help").Parse(p.tmpl)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, p)

	return buf.String()
}

func (p cmdHelp) RunSubcmd(args []string) error {
	if len(args) == 1 {
		fmt.Print(p.Usage())
		return nil
	}

	cmd := args[1]

	if p, ok := SubcmdRunners[cmd]; ok {
		fmt.Print(p.Usage())
		return nil
	}

	// original go cmds will fall through
	b, err := exec.Command("go", args...).CombinedOutput()
	fmt.Print(string(b))
	return err
}
