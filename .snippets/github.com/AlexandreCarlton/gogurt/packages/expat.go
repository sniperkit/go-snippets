package packages

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/alexandrecarlton/gogurt"
)

type Expat struct{}

func (expat Expat) Name() string {
	return "expat"
}

func (expat Expat) URL(version string) string {
	underscoredVersion := strings.Replace(version, ".", "_", 2)
	return fmt.Sprintf("https://github.com/libexpat/libexpat/archive/R_%s.tar.gz", underscoredVersion)
}

func (expat Expat) Build(config gogurt.Config) error {

	expatDir := filepath.Join(config.BuildDir(expat), "expat")
	// CMake will probably be the way to go in future,
	// but until then we'll stick with configure.
	buildconf := exec.Command("./buildconf.sh")
	buildconf.Dir = expatDir
	if err := buildconf.Run(); err != nil {
		return err
	}

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(expat),
		Args: []string{
			"--disable-shared",
			"--enable-static",
		},
		Dir: expatDir,
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Dir: expatDir,
	}.Cmd()
	return make.Run()
}

func (expat Expat) Install(config gogurt.Config) error {
	expatDir := filepath.Join(config.BuildDir(expat), "expat")
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			// We don't use 'install' as this will try to generate documentation using docbook.
			"installlib",
		},
		Dir: expatDir,
	}.Cmd()
	return makeInstall.Run()
}

func (expat Expat) Dependencies() []gogurt.Package {
	return []gogurt.Package{
	}
}
