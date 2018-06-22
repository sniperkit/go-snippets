package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Global struct{}

func (global Global) Name() string {
	return "global"
}

func (global Global) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/global/global-%s.tar.gz", version)
}

func (global Global) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(global),
		Args: []string{
			"--enable-static",
			"--disable-shared",
			"--with-universal-ctags=" + config.BinDir(UniversalCTags{}),
			"--with-ncurses=" + config.InstallDir(Ncurses{}),
			"--disable-gtagscscope", // for now.
			"--with-included-ltdl", // could use our own libtool?
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

func (global Global) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (global Global) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
		UniversalCTags{},
	}
}
