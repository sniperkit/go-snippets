package toolkit

import (
	"os"
	"path/filepath"
)

var (
	PathSeparator string = string(os.PathSeparator)
)

func PathDefault(removeSlash bool) string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if removeSlash == false {
		dir = dir + "/"
	}
	return dir
}

func IsFileNotExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}

	return false
}

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
