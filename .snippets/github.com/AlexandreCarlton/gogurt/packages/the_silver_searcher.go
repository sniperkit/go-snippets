package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type TheSilverSearcher struct{}

func (ag TheSilverSearcher) Name() string {
	return "the_silver_searcher"
}

func (ag TheSilverSearcher) URL(version string) string {
	return fmt.Sprintf("http://geoff.greer.fm/ag/releases/the_silver_searcher-%s.tar.gz", version)
}

func (ag TheSilverSearcher) Build(config gogurt.Config) error {

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(ag),
		Args: []string{
			"--enable-zlib",
			"--enable-lzma",
		},
		CFlags: []string{
			"-I" + config.IncludeDir(Pcre{}),
			"-I" + config.IncludeDir(XZ{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		CppFlags: []string{
			"-I" + config.IncludeDir(Pcre{}),
			"-I" + config.IncludeDir(XZ{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Pcre{}),
			"-L" + config.LibDir(XZ{}),
			"-L" + config.LibDir(Zlib{}),
		},
		Libs: []string{
			"-lpcre",
			"-llzma",
			"-lz",
		},
	}.Cmd()

	configure.Env = append(
		configure.Env,
		"PCRE_CFLAGS=-I" + config.IncludeDir(Pcre{}),
		"PCRE_LIBS=-L" + config.LibDir(Pcre{}) + " -lpcre",
		"LZMA_CFLAGS=-I" + config.IncludeDir(XZ{}),
		"LZMA_LIBS=-L" + config.LibDir(XZ{}) + " -llzma",
	)
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (ag TheSilverSearcher) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (ag TheSilverSearcher) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Pcre{},
		XZ{},
		Zlib{},
	}
}
