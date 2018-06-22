package session

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var (
	fakeUserID     bson.ObjectId
	fakeSessionID  bson.ObjectId
	fakeClassID    bson.ObjectId
	fakeInstanceID bson.ObjectId
)

func init() {
	fakeUserID = bson.NewObjectId()
	fakeClassID = bson.NewObjectId()
	fakeInstanceID = bson.NewObjectId()
}

func TestSessionCreate(t *testing.T) {
	c := Session{
		UserID:     fakeUserID,
		ClassID:    fakeClassID,
		InstanceID: fakeInstanceID,
	}

	if err := c.Create(); err != nil {
		t.Fatal(err)
	}

	if c.ID.Valid() == false {
		t.Fatal("Expecting a valid ID.")
	}

	fakeSessionID = c.ID
}

func TestSessionEdit(t *testing.T) {
	c, err := GetByID(fakeSessionID)

	if err != nil {
		t.Fatal(err)
	}

	if c.ID != fakeSessionID {
		t.Fatal("Expecting a valid ID.")
	}

	if c.UserID != fakeUserID {
		t.Fatal("Expecting a valid ID.")
	}

	if c.InstanceID != fakeInstanceID {
		t.Fatal("Expecting a valid ID.")
	}

	if err := c.Update(); err != nil {
		t.Fatal(err)
	}
}

func TestSessionRemove(t *testing.T) {
	c, err := GetByID(fakeSessionID)

	if err != nil {
		t.Fatal(err)
	}

	if err = c.Remove(); err != nil {
		t.Fatal(err)
	}

	if _, err = GetByID(fakeSessionID); err == nil {
		t.Fatal("Expecting an error.")
	}

}
