//
// gs.go adds Google Cloud Storage (gs://) support to storage.go
//
package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	// Google Cloud SDK/API
	gstorage "cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

// gsObjectInfo is a map so we can create a os.FileInfo compatible map from GS objects
type gsObjectInfo struct {
	Info map[string]interface{}
}

// GSConfigure is a function that configures a storage.Store for use with Google Cloud Storage
func gsConfigure(store *Store) (*Store, error) {
	// Set storage type to GS
	store.Type = GS

	// Initialization needed by Google Cloud Storage
	ctx := context.Background()
	client, err := gstorage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	store.Config["gsService"] = client

	// Basic Ops
	store.Create = func(fname string, rd io.Reader) error {
		return gsCreate(store, fname, rd)
	}
	store.Read = func(fname string) ([]byte, error) {
		return gsRead(store, fname)
	}
	store.Update = func(fname string, rd io.Reader) error {
		// NOTE: Create and Update and the same in GS, Update overwrites the existing object
		return gsCreate(store, fname, rd)
	}
	store.Delete = func(fname string) error {
		return gsRemove(store, fname)
	}

	// Extra ops for compatibilty with os.* and ioutil.*
	store.Stat = func(fname string) (os.FileInfo, error) {
		return gsStat(store, fname)
	}
	store.Mkdir = func(name string, perm os.FileMode) error {
		//NOTE: GS lacks the concept of directories, the fill path is the "Key" value info the bucket
		return nil
	}
	store.MkdirAll = func(path string, perm os.FileMode) error {
		//NOTE: GS lacks the concept of directories, the fill path is the "Key" value info the bucket
		return nil
	}
	store.Remove = func(fname string) error {
		return gsRemove(store, fname)
	}
	store.RemoveAll = func(fname string) error {
		return gsRemoveAll(store, fname)
	}
	store.ReadFile = func(fname string) ([]byte, error) {
		return gsRead(store, fname)
	}
	store.WriteFile = func(fname string, data []byte, perm os.FileMode) error {
		return gsCreate(store, fname, bytes.NewBuffer(data))
	}

	// Extended options for datatools and dataset

	// WriteFilter writes a file after running apply a filter function to its' file pointer
	// E.g. composing a tarball before uploading results to S3 or GS
	store.WriteFilter = func(finalPath string, processor func(*os.File) error) error {
		// Open temp file as file point
		tmp, err := ioutil.TempFile(os.TempDir(), path.Base(finalPath))
		if err != nil {
			return err
		}
		tmpName := tmp.Name()
		defer os.Remove(tmpName)

		// Envoke processor function
		err = processor(tmp)
		if err != nil {
			return err
		}
		err = tmp.Close()
		if err != nil {
			return err
		}
		// OK now we're ready to upload temp filename to final path
		buf, err := ioutil.ReadFile(tmpName)
		if err != nil {
			return err
		}
		return gsCreate(store, finalPath, bytes.NewReader(buf))
	}

	// Now the store is setup and we're ready to return
	return store, nil
}

// GSCreate takes a full path and a byte array of content and writes it to the bucket
// associated with the Store initialized.
func gsCreate(s *Store, fname string, rd io.Reader) error {
	if val, ok := s.Config["gsService"]; ok == true {
		gsSrv := val.(*gstorage.Client)
		val, ok = s.Config["GoogleBucket"]
		if ok == false {
			return fmt.Errorf("gsService not configured")
		}
		bucketName := val.(string)
		ctx := context.Background()
		wr := gsSrv.Bucket(bucketName).Object(fname).NewWriter(ctx)
		if _, err := io.Copy(wr, rd); err != nil {
			return err
		}
		if err := wr.Close(); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("gsService not configured")
}

// GSRead takes a full path and returns a byte array and error from the bucket read
func gsRead(s *Store, fname string) ([]byte, error) {
	if val, ok := s.Config["gsService"]; ok == true {
		gsSrv := val.(*gstorage.Client)
		val, ok = s.Config["GoogleBucket"]
		if ok == false {
			return nil, fmt.Errorf("gsService not configured")
		}
		bucketName := val.(string)
		ctx := context.Background()

		rd, err := gsSrv.Bucket(bucketName).Object(fname).NewReader(ctx)
		if err != nil {
			return nil, err
		}
		defer rd.Close()

		data, err := ioutil.ReadAll(rd)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("gsService not configured")
}

// GSRemove takes a full path and returns an error if delete not successful
func gsRemove(s *Store, fname string) error {
	if val, ok := s.Config["gsService"]; ok == true {
		gsSrv := val.(*gstorage.Client)
		val, ok = s.Config["GoogleBucket"]
		if ok == false {
			return fmt.Errorf("gsService not configured")
		}
		bucketName := val.(string)
		ctx := context.Background()

		return gsSrv.Bucket(bucketName).Object(fname).Delete(ctx)
	}
	return fmt.Errorf("gsService not configured")
}

// GSRemoveAll takes a path prefix and delete matching items (up to 1000) and returns an error if delete not successful
func gsRemoveAll(s *Store, prefixName string) error {
	var errors []string
	if val, ok := s.Config["gsService"]; ok == true {
		gsSrv := val.(*gstorage.Client)
		val, ok = s.Config["GoogleBucket"]
		if ok == false {
			return fmt.Errorf("gsService not configured")
		}
		bucketName := val.(string)
		ctx := context.Background()

		bucket := gsSrv.Bucket(bucketName)
		o := bucket.Objects(ctx, nil)
		for {
			attrs, err := o.Next()
			if err != nil && err != iterator.Done {
				return fmt.Errorf("Can't get next object, %s")
			}
			if err == iterator.Done {
				break
			}
			// Make sure we're at the right level of the pseudo path before we do a delete
			if strings.HasPrefix(attrs.Name, prefixName) {
				if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
					errors = append(errors, fmt.Sprintf("%s, %s", attrs.Name, err))
				}
			}
		}
		if len(errors) > 0 {
			return fmt.Errorf("%s", strings.Join(errors, "\n"))
		}
		return nil
	}
	return fmt.Errorf("gsService not configured")
}

// Create a gsObjectInfo struct from response Contents return by ListObjects on Google Cloud Storage
func gsToObjectInfo(o *gstorage.ObjectAttrs) *gsObjectInfo {
	doc := new(gsObjectInfo)
	doc.Info = map[string]interface{}{}
	//Q: What fields Do I need from gstorage.ObjectAttrs that can be combined to be an ETag?
	//doc.Info["ETag"] = o.ETag
	doc.Info["Key"] = o.Name
	doc.Info["LastModified"] = o.Updated
	doc.Info["Owner"] = o.Owner
	doc.Info["Size"] = o.Size
	doc.Info["StorageClass"] = o.StorageClass
	//Q: What other attributes should I pass through?
	return doc
}

// String returns a string representation of the object reported by ListObjects
func (d *gsObjectInfo) String() string {
	src, err := json.Marshal(d.Info)
	if err != nil {
		return fmt.Sprintf("%+v", d.Info)
	}
	return string(src)
}

// Name returns the Key after evaluating with path.Base() so we match os.FileInfo.Name()
// or an empty string
func (d *gsObjectInfo) Name() string {
	if val, ok := d.Info["Key"]; ok == true {
		p := val.(string)
		return path.Base(p)
	}
	return ""
}

// Size returns the size of an object reported by listing the object
// Or zero as a int64 if not available
func (d *gsObjectInfo) Size() int64 {
	if val, ok := d.Info["Size"]; ok == true {
		size := val.(int64)
		return size
	}
	return int64(0)
}

// ModTime returns the value of LastModied reported by listing the object
// or an empty Time object if not available
func (d *gsObjectInfo) ModTime() time.Time {
	if val, ok := d.Info["LastModified"]; ok == true {
		t := val.(*time.Time)
		return *t
	}
	return time.Time{}
}

// IsDir returns false, Google Cloud Storage doesn't support the concept of directories only keys in buckets
func (d *gsObjectInfo) IsDir() bool {
	return false
}

// Sys() returns an system dependant interface...
func (d *gsObjectInfo) Sys() interface{} {
	return nil
}

// Mode returns the file mode but this doens't map to Google Cloud Storage so we return zero always
func (d *gsObjectInfo) Mode() os.FileMode {
	return os.FileMode(0)
}

// GSStat takes a file name and returns a FileInfo and error value
func gsStat(s *Store, fname string) (os.FileInfo, error) {
	if val, ok := s.Config["gsService"]; ok == true {
		gsSrv := val.(*gstorage.Client)
		val, ok = s.Config["GoogleBucket"]
		if ok == false {
			return nil, fmt.Errorf("gsService not configured")
		}
		bucketName := val.(string)
		ctx := context.Background()

		o := gsSrv.Bucket(bucketName).Object(fname)
		attrs, err := o.Attrs(ctx)
		if err != nil {
			return nil, err
		}
		// See https://godoc.org/cloud.google.com/go/storage#ObjectAttrs for attribute mapping
		oInfo := gsToObjectInfo(attrs)
		return oInfo, nil

	}
	return nil, fmt.Errorf("gsService not configured")
}
