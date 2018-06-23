package packages
// requires pexpect (python module)

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type RR struct{}

func (rr RR) Name() string {
	return "rr"
}

func (rr RR) URL(version string) string {
	return fmt.Sprintf("https://github.com/mozilla/rr/archive/%s.tar.gz", version)
}

func (rr RR) Build(config gogurt.Config) error {
	buildDir := filepath.Join(config.BuildDir(rr), "build")
	os.Mkdir(buildDir, 0755)

	cmake := gogurt.CMakeCmd{
		Prefix: config.InstallDir(rr),
		SourceDir: config.BuildDir(rr),
		BuildDir: buildDir,
		CacheEntries: map[string]string{
			"disable32bit": "ON",
		},
		PkgConfigPaths: []string{
			config.PkgConfigLibDir(Zlib{}),
		},
	}.Cmd()

	if err := cmake.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Dir: buildDir,
	}.Cmd()
	return make.Run()
}

func (rr RR) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
		Dir: filepath.Join(config.BuildDir(rr), "build"),
	}.Cmd()
	return makeInstall.Run()
}

func (rr RR) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Zlib{},
	}
}
