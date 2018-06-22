package packages

// TODO: Include aspell-en dictionary, either here or as a separate package.
import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Aspell struct{}

func (aspell Aspell) Name() string {
	return "aspell"
}

func (aspell Aspell) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/aspell/aspell-%s.tar.gz", version)
}

func (aspell Aspell) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(aspell),
		Args: []string{
			"--enable-static",
			"--disable-shared",
		},
		CFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
		},
		CxxFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Ncurses{}),
		},
		Libs: []string{
			"-ltinfow",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (aspell Aspell) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (aspell Aspell) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
	}
}
