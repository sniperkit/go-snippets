package session

import (
	"gopkg.in/mgo.v2/bson"
	"robostats/errmsg"
	"robostats/storage"
	"time"
	"upper.io/db"
	"upper.io/i/v1/session/tokener"
)

type Session struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserID     bson.ObjectId `bson:"user_id" json:"user_id"`
	ClassID    bson.ObjectId `bson:"class_id" json:"class_id"`
	InstanceID bson.ObjectId `bson:"instance_id" json:"instance_id"`
	SessionKey string        `bson:"session_key" json:"session_key"`
	Data       interface{}   `bson:"data" json:"user_data"`
	StartTime  time.Time     `bson:"start_time,omitempty" json:"start_time"`
	EndTime    time.Time     `bson:"end_time,omitempty" json:"end_time"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
}

const (
	sessionCollectionName = "session"
)

var (
	// SessionCollection is the actual storage reference for sessions.
	SessionCollection db.Collection
)

func init() {
	SessionCollection = storage.C(sessionCollectionName)
}

// GetByID returns a session by the given ID.
func GetByID(id bson.ObjectId) (*Session, error) {
	var err error
	var s Session

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := SessionCollection.Find(db.Cond{
		"_id": id,
	})

	if k, _ := res.Count(); k < 1 {
		return nil, errmsg.ErrNoSuchItem
	}

	if err = res.One(&s); err != nil {
		return nil, err
	}

	return &s, err
}

// GetByInstanceID returns sessions associated with the given Instance ID.
func GetByInstanceID(id bson.ObjectId) ([]*Session, error) {
	var err error
	var c []*Session

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := SessionCollection.Find(db.Cond{
		"instance_id": id,
	})

	if k, _ := res.Count(); k < 1 {
		return nil, errmsg.ErrNoSuchItem
	}

	if err = res.All(&c); err != nil {
		return nil, err
	}

	return c, err
}

// GetByUserID returns sessions associated with the given User ID.
func GetByUserID(id bson.ObjectId) ([]*Session, error) {
	var err error
	var c []*Session

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := SessionCollection.Find(db.Cond{
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

// GetByClassID returns a instance by the given ID.
func GetByClassID(id bson.ObjectId) ([]*Session, error) {
	var err error
	var c []*Session

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := SessionCollection.Find(db.Cond{
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

func (s *Session) Remove() error {

	if s.ID.Valid() == false {
		return errmsg.ErrInvalidID
	}

	res := SessionCollection.Find(db.Cond{
		"_id": s.ID,
	})

	if k, _ := res.Count(); k < 1 {
		return errmsg.ErrNoSuchItem
	}

	return res.Remove()
}

// Update commits changes to permanent storage.
func (s *Session) Update() error {
	var err error

	if s.ID.Valid() == false {
		return errmsg.ErrInvalidID
	}

	if err = s.save(); err != nil {
		return err
	}

	return nil
}

// Create adds a new session to the database.
func (s *Session) Create() error {
	var err error
	s.ID = bson.ObjectId("")

	s.SessionKey = tokener.String(20)

	if err = s.save(); err != nil {
		return nil
	}

	return nil
}

// save updates or appends a session.
func (s *Session) save() error {

	if s.ID.Valid() {
		res := SessionCollection.Find(db.Cond{
			"_id": s.ID,
		})
		return res.Update(s)
	}

	s.CreatedAt = time.Now()

	id, err := SessionCollection.Append(s)

	if err != nil {
		return err
	}

	s.ID = id.(bson.ObjectId)

	return nil
}
