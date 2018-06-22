package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type TexInfo struct{}

func (texinfo TexInfo) Name() string {
	return "texinfo"
}

func (texinfo TexInfo) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/texinfo/texinfo-%s.tar.gz", version)
}

func (texinfo TexInfo) Build(config gogurt.Config) error {

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(texinfo),
		CppFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Ncurses{}),
		},
		Libs: []string{
			"-ltinfo",
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

func (texinfo TexInfo) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (texinfo TexInfo) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
	}
}
