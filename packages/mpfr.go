package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type MPFR struct{}

func (mpfr MPFR) Name() string {
	return "mpfr"
}

func (mpfr MPFR) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/mpfr/mpfr-%s.tar.gz", version)
}

func (mpfr MPFR) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(mpfr),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--with-gmp=" + config.InstallDir(GMP{}),
			"--enable-thread-safe",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (mpfr MPFR) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (mpfr MPFR) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		GMP{},
	}
}
