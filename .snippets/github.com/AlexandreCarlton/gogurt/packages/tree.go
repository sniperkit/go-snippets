package packages

import (
	"fmt"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type Tree struct{}

func (tree Tree) Name() string {
	return "tree"
}

func (tree Tree) URL(version string) string {
	return fmt.Sprintf("http://mama.indstate.edu/users/ice/tree/src/tree-%s.tgz", version)
}

func (tree Tree) Build(config gogurt.Config) error {
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (tree Tree) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			"install",
			"prefix=" + config.InstallDir(tree),
			"MANDIR=" + filepath.Join(config.InstallDir(tree), "share", "man", "man1"),
		},
	}.Cmd()
	return makeInstall.Run()
}

func (tree Tree) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
