package packages

// TODO: Migrate to meson (recommended Python 3 build system).
// ./configure support will soon be deprecated.

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type SSHFS struct{}

func (sshfs SSHFS) Name() string {
	return "sshfs"
}

func (sshfs SSHFS) URL(version string) string {
	return fmt.Sprintf("https://github.com/libfuse/sshfs/releases/download/sshfs-%s/sshfs-%s.tar.gz", version, version)
}

func (sshfs SSHFS) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(sshfs),
		PkgConfigPaths: []string{
			config.PkgConfigLibDir(FUSE{}),
			config.PkgConfigLibDir(GLib{}),
			config.PkgConfigLibDir(Pcre{}),
		},
		LdFlags: []string{
			"-ldl",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (sshfs SSHFS) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (sshfs SSHFS) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		FUSE{},
		GLib{},
	}
}
