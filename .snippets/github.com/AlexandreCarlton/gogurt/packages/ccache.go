package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type CCache struct{}

func (ccache CCache) Name() string {
	return "ccache"
}

func (ccache CCache) URL(version string) string {
	return fmt.Sprintf("https://www.samba.org/ftp/ccache/ccache-%s.tar.xz", version)
}

func (ccache CCache) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(ccache),
		CFlags: []string{
			"-I" + config.IncludeDir(Zlib{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Zlib{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (ccache CCache) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (ccache CCache) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Zlib{},
	}
}
