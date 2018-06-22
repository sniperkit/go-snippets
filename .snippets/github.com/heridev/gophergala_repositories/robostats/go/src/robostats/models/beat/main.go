package beat

import (
	"gopkg.in/mgo.v2/bson"
	labixbson "labix.org/v2/mgo/bson" // TODO: fill a mgo bug for this thing
	"robostats/errmsg"
	"robostats/storage"
	"time"
	"upper.io/db"
)

type Beat struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserID     bson.ObjectId `bson:"user_id" json:"user_id"`
	ClassID    bson.ObjectId `bson:"class_id" json:"class_id"`
	InstanceID bson.ObjectId `bson:"instance_id" json:"instance_id"`
	SessionID  bson.ObjectId `bson:"session_id" json:"session_id"`
	Data       interface{}   `bson:"data" json:"user_data"`
	LocalTime  int           `bson:"local_time" json:"local_time"`
	LatLng     [2]float64    `bson:"latlng" json:"latlng"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
}

const (
	beatCollectionName = "beat"
)

var (
	// BeatCollection is the actual storage reference for beats.
	BeatCollection db.Collection
)

func init() {
	BeatCollection = storage.C(beatCollectionName)
}

// GetByID returns a beat by the given ID.
func GetByID(id bson.ObjectId) (*Beat, error) {
	var err error
	var b Beat

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := BeatCollection.Find(db.Cond{
		"_id": id,
	})

	if k, _ := res.Count(); k < 1 {
		return nil, errmsg.ErrNoSuchItem
	}

	if err = res.One(&b); err != nil {
		return nil, err
	}

	return &b, err
}

// GetByClassID returns beats associated with the given Class ID.
func GetByClassID(id bson.ObjectId) ([]*Beat, error) {
	var err error
	var c []*Beat

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := BeatCollection.Find(db.Cond{
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

// GetBySessionID returns beats associated with the given Session ID.
func GetBySessionID(id bson.ObjectId) ([]*Beat, error) {
	var err error
	var c []*Beat

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := BeatCollection.Find(db.Cond{
		"session_id": id,
	}).Sort("local_time")

	if k, _ := res.Count(); k < 1 {
		return nil, errmsg.ErrNoSuchItem
	}

	if err = res.All(&c); err != nil {
		return nil, err
	}

	return c, err
}

// GetByUserID returns beats associated with the given User ID.
func GetByUserID(id bson.ObjectId) ([]*Beat, error) {
	var err error
	var c []*Beat

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := BeatCollection.Find(db.Cond{
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

// GetByInstanceID returns beats associated with the given Instance ID.
func GetByInstanceID(id bson.ObjectId) ([]*Beat, error) {
	var err error
	var c []*Beat

	if id.Valid() == false {
		return nil, errmsg.ErrInvalidID
	}

	res := BeatCollection.Find(db.Cond{
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

func (b *Beat) Remove() error {

	if b.ID.Valid() == false {
		return errmsg.ErrInvalidID
	}

	res := BeatCollection.Find(db.Cond{
		"_id": b.ID,
	})

	if k, _ := res.Count(); k < 1 {
		return errmsg.ErrNoSuchItem
	}

	return res.Remove()
}

// Update commits changes to permanent storage.
func (b *Beat) Update() error {
	var err error

	if b.ID.Valid() == false {
		return errmsg.ErrInvalidID
	}

	if err = b.save(); err != nil {
		return err
	}

	return nil
}

// Create adds a new beat to the database.
func (b *Beat) Create() error {
	var err error
	b.ID = bson.ObjectId("")

	if err = b.save(); err != nil {
		return nil
	}

	return nil
}

// save updates or appends a beat.
func (b *Beat) save() error {

	if b.ID.Valid() {
		res := BeatCollection.Find(db.Cond{
			"_id": b.ID,
		})
		return res.Update(b)
	}

	b.CreatedAt = time.Now()

	id, err := BeatCollection.Append(b)

	if err != nil {
		return err
	}

	b.ID = id.(bson.ObjectId)

	return nil
}

func (b *Beat) GetKeyValue(k string) (interface{}, error) {
	switch k {
	case "lat_lng":
		return b.LatLng, nil
	default:
		var ok bool
		var values labixbson.M

		if b.Data == nil {
			return nil, errmsg.ErrNoSuchKey
		}

		// WTF?
		if values, ok = b.Data.(labixbson.M); !ok {
			return nil, errmsg.ErrNoSuchKey
		}

		if _, ok = values[k]; ok {
			return values[k], nil
		}
	}
	return nil, errmsg.ErrNoSuchKey
}
