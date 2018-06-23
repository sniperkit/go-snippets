package feature

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"mcase"
	"rut"
	"sort"
)

type Feature struct {
	Group    string
	SubGroup string
	Name     string
	CCount   int
	ID       string
	RUTDB    map[string]*rut.RUT
	DUTCount int
	Cases    map[string]*mcase.Case
}

func Hash(name []byte) []byte {
	hash := sha1.New()
	return []byte(hex.EncodeToString(hash.Sum([]byte("featureFEATURE" + string(name)))))
}

func (f *Feature) Add(c *mcase.Case) error {
	_, ok := f.Cases[c.Name]
	if ok {
		return errors.New("Case: " + c.Name + "Alread exist for Group: " + c.Group + " Feature: " + c.Feature)
	}

	f.CCount++
	f.Cases[c.Name] = c

	return nil
}

func (f *Feature) Del(c *mcase.Case) error {
	_, ok := f.Cases[c.Name]
	if !ok {
		return errors.New("Cannot find  Case: " + c.Name + "for delete under Group: " + c.Group + " Feature: " + c.Feature)
	}

	delete(f.Cases, c.Name)
	f.CCount--
	return nil
}

func (f *Feature) Get(c *mcase.Case) (*mcase.Case, error) {
	old, ok := f.Cases[c.Name]
	if !ok {
		return nil, errors.New("Cannot find  Case: " + c.Name + "for Get under Group: " + c.Group + " Feature: " + c.Feature)
	}

	return old, nil
}

func (f *Feature) Dump() []*mcase.Case {
	cs := make([]*mcase.Case, 0, len(f.Cases))

	for _, c := range f.Cases {
		cs = append(cs, c)
	}

	//sort.Slice(cs, func(i, j int) bool { return cs[i].Name < cs[j].Name })
	sort.Stable(CaseSlice(cs))

	return cs
}

func (f *Feature) AddDUT(r *rut.RUT) error {
	if f.RUTDB == nil {
		f.RUTDB = make(map[string]*rut.RUT, 1)
	}
	if o, ok := f.RUTDB[r.Name]; ok {
		return fmt.Errorf("Same rut: %s already exist: %v", r.Name, o)
	}

	for _, c := range f.Cases {
		c.AddRUT(r)
	}

	f.RUTDB[r.Name] = r

	return nil
}

func (f *Feature) DelDUT(r *rut.RUT) error {
	if _, ok := f.RUTDB[r.Name]; !ok {
		return fmt.Errorf("DUT: %s does not exist", r.Name)
	}

	for _, c := range f.Cases {
		c.DelRUT(r)
	}
	delete(f.RUTDB, r.Name)
	return nil
}

func (f *Feature) ClearDUTs() {
	if f.RUTDB == nil {
		log.Println("Call clear RUT when there is no rut in feature DB")
	}

	for _, r := range f.RUTDB {
		for _, c := range f.Cases {
			c.DelRUT(r)
		}
		delete(f.RUTDB, r.Name)
	}
}

type CaseSlice []*mcase.Case

func (s CaseSlice) Len() int           { return len(s) }
func (s CaseSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s CaseSlice) Less(i, j int) bool { return s[i].Name < s[j].Name }
