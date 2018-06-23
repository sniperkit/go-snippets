package class

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var (
	fakeUserID  bson.ObjectId
	fakeClassID bson.ObjectId
)

func init() {
	fakeUserID = bson.NewObjectId()
}

func TestClassCreate(t *testing.T) {
	c := Class{
		UserID: fakeUserID,
		Name:   "Drones",
	}

	if err := c.Create(); err != nil {
		t.Fatal(err)
	}

	if c.APIKey == "" {
		t.Fatal("Expecting API key.")
	}

	if c.ID.Valid() == false {
		t.Fatal("Expecting a valid ID.")
	}

	fakeClassID = c.ID
}

func TestClassEdit(t *testing.T) {
	c, err := GetByID(fakeClassID)

	if err != nil {
		t.Fatal(err)
	}

	if c.ID != fakeClassID {
		t.Fatal("Expecting a valid ID.")
	}

	if c.UserID != fakeUserID {
		t.Fatal("Expecting a valid ID.")
	}

	c.Name = "Barooza"

	if err := c.Update(); err != nil {
		t.Fatal(err)
	}
}

func TestClassRemove(t *testing.T) {
	c, err := GetByID(fakeClassID)

	if err != nil {
		t.Fatal(err)
	}

	if c.Name != "Barooza" {
		t.Fatal("Expecting an actual modification.")
	}

	if err = c.Remove(); err != nil {
		t.Fatal(err)
	}

	if _, err = GetByID(fakeClassID); err == nil {
		t.Fatal("Expecting an error.")
	}

}
