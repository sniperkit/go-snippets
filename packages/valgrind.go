package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Valgrind struct{}

func (valgrind Valgrind) Name() string {
	return "valgrind"
}

func (valgrind Valgrind) URL(version string) string {
	return fmt.Sprintf("http://sourceware.org/pub/valgrind/valgrind-%s.tar.bz2", version)
}

func (valgrind Valgrind) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(valgrind),
		Args: []string{
			"--enable-only64bit",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (valgrind Valgrind) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (valgrind Valgrind) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}

