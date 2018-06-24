package adf

import (
	"github.com/juliengk/go-utils/filedir"
)

type DaemonConfig struct {
	App App
}

func (c *DaemonConfig) Init() error {
	c.App.Dir.Root = "/var/lib/metadata"

	if err := filedir.CreateDirIfNotExist(c.App.Dir.Root, false, 0700); err != nil {
		return err
	}

	return nil
}
