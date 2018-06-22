package instance

import (
	"gopkg.in/mgo.v2/bson"
	"robostats/errmsg"
	"robostats/storage"
	"time"
	"upper.io/db"
)

type Instance struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectId `bson:"user_id" json:"user_id"`
	ClassID   bson.ObjectId `bson:"class_id" json:"class_id"`
	Data      interface{}   `bson:"data" json:"user_data"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

const (
	instanceCollectionName = "instance"
)

var (
	// InstanceCollection is the actual storage reference for instances.
	InstanceCollection db.Collection
)

func init() {
	InstanceCollection = storage.C(instanceCollectionName)
}

// GetByID returns a instance by the given ID.
func GetByID(id bson.ObjectId) (*Instance, error) {
	var err error
	var i Instance

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := InstanceCollection.Find(db.Cond{
		"_id": id,
	})

	if k, _ := res.Count(); k < 1 {
		return nil, errmsg.ErrNoSuchItem
	}

	if err = res.One(&i); err != nil {
		return nil, err
	}

	return &i, err
}

// GetByClassID returns a instance by the given ID.
func GetByClassID(id bson.ObjectId) ([]*Instance, error) {
	var err error
	var c []*Instance

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := InstanceCollection.Find(db.Cond{
		"class_id": id,
	})

	if k, _ := res.Count(); k < 1 {
		return nil, errmsg.ErrNoSuchItem
	}

	if err = res.All(&c); err != nil {
		return nil, err
	}

	return c, err
}

// GetByUserID returns a instance by the given ID.
func GetByUserID(id bson.ObjectId) ([]*Instance, error) {
	var err error
	var c []*Instance

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := InstanceCollection.Find(db.Cond{
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

func (i *Instance) Remove() error {

	if i.ID.Valid() == false {
		return errmsg.ErrInvalidID
	}

	res := InstanceCollection.Find(db.Cond{
		"_id": i.ID,
	})

	if k, _ := res.Count(); k < 1 {
		return errmsg.ErrNoSuchItem
	}

	return res.Remove()
}

// Update commits changes to permanent storage.
func (i *Instance) Update() error {
	var err error

	if i.ID.Valid() == false {
		return errmsg.ErrInvalidID
	}

	if err = i.save(); err != nil {
		return err
	}

	return nil
}

// Create adds a new instance to the database.
func (i *Instance) Create() error {
	var err error
	i.ID = bson.ObjectId("")

	if err = i.save(); err != nil {
		return nil
	}

	return nil
}

// save updates or appends a instance.
func (i *Instance) save() error {

	if i.ID.Valid() {
		res := InstanceCollection.Find(db.Cond{
			"_id": i.ID,
		})
		return res.Update(i)
	}

	i.CreatedAt = time.Now()

	id, err := InstanceCollection.Append(i)

	if err != nil {
		return err
	}

	i.ID = id.(bson.ObjectId)

	return nil
}
