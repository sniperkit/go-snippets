package instance

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var (
	fakeUserID     bson.ObjectId
	fakeInstanceID bson.ObjectId
	fakeClassID    bson.ObjectId
)

func init() {
	fakeUserID = bson.NewObjectId()
	fakeClassID = bson.NewObjectId()
}

func TestInstanceCreate(t *testing.T) {
	c := Instance{
		UserID:  fakeUserID,
		ClassID: fakeClassID,
	}

	if err := c.Create(); err != nil {
		t.Fatal(err)
	}

	if c.ID.Valid() == false {
		t.Fatal("Expecting a valid ID.")
	}

	fakeInstanceID = c.ID
}

func TestInstanceEdit(t *testing.T) {
	c, err := GetByID(fakeInstanceID)

	if err != nil {
		t.Fatal(err)
	}

	if c.ID != fakeInstanceID {
		t.Fatal("Expecting a valid ID.")
	}

	if c.UserID != fakeUserID {
		t.Fatal("Expecting a valid ID.")
	}

	if err := c.Update(); err != nil {
		t.Fatal(err)
	}
}

func TestInstanceRemove(t *testing.T) {
	c, err := GetByID(fakeInstanceID)

	if err != nil {
		t.Fatal(err)
	}

	if err = c.Remove(); err != nil {
		t.Fatal(err)
	}

	if _, err = GetByID(fakeInstanceID); err == nil {
		t.Fatal("Expecting an error.")
	}

}
