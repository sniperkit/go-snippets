package goplugin

import (
	"github.com/gizak/gopp/plugin/ptype"
	"io"
	"os/exec"
	//_ "github.com/maruel/panicparse"
)

var info = ptype.PluginInfo{
	Name: "parsepanic",
	Uri:  "github.com/gizak/gopp/plugin/parsepanic",
}

var PLUGIN = Plugin{}

type Plugin struct{}

func (Plugin) DeclarePlugin() ptype.PluginInfo {
	return info
}

// TODO: use programmatic call
func (Plugin) RewriteOutput(in io.Reader, out io.Writer) error {
	cmd := exec.Command("panicparse")
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out

	return cmd.Run()
}
