package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Libevent struct{}

func (libevent Libevent) Name() string {
	return "libevent"
}

func (libevent Libevent) URL(version string) string {
	return fmt.Sprintf("https://github.com/libevent/libevent/releases/download/release-%s/libevent-%s.tar.gz", version, version)
}

func (libevent Libevent) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.Prefix,
		Args: []string{
			"--enable-static",
			"--disable-shared",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (libevent Libevent) Install(config gogurt.Config) error {
	make := gogurt.MakeCmd{
		Args: []string{
			"install",
			"prefix=" + config.InstallDir(libevent),
		},
	}.Cmd()
	return make.Run()
}

func (libevent Libevent) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
