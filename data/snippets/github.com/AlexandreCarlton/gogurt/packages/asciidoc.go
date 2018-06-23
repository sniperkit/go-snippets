package packages

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/alexandrecarlton/gogurt"
)

type AsciiDoc struct{}

func (asciidoc AsciiDoc) Name() string {
	return "asciidoc"
}

func (asciidoc AsciiDoc) URL(version string) string {
	return fmt.Sprintf("https://github.com/asciidoc/asciidoc/archive/%s.tar.gz", version)
}

func (asciidoc AsciiDoc) Build(config gogurt.Config) error {

	autoconf := exec.Command("autoconf")
	autoconf.Env = append(autoconf.Env, "PATH=" + config.BinDir(AutoConf{}) + ":" + os.Getenv("PATH"))
	if err := autoconf.Run(); err != nil {
		return err
	}
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(asciidoc),
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Paths: []string{
			config.BinDir(LibXSLT{}),
		},
	}.Cmd()
	return make.Run()
}

func (asciidoc AsciiDoc) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (asciidoc AsciiDoc) Dependencies() []gogurt.Package {
	// Note that this is a Python application, so some form of Python is needed to actually use it.
	return []gogurt.Package{
		LibXSLT{},
	}
}
