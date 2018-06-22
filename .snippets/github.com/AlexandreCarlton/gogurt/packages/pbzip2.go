package packages

import (
	"fmt"
	"strings"
	"github.com/alexandrecarlton/gogurt"
)

type PBZip2 struct{}

func (pbzip2 PBZip2) Name() string {
	return "pbzip2"
}

func (pbzip2 PBZip2) URL(version string) string {
	majorVersion := strings.Join(strings.Split(version, ".")[:2], ".")
	return fmt.Sprintf("https://launchpad.net/pbzip2/%s/%s/+download/pbzip2-%s.tar.gz", majorVersion, version, version)
}

func (pbzip2 PBZip2) Build(config gogurt.Config) error {
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Args: []string{
			"CXXFLAGS=-O2 -I" + config.IncludeDir(Bzip2{}),
			"LDFLAGS=-L" + config.LibDir(Bzip2{}),
		},
	}.Cmd()
	return make.Run()
}

func (pbzip2 PBZip2) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			"PREFIX=" + config.InstallDir(pbzip2),
			"install",
		}}.Cmd()
	return makeInstall.Run()
}

func (pbzip2 PBZip2) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Bzip2{},
	}
}
