package packages

// NOTE: 2.2.2 requires texi2html, which has been deprecated by texinfo.

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Stow struct{}

func (stow Stow) Name() string {
	return "stow"
}

func (stow Stow) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/stow/stow-%s.tar.gz", version)
}

func (stow Stow) Build(config gogurt.Config) error {

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(stow),
		Paths: []string{
			config.BinDir(TexInfo{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}

	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Paths: []string{
			config.BinDir(TexInfo{}),
		},
	}.Cmd()
	return make.Run()
}

func (stow Stow) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
		Paths: []string{
			config.BinDir(TexInfo{}),
		},
	}.Cmd()
	return makeInstall.Run()
}

func (stow Stow) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		TexInfo{},
	}
}
