package main

import (
	"os"
	"path/filepath"
)

func createFileWithPath(path string) (f *os.File, err error) {
	dir, _ := filepath.Split(path)

	if dir != "" {
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return nil, err
		}
	}

	f, err = os.Create(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}
