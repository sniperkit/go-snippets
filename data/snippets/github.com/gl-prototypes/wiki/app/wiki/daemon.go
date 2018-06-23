package wiki

import (
	"path"
	"runtime"
)

// InitializeDaemon is used to add the wikis templates to
// the console components views.
func (c *Component) InitializeDaemon() error {
	_, filename, _, _ := runtime.Caller(0)
	c.Console.Views.AddPath(path.Join(path.Dir(filename), "html"))
	return nil
}
