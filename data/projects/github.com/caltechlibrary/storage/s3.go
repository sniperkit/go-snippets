//
// s3.go adds s3:// (Amazon S3 storage) support to storage.go
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

	// 3rd Party Packages
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// s3ObjectInfo is a map so we can create a os.FileInfo compatible struct from S3 objects
type s3ObjectInfo struct {
	Info map[string]interface{}
}

// String returns a string representation of the object reported by ListObjects
func (d *s3ObjectInfo) String() string {
	src, err := json.Marshal(d.Info)
	if err != nil {
		return fmt.Sprintf("%+v", d.Info)
	}
	return string(src)
}

// Create a s3ObjectInfo struct from response Contents return by ListObjects on S3
func s3ToObjectInfo(o *s3.Object) *s3ObjectInfo {
	doc := new(s3ObjectInfo)
	doc.Info = map[string]interface{}{}
	doc.Info["ETag"] = o.ETag
	doc.Info["Key"] = o.Key
	doc.Info["LastModified"] = o.LastModified
	doc.Info["Owner"] = o.Owner
	doc.Info["Size"] = o.Size
	doc.Info["StorageClass"] = o.StorageClass
	return doc
}

// Name returns the Key after evaluating with path.Base() so we match os.FileInfo.Name()
// or an empty string
func (d *s3ObjectInfo) Name() string {
	if val, ok := d.Info["Key"]; ok == true {
		p := val.(*string)
		return path.Base(*p)
	}
	return ""
}

// Size returns the size of an object reported by listing the object
// Or zero as a int64 if not available
func (d *s3ObjectInfo) Size() int64 {
	if val, ok := d.Info["Size"]; ok == true {
		size := val.(*int64)
		return *size
	}
	return int64(0)
}

// ModTime returns the value of LastModied reported by listing the object
// or an empty Time object if not available
func (d *s3ObjectInfo) ModTime() time.Time {
	if val, ok := d.Info["LastModified"]; ok == true {
		t := val.(*time.Time)
		return *t
	}
	return time.Time{}
}

// Mode returns the file mode but this doens't map to S3 so we return zero always
func (d *s3ObjectInfo) Mode() os.FileMode {
	return os.FileMode(0)
}

// IsDir returns false, S3 doens't support the concept of directories only keys in buckets
func (d *s3ObjectInfo) IsDir() bool {
	return false
}

// Sys() returns an system dependant interface...
func (d *s3ObjectInfo) Sys() interface{} {
	return nil
}

// s3Configure is a function that configures a storage.Store for use with AWS S3
func s3Configure(store *Store) (*Store, error) {
	var (
		awsSDKLoadConfig       bool
		awsSharedConfigEnabled bool
		awsProfile             string

		sess *session.Session
		err  error
	)

	// Set storage type to S3
	store.Type = S3

	if val, ok := store.Config["AwsSDKLoadConfig"]; ok == true {
		awsSDKLoadConfig = val.(bool)
	}
	if val, ok := store.Config["AwsSharedConfigEnabled"]; ok == true {
		awsSharedConfigEnabled = val.(bool)
	}

	if awsSDKLoadConfig == true {
		var opts session.Options
		if awsProfile != "" {
			opts.Profile = awsProfile
		} else {
			opts.Profile = "default"
		}
		if awsSharedConfigEnabled == true {
			opts.SharedConfigState = session.SharedConfigEnable
		} else {
			opts.SharedConfigState = session.SharedConfigDisable
		}
		if val, ok := store.Config["AwsRegion"]; ok == true {
			opts.Config = aws.Config{Region: aws.String(val.(string))}
		}

		sess, err = session.NewSessionWithOptions(opts)
		if err != nil {
			return nil, err
		}
	} else {
		sess, err = session.NewSession()
		if err != nil {
			return nil, err
		}
	}
	store.Config["session"] = sess

	s3Svc := s3.New(sess)
	store.Config["s3Service"] = s3Svc

	// Basic Ops
	store.Create = func(fname string, rd io.Reader) error {
		return s3Create(store, fname, rd)
	}
	store.Read = func(fname string) ([]byte, error) {
		return s3Read(store, fname)
	}
	store.Update = func(fname string, rd io.Reader) error {
		// NOTE: Create and Update and the same in S3, Update overwrites the existing object
		return s3Create(store, fname, rd)
	}
	store.Delete = func(fname string) error {
		return s3Remove(store, fname)
	}

	// Extra ops for compatibilty with os.* and ioutil.*
	store.Stat = func(fname string) (os.FileInfo, error) {
		return s3Stat(store, fname)
	}
	store.Mkdir = func(name string, perm os.FileMode) error {
		//NOTE: S3 lacks the concept of directories, the fill path is the "Key" value info the bucket
		return nil
	}
	store.MkdirAll = func(path string, perm os.FileMode) error {
		//NOTE: S3 lacks the concept of directories, the fill path is the "Key" value info the bucket
		return nil
	}
	store.Remove = func(fname string) error {
		return s3Remove(store, fname)
	}
	store.RemoveAll = func(fname string) error {
		return s3RemoveAll(store, fname)
	}
	store.ReadFile = func(fname string) ([]byte, error) {
		return s3Read(store, fname)
	}
	store.WriteFile = func(fname string, data []byte, perm os.FileMode) error {
		return s3Create(store, fname, bytes.NewBuffer(data))
	}

	// Extended options for datatools and dataset

	// WriteFilter writes a file after running apply a filter function to its' file pointer
	// E.g. composing a tarball before uploading results to S3
	store.WriteFilter = func(finalPath string, processor func(fp *os.File) error) error {
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

		// Now we're ready to upload results
		buf, err := ioutil.ReadFile(tmpName)
		if err != nil {
			return err
		}
		return s3Create(store, finalPath, bytes.NewReader(buf))
	}

	// Now the store is setup and we're ready to return
	return store, nil
}

// S3Stat takes a file name and returns a FileInfo and error value
func s3Stat(s *Store, fname string) (os.FileInfo, error) {
	if val, ok := s.Config["s3Service"]; ok == true {
		s3Svc := val.(s3iface.S3API)
		if _, ok := s.Config["AwsBucket"]; ok == false {
			return nil, fmt.Errorf("Bucket not defined for %s", fname)
		}
		bucketName := s.Config["AwsBucket"].(string)
		statParams := &s3.ListObjectsInput{
			Bucket:  &bucketName,
			Prefix:  &fname,
			MaxKeys: aws.Int64(1),
		}
		res, err := s3Svc.ListObjects(statParams)
		if err != nil {
			return nil, err
		}
		// NOTE: Only return the fname we're looking for not the other ones with matching prefix
		for _, obj := range res.Contents {
			oInfo := s3ToObjectInfo(obj)
			if strings.Compare(oInfo.Name(), path.Base(fname)) == 0 {
				return oInfo, nil
			}
		}
		return nil, fmt.Errorf("%s not found", fname)
	}
	return nil, fmt.Errorf("s3Service not configured")
}

// Create takes a full path and a byte array of content and writes it to the bucket
// associated with the Store initialized.
func s3Create(s *Store, fname string, rd io.Reader) error {
	if val, ok := s.Config["s3Service"]; ok == true {
		s3Svc := val.(s3iface.S3API)
		if _, ok := s.Config["AwsBucket"]; ok == false {
			return fmt.Errorf("Bucket not defined for %s", fname)
		}
		bucketName := s.Config["AwsBucket"].(string)
		upParams := &s3manager.UploadInput{
			Bucket: &bucketName,
			Key:    &fname,
			Body:   rd,
		}
		uploader := s3manager.NewUploaderWithClient(s3Svc)
		_, err := uploader.Upload(upParams)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("s3Service not configured")
}

// s3Read takes a full path and returns a byte array and error from the bucket read
func s3Read(s *Store, fname string) ([]byte, error) {
	if val, ok := s.Config["s3Service"]; ok == true {
		s3Svc := val.(s3iface.S3API)
		if _, ok := s.Config["AwsBucket"]; ok == false {
			return nil, fmt.Errorf("Bucket not defined for %s", fname)
		}
		bucketName := s.Config["AwsBucket"].(string)
		downloadParams := &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &fname,
		}
		buf := &aws.WriteAtBuffer{}
		downloader := s3manager.NewDownloaderWithClient(s3Svc)
		_, err := downloader.Download(buf, downloadParams)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("s3Service not configured")
}

// s3Remove takes a full path and returns an error if delete not successful
func s3Remove(s *Store, fname string) error {
	if val, ok := s.Config["s3Service"]; ok == true {
		s3Svc := val.(s3iface.S3API)
		if _, ok := s.Config["AwsBucket"]; ok == false {
			return fmt.Errorf("Bucket not defined for %s", fname)
		}
		bucketName := s.Config["AwsBucket"].(string)
		deleteParams := &s3.DeleteObjectInput{
			Bucket: &bucketName,
			Key:    &fname,
		}
		_, err := s3Svc.DeleteObject(deleteParams)
		return err
	}
	return fmt.Errorf("s3Service not configured")
}

// s3RemoveAll takes a path prefix and delete matching items (up to 1000) and returns an error if delete not successful
func s3RemoveAll(s *Store, prefixName string) error {
	if val, ok := s.Config["s3Service"]; ok == true {
		s3Svc := val.(s3iface.S3API)
		if _, ok := s.Config["AwsBucket"]; ok == false {
			return fmt.Errorf("Bucket not defined for %s", prefixName)
		}
		bucketName := s.Config["AwsBucket"].(string)
		// FIXME: Get a list of objects, then delate each one
		statParams := &s3.ListObjectsInput{
			Bucket: &bucketName,
			Prefix: &prefixName,
		}
		// S3 ListObjects returns an maximum of 1000 objects, I am using an outer loop to handle
		// the case of where the prefix matches more than 1000 objects.
		res, err := s3Svc.ListObjects(statParams)
		if err != nil {
			return err
		}
		cnt := len(res.Contents)
		for cnt > 0 {
			// NOTE: Only return the fname we're looking for not the other ones with matching prefix
			for _, obj := range res.Contents {
				deleteParams := &s3.DeleteObjectInput{
					Bucket: &bucketName,
					Key:    obj.Key,
				}
				_, err := s3Svc.DeleteObject(deleteParams)
				if err != nil {
					return err
				}
			}
			res, err := s3Svc.ListObjects(statParams)
			if err != nil {
				return err
			}
			cnt = len(res.Contents)
		}
		return nil
	}
	return fmt.Errorf("s3Service not configured")
}
