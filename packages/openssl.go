package packages

import (
	"fmt"
	"os/exec"

	"github.com/alexandrecarlton/gogurt"
)

type OpenSSL struct{}

func (openssl OpenSSL) Name() string {
	return "openssl"
}

func (openssl OpenSSL) URL(version string) string {
	return fmt.Sprintf("https://www.openssl.org/source/openssl-%s.tar.gz", version)
}

func (openssl OpenSSL) Build(config gogurt.Config) error {
	zlib := Zlib{}
	configure := exec.Command(
		"./config",
		"--prefix=" + config.InstallDir(openssl),
		"no-shared",
		"--with-zlib-include=" + config.IncludeDir(zlib),
		"--with-zlib-lib=" + config.LibDir(zlib),
	)
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (openssl OpenSSL) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
		Jobs: config.NumCores,
	}.Cmd()
	return makeInstall.Run()
}

func (openssl OpenSSL) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Zlib{},
	}
}

