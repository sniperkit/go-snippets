// Helpers for working with files and filesystems
package fileutil

import (
	"net/http"
	"os"
	"regexp"
)

type RewriteFileSystem struct {
	FileSystem http.FileSystem
	Find       *regexp.Regexp
	Replace    string
	MustMatch  bool
}

func (self RewriteFileSystem) Open(name string) (http.File, error) {
	if self.Find != nil {
		if self.Find.MatchString(name) {
			name = self.Find.ReplaceAllString(name, self.Replace)
		} else if self.MustMatch {
			return nil, os.ErrNotExist
		}
	}

	return self.FileSystem.Open(name)
}
