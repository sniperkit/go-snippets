package packages

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexandrecarlton/gogurt"
)

type Pigz struct{}

func (pigz Pigz) Name() string {
	return "pigz"
}

func (pigz Pigz) URL(version string) string {
	return fmt.Sprintf("https://zlib.net/pigz/pigz-%s.tar.gz", version)
}

func (pigz Pigz) Build(config gogurt.Config) error {
	gogurt.CopyFile(filepath.Join(config.IncludeDir(Zlib{}), "zlib.h"), config.BuildDir(pigz))
	gogurt.CopyFile(filepath.Join(config.IncludeDir(Zlib{}), "zconf.h"), config.BuildDir(pigz))
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Args: []string{
			"LDFLAGS=-static -L" + config.LibDir(Zlib{}),
		},
	}.Cmd()
	return make.Run()
}

func (pigz Pigz) Install(config gogurt.Config) error {
	os.MkdirAll(config.BinDir(pigz), 0755)
	return gogurt.CopyFile(filepath.Join(config.BuildDir(pigz), "pigz"), config.BinDir(pigz))
}

func (pigz Pigz) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Zlib{},
	}
}
