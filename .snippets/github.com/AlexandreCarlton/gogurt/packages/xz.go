package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type XZ struct{}

func (xz XZ) Name() string {
	return "xz"
}

func (xz XZ) URL(version string) string {
	return fmt.Sprintf("https://tukaani.org/xz/xz-%s.tar.gz", version)
}

func (xz XZ) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(xz),
		Args: []string{
			"--disable-shared",
			"--enable-static",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Paths: []string{},
	}.Cmd()
	return make.Run()
}

func (xz XZ) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (xz XZ) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
