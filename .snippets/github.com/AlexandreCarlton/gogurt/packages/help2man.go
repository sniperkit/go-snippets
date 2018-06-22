package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Help2Man struct{}

func (help2man Help2Man) Name() string {
	return "help2man"
}

func (help2man Help2Man) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/help2man/help2man-%s.tar.xz", version)
}

func (help2man Help2Man) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(help2man),
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (help2man Help2Man) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (help2man Help2Man) Dependencies() []gogurt.Package {
	return []gogurt.Package{
	}
}
