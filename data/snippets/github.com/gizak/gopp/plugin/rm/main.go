package goplugin

import (
	"fmt"
	"go/build"
	"os"
	"path"
)

type aplugin struct{}

var PLUGIN = aplugin{}

func (p aplugin) Subcmd() string {
	return "rm"
}

func (p aplugin) Descrip() string {
	return "remove installed package from GOPATH/GOROOT"
}

func (p aplugin) Usage() string {
	return "Usage:\n\tgopp rm [pacakge]\n"
}

func (p aplugin) RunSubcmd(args []string) error {
	if len(args) < 2 {
		fmt.Print(p.Usage())
		os.Exit(1)
	}

	pkgpath := args[1]
	pkg, err := build.Import(pkgpath, "", build.FindOnly)

	if err != nil {
		return err
	}

	fmt.Printf("%+v", pkg)

	// rm from src tree
	os.Remove(pkg.Dir)

	// rm obj
	os.Remove(pkg.PkgObj)

	// rm bin
	if pkg.IsCommand() {
		bin := path.Join(pkg.BinDir, path.Base(pkg.Dir))
		os.Remove(bin)
	}

	return nil
}
