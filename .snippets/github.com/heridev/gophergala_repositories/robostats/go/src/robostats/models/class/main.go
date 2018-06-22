package class

import (
	"gopkg.in/mgo.v2/bson"
	"robostats/errmsg"
	"robostats/storage"
	"time"
	"upper.io/db"
	"upper.io/i/v1/session/tokener"
)

type Class struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectId `bson:"user_id" json:"user_id"`
	Name      string        `bson:"name" json:"name"`
	APIKey    string        `bson:"api_key" json:"api_key"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

const (
	classCollectionName = "class"
)

var (
	// ClassCollection is the actual storage reference for classs.
	ClassCollection db.Collection
)

func init() {
	ClassCollection = storage.C(classCollectionName)
}

// GetByUserID returns a class by the given ID.
func GetByUserID(id bson.ObjectId) ([]*Class, error) {
	var err error
	var c []*Class

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := ClassCollection.Find(db.Cond{
		"user_id": id,
	})

	if k, _ := res.Count(); k < 1 {
		return nil, errmsg.ErrNoSuchItem
	}

	if err = res.All(&c); err != nil {
		return nil, err
	}

	return c, err
}

// GetByID returns a class by the given ID.
func GetByID(id bson.ObjectId) (*Class, error) {
	var err error
	var c Class

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := ClassCollection.Find(db.Cond{
		"_id": id,
	})

	if k, _ := res.Count(); k < 1 {
		return nil, errmsg.ErrNoSuchItem
	}

	if err = res.One(&c); err != nil {
		return nil, err
	}

	return &c, err
}

func (c *Class) Remove() error {

	if c.ID.Valid() == false {
		return errmsg.ErrInvalidID
	}

	res := ClassCollection.Find(db.Cond{
		"_id": c.ID,
	})

	if k, _ := res.Count(); k < 1 {
		return errmsg.ErrNoSuchItem
	}

	return res.Remove()
}

// Update commits changes to permanent storage.
func (c *Class) Update() error {
	var err error

	if c.ID.Valid() == false {
		return errmsg.ErrInvalidID
	}

	if err = c.save(); err != nil {
		return err
	}

	return nil
}

// Create adds a new class to the database.
func (c *Class) Create() error {
	var err error
	c.ID = bson.ObjectId("")

	c.APIKey = tokener.String(40)

	if err = c.save(); err != nil {
		return nil
	}

	return nil
}

// save updates or appends a class.
func (c *Class) save() error {

	if c.ID.Valid() {
		res := ClassCollection.Find(db.Cond{
			"_id": c.ID,
		})
		return res.Update(c)
	}

	c.CreatedAt = time.Now()

	id, err := ClassCollection.Append(c)

	if err != nil {
		return err
	}

	c.ID = id.(bson.ObjectId)

	return nil
}
