package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type GetText struct{}

func (gettext GetText) Name() string {
	return "gettext"
}

func (gettext GetText) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/gettext/gettext-%s.tar.xz", version)
}

func (gettext GetText) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(gettext),
		Args: []string{
			"--enable-static",
			"--disable-shared",
			"--with-included-glib",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (gettext GetText) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (gettext GetText) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}

