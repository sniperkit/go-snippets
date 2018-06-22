package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type MPC struct{}

func (mpc MPC) Name() string {
	return "mpc"
}

func (mpc MPC) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/mpc/mpc-%s.tar.gz", version)
}

func (mpc MPC) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(mpc),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--with-gmp=" + config.InstallDir(GMP{}),
			"--with-mpfr=" + config.InstallDir(MPFR{}),
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

func (mpc MPC) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (mpc MPC) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		GMP{},
		MPFR{},
	}
}
