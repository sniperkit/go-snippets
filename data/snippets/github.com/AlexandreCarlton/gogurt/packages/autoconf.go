package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type AutoConf struct{}

func (autoconf AutoConf) Name() string {
	return "autoconf"
}

func (autoconf AutoConf) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/autoconf/autoconf-%s.tar.gz", version)
}

func (autoconf AutoConf) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(autoconf),
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (autoconf AutoConf) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (autoconf AutoConf) Dependencies() []gogurt.Package {
	return []gogurt.Package{

	}
}
