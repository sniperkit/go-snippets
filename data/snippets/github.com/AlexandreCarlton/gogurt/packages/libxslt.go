package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type LibXSLT struct{}

func (libxslt LibXSLT) Name() string {
	return "libxslt"
}

func (libxslt LibXSLT) URL(version string) string {
	return fmt.Sprintf("http://xmlsoft.org/sources/libxslt-%s.tar.gz", version)
}

func (libxslt LibXSLT) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(libxslt),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--with-libxml-prefix=" + config.InstallDir(LibXML2{}),
			"--without-python", // TODO: When we start building Python with expat, set this config.InstallDir(Python2{})
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (libxslt LibXSLT) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (libxslt LibXSLT) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		LibXML2{},
	}
}
