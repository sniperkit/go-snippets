package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type AutoMake struct{}

func (automake AutoMake) Name() string {
	return "automake"
}

func (automake AutoMake) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/automake/automake-%s.tar.gz", version)
}

func (automake AutoMake) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(automake),
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (automake AutoMake) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
		Jobs: config.NumCores,
	}.Cmd()
	return makeInstall.Run()
}

func (automake AutoMake) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
