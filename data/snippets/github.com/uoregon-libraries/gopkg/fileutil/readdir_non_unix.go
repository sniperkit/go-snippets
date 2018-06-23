//+build !linux

package fileutil

import "os"

// Readdir wraps os.File's Readdir to handle common operations we need for
// getting a list of file info structures
func Readdir(path string) ([]os.FileInfo, error) {
	var d *os.File
	var err error

	d, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	var items []os.FileInfo
	items, err = d.Readdir(-1)
	d.Close()
	return items, err
}
