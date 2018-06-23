package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Zsh struct{}

func (zsh Zsh) Name() string {
	return "zsh"
}

func (zsh Zsh) URL(version string) string {
	return fmt.Sprintf("http://www.zsh.org/pub/zsh-%s.tar.xz", version)
}

func (zsh Zsh) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(zsh),
		Args: []string{
			"--enable-cap",
			"--enable-pcre",
			"--enable-multibyte",
			"--with-term-lib=tinfow",
		},
		CFlags: []string{
			"-I" + config.IncludeDir(Pcre{}),
			"-I" + config.IncludeDir(Ncurses{}),
		},
		CppFlags: []string{
			"-I" + config.IncludeDir(Pcre{}),
			"-I" + config.IncludeDir(Ncurses{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Pcre{}),
			"-L" + config.LibDir(Ncurses{}),
		},
		Paths: []string{
			config.BinDir(Pcre{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (zsh Zsh) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (zsh Zsh) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
		Pcre{},
	}
}
