package beat

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var (
	fakeUserID     bson.ObjectId
	fakeBeatID     bson.ObjectId
	fakeClassID    bson.ObjectId
	fakeInstanceID bson.ObjectId
	fakeSessionID  bson.ObjectId
)

func init() {
	fakeUserID = bson.NewObjectId()
	fakeClassID = bson.NewObjectId()
	fakeInstanceID = bson.NewObjectId()
	fakeSessionID = bson.NewObjectId()
}

func TestBeatCreate(t *testing.T) {
	c := Beat{
		UserID:     fakeUserID,
		ClassID:    fakeClassID,
		InstanceID: fakeInstanceID,
		SessionID:  fakeSessionID,
	}

	if err := c.Create(); err != nil {
		t.Fatal(err)
	}

	if c.ID.Valid() == false {
		t.Fatal("Expecting a valid ID.")
	}

	fakeBeatID = c.ID
}

func TestBeatEdit(t *testing.T) {
	c, err := GetByID(fakeBeatID)

	if err != nil {
		t.Fatal(err)
	}

	if c.ID != fakeBeatID {
		t.Fatal("Expecting a valid ID.")
	}

	if c.UserID != fakeUserID {
		t.Fatal("Expecting a valid ID.")
	}

	if c.InstanceID != fakeInstanceID {
		t.Fatal("Expecting a valid ID.")
	}

	if c.SessionID != fakeSessionID {
		t.Fatal("Expecting a valid ID.")
	}

	if err := c.Update(); err != nil {
		t.Fatal(err)
	}
}

func TestBeatRemove(t *testing.T) {
	c, err := GetByID(fakeBeatID)

	if err != nil {
		t.Fatal(err)
	}

	if err = c.Remove(); err != nil {
		t.Fatal(err)
	}

	if _, err = GetByID(fakeBeatID); err == nil {
		t.Fatal("Expecting an error.")
	}

}
