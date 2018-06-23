package packages

import (
	"fmt"
	"os/exec"
	"github.com/alexandrecarlton/gogurt"
)

type UniversalCTags struct{}

func (ctags UniversalCTags) Name() string {
	return "universal-ctags"
}

func (ctags UniversalCTags) URL(version string) string {
	return fmt.Sprintf("https://github.com/universal-ctags/ctags/archive/%s.tar.gz", version)
}

func (ctags UniversalCTags) Build(config gogurt.Config) error {

	bootstrap := exec.Command("./autogen.sh")
	if err := bootstrap.Run(); err != nil {
		return err
	}

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(ctags),
		Paths: []string{
			config.BinDir(AutoConf{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}

	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (ctags UniversalCTags) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (ctags UniversalCTags) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		AutoConf{},
	}
}

