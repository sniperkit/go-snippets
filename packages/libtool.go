package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type LibTool struct{}

func (libtool LibTool) Name() string {
	return "libtool"
}

func (libtool LibTool) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/libtool/libtool-%s.tar.gz", version)
}

func (libtool LibTool) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(libtool),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--enable-ltdl-install",
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

func (libtool LibTool) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (libtool LibTool) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		// TODO: Requires help2man
	}
}
