package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type PkgConfig struct{}

func (pkgconfig PkgConfig) Name() string {
	return "pkg-config"
}

func (pkgconfig PkgConfig) URL(version string) string {
	return fmt.Sprintf("https://pkgconfig.freedesktop.org/releases/pkg-config-%s.tar.gz", version)
}

func (pkgconfig PkgConfig) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(pkgconfig),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--with-internal-glib",
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

func (pkgconfig PkgConfig) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			"install",
		},
	}.Cmd()
	return makeInstall.Run()
}

func (pkgconfig PkgConfig) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
