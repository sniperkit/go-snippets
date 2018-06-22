package packages

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type NodeJS struct{}

func (nodejs NodeJS) Name() string {
	return "nodejs"
}

func (nodejs NodeJS) URL(version string) string {
	return fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%s-linux-x64.tar.xz", version, version)
}

func (nodejs NodeJS) Build(config gogurt.Config) error {
	// It's prebuilt!
	return nil
}

func (nodejs NodeJS) Install(config gogurt.Config) error {
	err := os.MkdirAll(config.InstallDir(nodejs), 0755)
	if err != nil {
		return err
	}

	err = gogurt.CopyFile(filepath.Join(config.BuildDir(nodejs), "bin"), config.BinDir(nodejs))
	if err != nil {
		return err
	}

	err = gogurt.CopyFile(filepath.Join(config.BuildDir(nodejs), "include"), config.IncludeDir(nodejs))
	if err != nil {
		return err
	}

	err = gogurt.CopyFile(filepath.Join(config.BuildDir(nodejs), "lib"), config.LibDir(nodejs))
	if err != nil {
		return err
	}

	return gogurt.CopyFile(filepath.Join(config.BuildDir(nodejs), "share"), filepath.Join(config.InstallDir(nodejs), "share"))
}

func (nodejs NodeJS) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}

