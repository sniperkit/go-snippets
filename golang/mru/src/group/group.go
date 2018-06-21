package group

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"feature"
	"mcase"
	"sort"
	"subgroup"
)

type Group struct {
	Name      string
	SGCount   int
	CCount    int
	ID        string
	SubGroups map[string]*subgroup.SubGroup
}

func Hash(name []byte) []byte {
	hash := sha1.New()
	return []byte(hex.EncodeToString(hash.Sum([]byte("groupGROUP" + string(name)))))
}

func (g *Group) Add(c *mcase.Case) error {
	sg, ok := g.SubGroups[c.SubGroup]
	if !ok {
		g.SubGroups[c.SubGroup] = &subgroup.SubGroup{
			Group:    g.Name,
			Name:     c.SubGroup,
			ID:       string(subgroup.Hash([]byte(c.SubGroup))),
			Features: make(map[string]*feature.Feature, 1),
		}
		g.SGCount++
		sg, _ = g.SubGroups[c.SubGroup]
	}

	err := sg.Add(c)
	if err != nil {
		return err
	}

	g.CCount++

	return nil
}

func (g *Group) Del(c *mcase.Case) error {
	sg, ok := g.SubGroups[c.SubGroup]
	if !ok {
		return errors.New("Cannot find Feature: " + c.Feature + " in Group: " + c.Group + " for delete case: " + c.Name)
	}

	err := sg.Del(c)
	if err != nil {
		return err
	}

	if len(sg.Features) == 0 {
		delete(g.SubGroups, c.SubGroup)
		g.SGCount--
	}

	g.CCount--
	return nil
}

func (g *Group) Get(c *mcase.Case) (*mcase.Case, error) {
	sg, ok := g.SubGroups[c.SubGroup]
	if !ok {
		return nil, errors.New("Cannot find Feature: " + c.Feature + " in Group: " + c.Group + " for Get case: " + c.Name)
	}

	return sg.Get(c)
}

func (g *Group) Dump() []*mcase.Case {
	result := make([]*mcase.Case, 0, 10)
	sgs := make([]*subgroup.SubGroup, 0, len(g.SubGroups))

	for _, sg := range g.SubGroups {
		sgs = append(sgs, sg)
	}

	//sort.Slice(fs, func(i, j int) bool { return fs[i].Name < fs[j].Name })
	sort.Stable(SubGroupSlice(sgs))
	for _, sg := range sgs {
		result = append(result, sg.Dump()...)
	}

	return result
}

func (g *Group) DumpSubGroup(sgroup string) ([]*mcase.Case, error) {
	sg, ok := g.SubGroups[sgroup]
	if !ok {
		return nil, errors.New("Cannot find SubGroup: " + sgroup + " for dump")
	}

	return sg.Dump(), nil
}

func (g *Group) DumpFeature(sgroup, feature string) ([]*mcase.Case, error) {
	sg, ok := g.SubGroups[sgroup]
	if !ok {
		return nil, errors.New("Cannot find SubGroup: " + sgroup + " for dump")
	}

	return sg.DumpFeature(feature)
}

type SubGroupSlice []*subgroup.SubGroup

func (s SubGroupSlice) Len() int           { return len(s) }
func (s SubGroupSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SubGroupSlice) Less(i, j int) bool { return s[i].Name < s[j].Name }
