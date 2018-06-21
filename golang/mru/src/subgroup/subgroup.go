package subgroup

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"feature"
	"mcase"
	"sort"
)

type SubGroup struct {
	Group    string
	Name     string
	FCount   int
	CCount   int
	ID       string
	Features map[string]*feature.Feature
}

func Hash(name []byte) []byte {
	hash := sha1.New()
	return []byte(hex.EncodeToString(hash.Sum([]byte("subgroupSUBGROUP" + string(name)))))
}

func (sg *SubGroup) Add(c *mcase.Case) error {
	f, ok := sg.Features[c.Feature]
	if !ok {
		sg.Features[c.Feature] = &feature.Feature{
			Group:    sg.Group,
			SubGroup: sg.Name,
			Name:     c.Feature,
			ID:       string(feature.Hash([]byte(c.Feature))),
			Cases:    make(map[string]*mcase.Case, 1),
		}
		sg.FCount++
		f, _ = sg.Features[c.Feature]
	}

	err := f.Add(c)
	if err != nil {
		return err
	}

	sg.CCount++

	return nil
}

func (sg *SubGroup) Del(c *mcase.Case) error {
	f, ok := sg.Features[c.Feature]
	if !ok {
		return errors.New("Cannot find Feature: " + c.Feature + " in SubGroup: " + c.SubGroup + " for delete case: " + c.Name)
	}

	err := f.Del(c)
	if err != nil {
		return err
	}

	if len(f.Cases) == 0 {
		delete(sg.Features, c.Feature)
		sg.FCount--
	}

	sg.CCount--

	return nil
}

func (sg *SubGroup) Get(c *mcase.Case) (*mcase.Case, error) {
	f, ok := sg.Features[c.Feature]
	if !ok {
		return nil, errors.New("Cannot find Feature: " + c.Feature + " in SubGroup: " + c.SubGroup + " for Get case: " + c.Name)
	}

	return f.Get(c)
}

func (sg *SubGroup) Dump() []*mcase.Case {
	result := make([]*mcase.Case, 0, 10)
	fs := make([]*feature.Feature, 0, len(sg.Features))

	for _, f := range sg.Features {
		fs = append(fs, f)
	}

	//sort.Slice(fs, func(i, j int) bool { return fs[i].Name < fs[j].Name })
	sort.Stable(FeatureSlice(fs))
	for _, f := range fs {
		result = append(result, f.Dump()...)
	}

	return result
}

func (sg *SubGroup) DumpFeature(feature string) ([]*mcase.Case, error) {
	f, ok := sg.Features[feature]
	if !ok {
		return nil, errors.New("Cannot find Group: " + feature + " for dump")
	}

	return f.Dump(), nil
}

type FeatureSlice []*feature.Feature

func (s FeatureSlice) Len() int           { return len(s) }
func (s FeatureSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s FeatureSlice) Less(i, j int) bool { return s[i].Name < s[j].Name }
