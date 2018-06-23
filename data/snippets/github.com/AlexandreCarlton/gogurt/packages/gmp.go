package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type GMP struct{}

func (gmp GMP) Name() string {
	return "gmp"
}

func (gmp GMP) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/gmp/gmp-%s.tar.xz", version)
}

func (gmp GMP) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(gmp),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--enable-cxx",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (gmp GMP) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (gmp GMP) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
