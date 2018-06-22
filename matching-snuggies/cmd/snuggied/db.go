package main

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/gophergala/matching-snuggies/slicerjob"
)

func b(s string) []byte {
	return []byte(s)
}

var DB *bolt.DB

const (
	dbJobs       = "jobs"
	dbMeshFiles  = "meshFiles"
	dbGCodeFiles = "gCodeFiles"
)

func loadDB(path string) *bolt.DB {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(b(dbMeshFiles))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(b(dbGCodeFiles))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(b(dbJobs))
		if err != nil {
			return err
		}
		return nil
	})
	return db
}

func PutMeshFile(key string, path string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b(dbMeshFiles)).
			Put(b(key), b(path))
	})
}

func PutGCodeFile(key string, value string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b(dbGCodeFiles))
		if bucket == nil {
			return fmt.Errorf("%v bucket doesn't exist!", dbGCodeFiles)
		}
		return bucket.Put(b(key), b(value))
	})
}

func PutJob(key string, job *slicerjob.Job) error {
	jsonJob, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return DB.Update(func(tx *bolt.Tx) error {
		bucketName := "jobs"
		bucket := tx.Bucket(b(bucketName))
		if bucket == nil {
			return fmt.Errorf("%v bucket doesn't exist!", bucketName)
		}
		return bucket.Put(b(key), jsonJob)
	})
}

func ViewMeshFile(key string) (path string, err error) {
	err = DB.View(func(tx *bolt.Tx) error {
		path = string(tx.Bucket(b(dbMeshFiles)).Get(b(key)))
		return nil
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

func ViewGCodeFile(key string) (val string, err error) {
	err = DB.View(func(tx *bolt.Tx) error {
		val = string(tx.Bucket(b(dbGCodeFiles)).Get(b(key)))
		return nil
	})
	if err != nil {
		return "", err
	}
	return val, nil
}

func ViewJob(key string) (*slicerjob.Job, error) {
	var job = new(slicerjob.Job)
	err := DB.View(func(tx *bolt.Tx) error {
		jsonJob := tx.Bucket(b(dbJobs)).Get(b(key))
		return json.Unmarshal(jsonJob, job)
	})
	return job, err
}

func CancelJob(id string) error {
	job, err := ViewJob(id)
	if err != nil {
		return err
	}
	job.Status = slicerjob.Cancelled
	return PutJob(id, job)
}

func DeleteJob(id string) error {
	bucket := "jobs"
	err := DB.View(func(tx *bolt.Tx) error {
		err := tx.Bucket(b(bucket)).Delete(b(id))
		return err
	})
	return err
}

func DeleteGCodeFile(id string) error {
	bucket := "gCodeFiles"
	err := DB.View(func(tx *bolt.Tx) error {
		err := tx.Bucket(b(bucket)).Delete(b(id))
		return err
	})
	return err
}
