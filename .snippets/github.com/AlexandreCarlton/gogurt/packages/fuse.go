package packages

// TODO: Migrate to meson (recommended Python 3 build system).
// ./configure support will soon be deprecated.

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type FUSE struct{}

func (fuse FUSE) Name() string {
	return "fuse"
}

func (fuse FUSE) URL(version string) string {
	return fmt.Sprintf("https://github.com/libfuse/libfuse/releases/download/fuse-%s/fuse-%s.tar.gz", version, version)
}

func (fuse FUSE) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(fuse),
		Args: []string{
			"--disable-shared",
			"--enable-static",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (fuse FUSE) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (fuse FUSE) Dependencies() []gogurt.Package {
	return []gogurt.Package{

	}
}
