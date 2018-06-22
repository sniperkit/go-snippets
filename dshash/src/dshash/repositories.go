package dshash

import (
	"appengine"
	"appengine/datastore"
	"errors"
)

type PersonRepository struct {
	appengine.Context
}

func (pr *PersonRepository) Find(person *Person) (*Person, error) {
	q := datastore.NewQuery("Person").Filter("Handler =", person.Handler)

	foundPersons := []Person{}
	_, err := q.GetAll(pr.Context, &foundPersons)

	if err != nil {
		return nil, err
	}

	if len(foundPersons) > 1 {
		return nil, errors.New("user stored multiple times")
	} else if len(foundPersons) < 1 {
		return nil, nil
	}

	return &foundPersons[0], nil
}

func (pr *PersonRepository) Save(person *Person) error {
	key := datastore.NewIncompleteKey(pr.Context, "Person", nil)

	_, err := datastore.Put(pr.Context, key, person)

	return err
}
