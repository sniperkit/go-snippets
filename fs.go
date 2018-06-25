//
// fs.go defines local file system support for storage.go
//
package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// FSConfigure initialize a store to a local disc system type
func fsConfigure(store *Store) (*Store, error) {
	store.Type = FS

	// Basic CRUD ops
	store.Create = func(fname string, rd io.Reader) error {
		return fsCreate(store, fname, rd)
	}
	store.Read = func(fname string) ([]byte, error) {
		return ioutil.ReadFile(fname)
	}
	store.Update = func(fname string, rd io.Reader) error {
		return fsUpdate(store, fname, rd)
	}
	store.Delete = func(fname string) error {
		return os.Remove(fname)
	}

	// Extra ops for compatibility with os.* and ioutil.*
	store.Stat = func(fname string) (os.FileInfo, error) {
		return os.Stat(fname)
	}
	store.Mkdir = func(name string, perm os.FileMode) error {
		return os.Mkdir(name, perm)
	}
	store.MkdirAll = func(path string, perm os.FileMode) error {
		return os.MkdirAll(path, perm)
	}
	store.Remove = func(name string) error {
		return os.Remove(name)
	}
	store.RemoveAll = func(path string) error {
		return os.RemoveAll(path)
	}
	store.ReadFile = func(fname string) ([]byte, error) {
		return ioutil.ReadFile(fname)
	}
	store.WriteFile = func(fname string, data []byte, perm os.FileMode) error {
		return ioutil.WriteFile(fname, data, perm)
	}

	// Extended ops for datatools and dataset

	// WriteFilter writes a file after running/applying a filter function to its' file pointer
	// E.g. composing a tarball before storing
	store.WriteFilter = func(finalPath string, processor func(*os.File) error) error {
		// Open temp file as file point
		tmp, err := ioutil.TempFile(os.TempDir(), path.Base(finalPath))
		if err != nil {
			return err
		}
		tmpName := tmp.Name()

		// Envoke processor function
		err = processor(tmp)
		if err != nil {
			return err
		}
		err = tmp.Close()
		if err != nil {
			return err
		}
		// Now we're ready to but the processed file in its place
		if err := os.Rename(tmpName, finalPath); err != nil {
			// If rename files try copy then delete.
			in, err := os.Open(tmpName)
			if err != nil {
				return err
			}
			defer in.Close()
			out, err := os.Create(finalPath)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, in)
			if err != nil {
				return err
			}
			return os.Remove(tmpName)
		}
		return nil
	}

	// Now the store is setup and we're ready to return
	return store, nil
}

// fsCreate creates a new file on the file system with a given name from the byte array.
func fsCreate(s *Store, fname string, rd io.Reader) error {
	// FIXME: FSCreate should create the path elements only if necessary
	dname := path.Dir(fname)
	os.MkdirAll(dname, 0775)
	wr, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer wr.Close()
	_, err = io.Copy(wr, rd)
	if err != nil {
		return fmt.Errorf("%s, %s", fname, err)
	}
	return nil
}

// fsUpdate replaces a file on the file system with the contents fo byte array returning error.
// It will truncate the file if necessary.
func fsUpdate(s *Store, fname string, rd io.Reader) error {
	wr, err := os.OpenFile(fname, os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer wr.Close()
	_, err = io.Copy(wr, rd)
	if err != nil {
		return fmt.Errorf("%s, %s", fname, err)
	}
	return nil
}
