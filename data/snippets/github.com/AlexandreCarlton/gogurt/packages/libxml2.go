package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type LibXML2 struct{}

func (libxml2 LibXML2) Name() string {
	return "libxml2"
}

func (libxml2 LibXML2) URL(version string) string {
	return fmt.Sprintf("http://xmlsoft.org/sources/libxml2-%s.tar.gz", version)
}

func (libxml2 LibXML2) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(libxml2),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--with-history", // enable history support to xmllint shell
			"--with-readline=" + config.InstallDir(ReadLine{}),
			"--with-lzma=" + config.InstallDir(XZ{}),
			"--with-zlib=" + config.InstallDir(Zlib{}),
			"--without-python", // TODO: When we start building Python with expat, set this config.InstallDir(Python2{})
		},
		Paths: []string{
			config.BinDir(Python2{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (libxml2 LibXML2) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (libxml2 LibXML2) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		ReadLine{},
		XZ{},
		Zlib{},
	}
}
