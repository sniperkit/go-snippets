// Package docindex fetches and lists text documents (with metadata "front
// matter" at the top of each file) in a directory tree.
package docindex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
)

// Index represents a directory containing text documents with metadata.
type Index struct {
	dir      string
	metaType reflect.Type

	docsLock  sync.Mutex
	docMeta   map[string]interface{}
	docData   map[string][]byte
	filenames []string
}

// Open reads all documents in dir and its subdirectories. Metadata JSON is
// decoded into structs whose type is pointer-to-metadataType.
//
// If dir is not an existing directory, an error satisfying os.IsNotExist is
// returned.
func Open(dir string, metadataType interface{}) (*Index, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("not a directory: %q", dir)
	}
	idx := &Index{dir: dir, metaType: reflect.TypeOf(metadataType)}
	err = idx.Reload()
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// Reload refreshes the list of documents in the index.
func (idx *Index) Reload() error {
	idx.docsLock.Lock()
	defer idx.docsLock.Unlock()

	// Load all docs in dir.
	idx.docMeta = make(map[string]interface{})
	idx.docData = make(map[string][]byte)
	err := filepath.Walk(idx.dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.Mode().IsRegular() {
			err = idx.loadFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Cache list of filenames.
	idx.filenames = make([]string, len(idx.docData))
	i := 0
	for filename, _ := range idx.docData {
		idx.filenames[i] = filename
		i++
	}
	sort.Strings(idx.filenames)

	return nil
}

func (idx *Index) loadFile(file string) error {
	// Assume the caller holds idx.docsLock.

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	delim := []byte("\n\n")
	mdStart := bytes.Index(data, delim)
	if mdStart >= 0 {
		md := data[:mdStart]
		mdObjPtr := reflect.New(idx.metaType).Interface()
		err = json.Unmarshal(md, mdObjPtr)
		if err != nil {
			return err
		}
		idx.docMeta[file] = mdObjPtr
	}
	idx.docData[file] = data[mdStart+len(delim):]
	return nil
}

type notFoundError struct {
	file string
}

func (e notFoundError) Error() string {
	return fmt.Sprintf("doc not found: %q", e.file)
}

func IsNotFound(err error) bool {
	_, ok := err.(notFoundError)
	return ok
}

// Filenames returns a list of document filenames that were loaded in the
// last call to Open or Reload.
func (idx *Index) Filenames() []string {
	return idx.filenames
}

// AllMetadata stores the metadata for each document in meta, which should be a
// map with string keys and values whose type is pointer-to-metadataType.
func (idx *Index) AllMetadata(metaMap interface{}) {
	idx.docsLock.Lock()
	defer idx.docsLock.Unlock()

	metaMapV := reflect.ValueOf(metaMap)
	if metaMapV.IsNil() {
		panic("nil metaMap")
	}

	for _, filename := range idx.Filenames() {
		var metaV reflect.Value
		if meta, present := idx.docMeta[filename]; present {
			metaV = reflect.ValueOf(meta)
		} else {
			metaV = reflect.New(idx.metaType)
		}
		metaMapV.SetMapIndex(reflect.ValueOf(filename), metaV)
	}
}

// Doc returns the data in the document named by filename. If meta is a pointer
// a struct of the metadataType used in Open, the document metadata is written
// to that struct.
func (idx *Index) Doc(filename string, meta interface{}) ([]byte, error) {
	idx.docsLock.Lock()
	defer idx.docsLock.Unlock()

	if _, exists := idx.docData[filename]; !exists {
		return nil, notFoundError{filename}
	}

	metaVPtr := reflect.ValueOf(meta)
	if !metaVPtr.IsNil() {
		metaV := metaVPtr.Elem()
		if docMeta, hasMeta := idx.docMeta[filename]; hasMeta {
			metaV.Set(reflect.ValueOf(docMeta).Elem())
		} else {
			metaV.Set(reflect.Zero(idx.metaType))
		}
	}

	return idx.docData[filename], nil
}
