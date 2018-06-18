package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gizak/gopp/plugin"
)

func Run() error {
	subcmd := ""
	if len(os.Args) > 1 {
		subcmd = os.Args[1]
	}

	plugin.Init()

	if subcmd == "" {
		b, err := exec.Command("gopp", "help").CombinedOutput()
		fmt.Print(string(b))
		return err
	}

	// plugin added command
	if runner, ok := plugin.SubcmdRunners[subcmd]; ok {
		if err := runner.RunSubcmd(os.Args[1:]); err != nil {
			return err
		}
	} else { // go command
		cmd := exec.Command("go", os.Args[1:]...)
		cmd.Stdin = os.Stdin
		out, _ := cmd.CombinedOutput()
		combpipe := bytes.NewReader(out)

		var lastr io.Reader
		var nextw io.Writer
		var b bytes.Buffer

		for i, orw := range plugin.OutputRewriters {
			lastr = &b
			nextw = &b

			if i == 0 {
				lastr = combpipe
			}

			if i == len(plugin.OutputRewriters)-1 {
				nextw = os.Stdout
			}

			if err := orw.RewriteOutput(lastr, nextw); err != nil {
				return err
			}
		}
	}

	return nil
}
