package datastore

import (
	"errors"
	"reflect"
	"testing"
)

type mockObject struct {
	Name      string
	Attribute string
}

func (m mockObject) Key() string {
	return m.Name
}

func (m mockObject) Type() string {
	return "MockObject"
}

type mockObject2 struct {
	Name      string
	Attribute string
}

func (m mockObject2) Key() string {
	return m.Name
}

func (m mockObject2) Type() string {
	return "MockObject2"
}

func (m mockObject) Validate(context Context) error {
	if len(m.Name) == 0 {
		return errors.New("mockObject must have a name")
	}
	nomore := mockObject{Name: "nomore"}
	if err := Get(context, &nomore); err != ErrNotFound {
		return errors.New("nomore objects allowed because nomore object found")
	}
	return nil
}

func (m mockObject2) Validate(context Context) error {
	if len(m.Name) == 0 {
		return errors.New("mockObject2 must have a name")
	}
	nomore := mockObject2{Name: "nomore"}
	if err := Get(context, &nomore); err != ErrNotFound {
		return errors.New("nomore objects allowed because nomore object found")
	}
	return nil
}

func TestStores(t *testing.T) {
	storageContext, _ := NewMemoryContext()

	ob1 := mockObject{
		Name:      "test",
		Attribute: "foo",
	}

	if exists, _ := Exists(storageContext, ob1); exists {
		t.Fatalf("Did not expect to find %v", ob1)
	}

	if err := Put(storageContext, ob1); err != nil {
		t.Fatalf("Could not store object %s: %s", ob1, err)
	}

	ob2 := mockObject{Name: "test"}
	if err := Get(storageContext, &ob2); err != nil {
		t.Fatalf("Could not get object %s: %s", ob1, err)
	}

	if !reflect.DeepEqual(ob1, ob2) {
		t.Fatalf("ob1 '%v' != ob2 '%v'", ob1, ob2)
	}

	notfindable := mockObject{Name: "does not exist"}
	if err := Get(storageContext, notfindable); err == nil {
		t.Fatalf("Should not have not found %s", notfindable)
	}

	if exists, _ := Exists(storageContext, notfindable); exists {
		t.Fatalf("Did not expect to find %v", ob1)
	}

	if err := Get(storageContext, mockObject2{Name: "does not exist"}); err == nil {
		t.Fatalf("Should not have not found %s", notfindable)
	}

	if err := Delete(storageContext, mockObject2{Name: "does not exist"}); err != nil {
		t.Fatalf("Unexexpected error when deleting non-existent object %s: %s", notfindable, err)
	}

	nomore := mockObject{Name: "nomore"}
	if err := Put(storageContext, nomore); err != nil {
		t.Fatalf("Could not add 'nomore' object", err)
	}

	if err := Put(storageContext, ob1); err == nil {
		t.Fatalf("expected error after nomore object")
	}

	if err := Delete(storageContext, ob1); err != nil {
		t.Fatalf("Unexexpected error when deleting object %s: %s", ob1, err)
	}
}
