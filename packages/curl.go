package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

// Note that this does not generate a completely static binary,
// but one that links to system libraries, like linux-vdso.so.1 and libc.so.6.

type Curl struct{}

func (curl Curl) Name() string {
	return "curl"
}

func (curl Curl) URL(version string) string {
	return fmt.Sprintf("https://curl.haxx.se/download/curl-%s.tar.gz", version)
}

func (curl Curl) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(curl),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--with-ssl=" + config.InstallDir(OpenSSL{}),
			"--with-zlib=" + config.InstallDir(Zlib{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (curl Curl) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (curl Curl) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		OpenSSL{},
		Zlib{},
	}
}
