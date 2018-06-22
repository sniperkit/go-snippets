package db

import (
	"os"
	"os/user"
	"time"

	"github.com/boltdb/bolt"
)

type Database struct {
	Instance         *bolt.DB
	Bucket           []byte
	Datafile         string
	CurrentDirectory string
}

// Initalize the database and store the instance in the struct.
func Init() (d *Database, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	d = &Database{}
	d.Datafile = u.HomeDir + "/.config/dtags/dtags.db"
	d.Bucket = []byte("dtags")
	d.CurrentDirectory, _ = os.Getwd()

	// make all the directorys
	if err = os.MkdirAll(u.HomeDir+"/.config/dtags", 0755); err != nil {
		return nil, err
	}

	// Open the .db data file
	// It will be created if it doesn't exist.
	d.Instance, err = bolt.Open(d.Datafile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = d.Instance.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(d.Bucket)
		return err
	})

	return d, err
}

// Add a new key to the database with a given path.
func (d *Database) AddKey(k string, v string) error {
	return d.Instance.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(d.Bucket).Put([]byte(k), []byte(v))
	})
}

// Delete a key from the database.
func (d *Database) DeleteKey(k string) error {
	return d.Instance.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(d.Bucket).Delete([]byte(k))
	})
}

// Get all of the tags that currently exist in the database.
func (d *Database) GetTags() []string {
	var tags []string
	for key := range d.All() {
		tags = append(tags, key)
	}

	return tags
}

// Get the current tags that belong to the directory the command
// is being run in.
func (d *Database) GetCurrentTags() []string {
	var tags []string
	for key, value := range d.All() {
		if value == d.CurrentDirectory {
			tags = append(tags, key)
		}
	}

	return tags
}

// Get the value of a specific tag stored in the database.
func (d *Database) GetValue(k string) string {
	var val []byte
	_ = d.Instance.View(func(tx *bolt.Tx) error {
		val = tx.Bucket(d.Bucket).Get([]byte(k))

		return nil
	})

	return string(val)
}

// Return all of the tags currently stored in the database with their
// associated values.
func (d *Database) All() map[string]string {
	m := make(map[string]string)
	_ = d.Instance.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(d.Bucket).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			m[string(k)] = string(v)
		}

		return nil
	})

	return m
}

// Check if a given key exists in the database.
func (d *Database) Exists(k string) bool {
	return d.GetValue(k) != ""
}
