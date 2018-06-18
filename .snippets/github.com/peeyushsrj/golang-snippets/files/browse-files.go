package golangSnippets

import (
	"os"
	"path/filepath"
	"strings"
)

//BrowseXFiles take file type i.e. `.mp3` & directory as input and return array of filepaths and error.
func BrowseXFiles(x string, root string) ([]string, error) {
	var arr []string
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), x) { //.mp3
				arr = append(arr, path)
			}
		}
		return nil
	})
	if err != nil {
		return arr, err
	}
	return arr, nil
}
