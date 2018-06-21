package casedb

import (
	"github.com/Workiva/go-datastructures/trie/ctrie"
	"log"
)

type CaseDB struct {
	Device string
	DB     *ctrie.Ctrie
}

func (cdb *CaseDB) New(device string) *CaseDB {
	return &CaseDB{
		Device: device,
		DB:     ctrie.New(nil),
	}
}

func (cdb *CaseDB) Add(c *Case) error {
	if cdb.IsExist(c) {
		return errors.New("Same case alread exist")
	}

	cdb.Insert(cdb.Key(c), c)
}

func (cdb *CaseDB) Del(c *Case) error {
	if cdb.IsExist(c) {
		return errors.New("Same case alread exist")
	}

	cdb.Remove(cdb.Key(c))
}

func (cdb *CaseDB) Get(c *Case) (*Case, error) {
	if cdb.IsCaseValid(c) {
		if i, ok := cdb.Lookup(cdb.Key(c)); ok {
			if r, ok := i.(*Case); ok {
				return r, nil
			}
			log.Println("Invalid result when Get case from DB")
		}
	}

	return nil, errors.New("Cannot find case")
}

func (cdb *CaseDB) Run(c *Case) error {

}

func (cdb *CaseDB) IsExist(c *Case) error {
	if cdb.isCaseValid(c) {
		if _, ok := cdb.Lookup([]byte(cdb.Key(c))); ok {
			return true
		}
	}

	return false
}

func (cdb *CaseDB) isCaseValid(c *Case) bool {
	if c.Group == "" {
		return false
	}

	return true
}

func (cdb *CaseDB) Key(c *Case) []byte {
	var key = ""
	if cdb.isCaseValid(c) {
		key += c.Group
		if c.Feature != "" {
			key += c.Feature
		} else {
			return []byte(key)
		}

		if c.Name != "" {
			key += c.Name
		}
	}

	return []byte(key)
}

func (cdb *CaseDB) GetAll() <-chan *Case {
	ch := make(chan *Case)
	go func(in chan<- *Case) {
		for i := range cdb.Iterator(nil) {
			if c, ok := i.(*Case); ok {
				ch <- c
			}
		}
		defer close(ch)
	}(ch)

	return ch
}
