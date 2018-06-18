// Package fileutil holds various general utilities for working with the
// filesystem more easily
package fileutil

import (
	"os"
	"path/filepath"
	"sort"
)

// IsDir returns true if the given path exists and is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	// This means something weird happened that we probably want to report (often
	// a permissions issue), but the function's purpose is simplicity, so we
	// consider this a non-file.
	if err != nil {
		return false
	}

	return info.IsDir()
}

// IsFile returns true if the given path exists and is a regular file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	// This means something weird happened that we probably want to report (often
	// a permissions issue), but the function's purpose is simplicity, so we
	// consider this a non-file.
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}

// Exists returns true if the given path exists and has no errors.  All errors
// are treated as the path not existing in order to avoid trying to determine
// what to do to handle the unknown errors which may be returned.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// MustNotExist is used when we need to be absolutely certain a path doesn't
// exist, such as when a directory's existence means a duplicate operation
// occurred.
func MustNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

// ReaddirSorted calls Readdir and sorts the results
func ReaddirSorted(path string) ([]os.FileInfo, error) {
	var fi, err = Readdir(path)
	if err == nil {
		sort.Sort(byName(fi))
	}

	return fi, err
}

// byName implements sort.Interface for sorting os.FileInfo data by name
type byName []os.FileInfo

func (n byName) Len() int           { return len(n) }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byName) Less(i, j int) bool { return n[i].Name() < n[j].Name() }

// SortFileInfos sorts a slice of os.FileInfo data by the underlying filename
func SortFileInfos(list []os.FileInfo) {
	sort.Sort(byName(list))
}

// FindIf iterates over all directory entries in the given path, running the
// given selector on each, and returning a list of those for which the selector
// returned true.
//
// Symlinks are resolved to their real file for the selector function, but the
// path added to the return will be a path to the symlink, not its target.
//
// Filesystem errors, including permission errors, will cause FindIf to halt
// and return an empty list and the error.
func FindIf(path string, selector func(i os.FileInfo) bool) ([]string, error) {
	var results []string
	var items, err = ReaddirSorted(path)
	if err != nil {
		return nil, err
	}

	for _, i := range items {
		var fName = i.Name()
		var path = filepath.Join(path, fName)
		var realPath = path
		if i.Mode()&os.ModeSymlink != 0 {
			realPath, err = os.Readlink(path)
			if err != nil {
				return nil, err
			}
			// Symlinks kind of suck - they can be absolute or relative, and if
			// they're relative we have to make them absolute....
			if !filepath.IsAbs(realPath) {
				realPath = filepath.Join(path, realPath)
			}

			i, err = os.Stat(realPath)
			if err != nil {
				return nil, err
			}
		}
		realPath = filepath.Clean(realPath)

		// See if the selector allows this file to be put in the list
		if !selector(i) {
			continue
		}

		results = append(results, path)
	}

	return results, nil
}

// FindFiles returns a list of all entries in a given path which are *not*
// directories or symlinks to directories.  For the purpose of this function,
// we define "files" as "things from which we can directly read data".
func FindFiles(path string) ([]string, error) {
	return FindIf(path, func(i os.FileInfo) bool {
		return !i.IsDir()
	})
}

// FindDirectories returns a list of all directories or symlinks to directories
// within the given path
func FindDirectories(path string) ([]string, error) {
	return FindIf(path, func(i os.FileInfo) bool {
		return i.IsDir()
	})
}

// Find traverses the filesystem to the given depth, returning only the items
// that are found at that depth.  Traverses symlinks if they are directories.
// Returns the first error found if any occur.
func Find(root string, depth int) ([]string, error) {
	var paths = []string{root}
	var newPaths []string
	for depth > 0 {
		for _, p := range paths {
			var appendList, err = FindDirectories(p)
			if err != nil {
				return nil, err
			}
			newPaths = append(newPaths, appendList...)
		}
		paths = newPaths
		newPaths = nil
		depth--
	}

	return paths, nil
}
