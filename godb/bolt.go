package godb

import (
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

var boltDB *bolt.DB
var boltDefaultBucketName = "default"
var boltBucket *bolt.Bucket

func init() {
	var err error
	boltDB, err = bolt.Open("my.db", 0644, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("init int bolt...")

	err = boltDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(boltDefaultBucketName))
		if err != nil {
			return err
		}
		log.Printf("successfully create bucket - [%s]\n", boltDefaultBucketName)

		boltBucket = b
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}

// add a <key, value> pair to bolt
// only bytes are allowed as key and value
// _metadata.db is database name under directory
// bucket is 'default' for simplicity.
func addBoltKV(key string, value []byte) error {
	err := boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("default"))
		if b == nil {
			return errors.New("can't get [default] bucket in boltdb")
		}
		err := b.Put([]byte(key), value)
		return err
	})

	return err
}

func deleteBoltKV(key string) error {
	err := boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("default"))
		if b == nil {
			return errors.New("can't get [default] bucket in boltdb")
		}
		err := b.Delete([]byte(key))
		return err
	})

	return err
}

func getBoltKV(key string) []byte {
	var val []byte
	boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("default"))
		if b == nil {
			return errors.New("can't get [default] bucket in boltdb")
		}
		value := b.Get([]byte(key))
		// Must use copy to store the value, because it is just valid in this transaction
		// copy(val, value)
		if len(value) == 0 {
			val = nil
		} else {
			val = make([]byte, len(value))
			copy(val, value)
		}
		return nil
	})

	return val
}

func getAllBoltKV() map[string][]byte {
	m := make(map[string][]byte)

	boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("default"))
		if b == nil {
			return errors.New("can't get [default] bucket in boltdb")
		}
		b.ForEach(func(k, v []byte) error {
			log.Println(string(k))
			// must copy key, value to out varialbe (out of this transaction scope)
			// to use the value, or creat a new one and use it immediately
			if len(v) == 0 {
				m[string(k)] = nil
			} else {
				value := make([]byte, len(v))
				copy(value, v)
				m[string(k)] = value
			}

			return nil
		})

		return nil
	})

	return m
}

// EmptyBoltBucket -- empty a bucket
// if not exist, create one
func EmptyBoltBucket(name string) error {

	err := boltDB.Update(func(tx *bolt.Tx) error {
		namebytes := []byte(name)
		bucket := tx.Bucket(namebytes)
		if bucket != nil {
			err := tx.DeleteBucket(namebytes)
			if err != nil {
				log.Printf("empty bucket - [%s] error\n", name)
				return err
			}
		}

		_, err := tx.CreateBucket(namebytes)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
