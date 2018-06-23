package goplugin

import (
	"github.com/gizak/colorgo"
	"io"
	"os"
)

var PLUGIN = aplugin{}

type aplugin struct{}

func (ap aplugin) RewriteOutput(in io.Reader, out io.Writer) error {

	if len(os.Args) > 1 && (os.Args[1] == "build" || os.Args[1] == "test") {
		return colorgo.Colorize(in, out)

	}
	_, err := io.Copy(out, in)
	return err
}
