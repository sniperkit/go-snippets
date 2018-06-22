package packages

// golang of versions > 1.4 require an existing go installation.
// The recomended course of action to build from source is to download
// a pre-built installation of golang to use in building the newer one.
// As such, we just download the newest one.

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type Go struct{}

func (golang Go) Name() string {
	return "go"
}

func (golang Go) URL(version string) string {
	return fmt.Sprintf("https://storage.googleapis.com/golang/go%s.linux-amd64.tar.gz", version)
}

func (golang Go) Build(config gogurt.Config) error {
	// It's prebuilt!
	return nil
}

func (golang Go) Install(config gogurt.Config) error {
	err := os.MkdirAll(config.InstallDir(golang), 0755)
	if err != nil {
		return err
	}
	return gogurt.CopyFile(
		filepath.Join(config.BuildDir(golang), "bin"),
		config.BinDir(golang),
	)
}

func (golang Go) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
