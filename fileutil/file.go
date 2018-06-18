package fileutil

import (
	"os"

	isatty "github.com/mattn/go-isatty"
)

func IsTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}
