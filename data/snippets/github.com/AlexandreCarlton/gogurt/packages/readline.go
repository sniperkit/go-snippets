package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type ReadLine struct{}

func (readline ReadLine) Name() string {
	return "readline"
}

func (readline ReadLine) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/readline/readline-%s.tar.gz", version)
}

func (readline ReadLine) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(readline),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--enable-multibyte",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (readline ReadLine) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (readline ReadLine) Dependencies() []gogurt.Package {
	return []gogurt.Package{

	}
}
