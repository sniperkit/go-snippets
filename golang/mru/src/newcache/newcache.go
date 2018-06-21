package newcache

import (
	"context"
	"encoding/json"
	"errors"
	"feature"
	"fmt"
	"github.com/Workiva/go-datastructures/trie/ctrie"
	"group"
	"io/ioutil"
	"log"
	"mcase"
	"result"
	"rut"
	"subgroup"
	"sync"
	"task"
	"util"
)

type NewCache struct {
	Name   string
	Node   *Node
	Groups map[string]*group.Group
	*ctrie.Ctrie
}

type Node struct {
	Key      string `json:"text"`
	ID       string
	Type     string
	Link     string  `json:"href"`
	Children []*Node `json:"nodes"`
	Data     interface{}
}

func New(name string) *NewCache {
	return &NewCache{
		Name: name,
		Node: &Node{
			Key:      name,
			Type:     "ROOT",
			Link:     "caseinfo",
			Children: make([]*Node, 0, 1),
		},
		Ctrie:  ctrie.New(nil),
		Groups: make(map[string]*group.Group, 1),
	}
}

func (tr *NewCache) Save() {
	cases := tr.GetAllCase()
	js, err := json.Marshal(cases)
	if err != nil {
		log.Println("Cannot format db for debug")
		return
	}

	//This is very stupid, but I'm really very tired.
	if tr.Name == util.GenerateSessionIDByUserNameAndPassword("mra", "") {
		util.SaveToFile("asset/db/common/COMMON.json", js)
	}
	util.SaveToFile("asset/db/personal/"+tr.Name+".json", js)
}

func (tr *NewCache) Restore() {
	tr.RestoreCaseFromFile("asset/db/common/COMMON.json")
	if tr.Name != util.GenerateSessionIDByUserNameAndPassword("mra", "") {
		tr.RestoreCaseFromFile("asset/db/personal/" + tr.Name + ".json")
	}
}

func (tr *NewCache) SetTopologyByID(id, name string, content []byte) error {
	//Just implement Set Case ID
	c, err := tr.GetCaseByID(id)
	if err != nil {
		return err
	}

	log.Printf("Save %s\n", name)
	c.SetTopology(name, content)

	defer tr.Save()

	return nil
}

func (tr *NewCache) RestoreCaseFromFile(name string) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		log.Println("Cannot read case db")
		return
	}

	cases := make([]*mcase.Case, 0, 1)
	err = json.Unmarshal(data, &cases)
	if err != nil {
		panic(err.Error())
	}

	for _, c := range cases {
		tr.RegisterCase(c)
	}
}

func (tr *NewCache) TreeView() *Node {
	return tr.Node
}

func (tr *NewCache) IsCaseValid(c *mcase.Case) bool {
	if c.Group == "" ||
		c.SubGroup == "" ||
		c.Feature == "" ||
		c.Name == "" ||
		c.DUTCount <= 0 {
		return false
	}
	return true
}

func (tr *NewCache) RegisterCase(c *mcase.Case) error {
	if !tr.IsCaseValid(c) {
		return fmt.Errorf("Case: is invalid")
	}

	if tr.isNodeExist(mcase.Hash(tr.CaseKey(c))) {
		return errors.New("Same case alread exist")
	}

	if !tr.isNodeExist(group.Hash(tr.GetGroupKeyByCase(c))) {
		newgroup := &group.Group{
			Name:      c.Group,
			ID:        string(group.Hash(tr.GetGroupKeyByCase(c))),
			SubGroups: make(map[string]*subgroup.SubGroup, 1),
		}
		newnode := &Node{
			Key:      c.Group,
			ID:       string(group.Hash(tr.GetGroupKeyByCase(c))),
			Type:     "GROUP",
			Link:     "groupinfo?id=" + string(group.Hash(tr.GetGroupKeyByCase(c))),
			Children: make([]*Node, 0, 1),
			Data:     newgroup,
		}

		tr.Groups[c.Group] = newgroup
		tr.Insert(group.Hash(tr.GetGroupKeyByCase(c)), newnode)
		tr.Node.Children = append(tr.Node.Children, newnode)
	}

	if !tr.isNodeExist(subgroup.Hash(tr.GetSubGroupKeyByCase(c))) {
		newsubgroup := &subgroup.SubGroup{
			Group:    c.Group,
			Name:     c.SubGroup,
			ID:       string(subgroup.Hash(tr.GetSubGroupKeyByCase(c))),
			Features: make(map[string]*feature.Feature, 1),
		}

		newnode := &Node{
			Key:      c.SubGroup,
			ID:       string(subgroup.Hash(tr.GetSubGroupKeyByCase(c))),
			Type:     "SUBGROUP",
			Link:     "subgroupinfo?id=" + string(subgroup.Hash(tr.GetSubGroupKeyByCase(c))),
			Children: make([]*Node, 0, 1),
			Data:     newsubgroup,
		}
		tr.Insert(subgroup.Hash(tr.GetSubGroupKeyByCase(c)), newnode)

		n, err := tr.Get(group.Hash(tr.GetGroupKeyByCase(c)))
		if err != nil {
			panic(err)
		}

		if g, ok := n.Data.(*group.Group); ok {
			g.SGCount++
			g.SubGroups[c.SubGroup] = newsubgroup
		} else {
			log.Printf("!!!!!!: %s is not a group\n", tr.GetGroupKeyByCase(c))
		}
		n.Children = append(n.Children, newnode)
	}

	if !tr.isNodeExist(feature.Hash(tr.GetFeatureKeyByCase(c))) {
		newfeature := &feature.Feature{
			Group:    c.Group,
			SubGroup: c.SubGroup,
			Name:     c.Feature,
			DUTCount: c.DUTCount,
			ID:       string(feature.Hash(tr.GetFeatureKeyByCase(c))),
			Cases:    make(map[string]*mcase.Case, 1),
		}

		newnode := &Node{
			Key:      c.Feature,
			ID:       string(feature.Hash(tr.GetFeatureKeyByCase(c))),
			Type:     "FEATURE",
			Link:     "featureinfo?id=" + string(feature.Hash(tr.GetFeatureKeyByCase(c))),
			Children: make([]*Node, 0, 1),
			Data:     newfeature,
		}
		tr.Insert(feature.Hash(tr.GetFeatureKeyByCase(c)), newnode)

		n, err := tr.Get(subgroup.Hash(tr.GetSubGroupKeyByCase(c)))
		if err != nil {
			panic(err)
		}
		if sg, ok := n.Data.(*subgroup.SubGroup); ok {
			sg.FCount++
			sg.Features[c.Feature] = newfeature
		} else {
			log.Printf("!!!!!!: %s is not a subgroup\n", tr.GetSubGroupKeyByCase(c))
		}

		n.Children = append(n.Children, newnode)
	} else { //I want to run all the case of a feature at the same time. So the DUT count must be the largest count of all the cases
		f, _ := tr.GetFeatureByID(string(feature.Hash(tr.GetFeatureKeyByCase(c))))
		if f.DUTCount < c.DUTCount {
			f.DUTCount = c.DUTCount
		}
	}

	//Case Node no need children.
	if !tr.isNodeExist(mcase.Hash(tr.CaseKey(c))) {
		if c.Tasks == nil {
			c.Tasks = make([]*task.Task, 0, 1)
		}
		c.ID = string(mcase.Hash(tr.CaseKey(c)))
		newnode := &Node{
			Key:  c.Name,
			ID:   string(mcase.Hash(tr.CaseKey(c))),
			Type: "CASE",
			Link: "caseinfo?id=" + string(mcase.Hash(tr.CaseKey(c))),
			Data: c,
		}
		tr.Insert(mcase.Hash(tr.CaseKey(c)), newnode)

		n, err := tr.Get(feature.Hash(tr.GetFeatureKeyByCase(c)))
		if err != nil {
			panic(err)
		}

		if f, ok := n.Data.(*feature.Feature); ok {
			f.CCount++
			f.Cases[c.Name] = c
		} else {
			log.Printf("!!!!!!: %s is not a feature\n", tr.GetFeatureKeyByCase(c))
		}

		n.Children = append(n.Children, newnode)
		return nil
	}

	return fmt.Errorf("Case: %s already exist!")
}

func (tr *NewCache) AddCase(c *mcase.Case) error {
	if err := tr.RegisterCase(c); err != nil {
		return err
	}
	tr.Save()

	return nil
}

func (tr *NewCache) DelCase(c *mcase.Case) error {
	if !tr.isNodeExist(mcase.Hash(tr.CaseKey(c))) {
		return fmt.Errorf("Case %s does not exist!", c.Name)
	}

	if _, ok := tr.Remove(mcase.Hash(tr.CaseKey(c))); !ok {
		return fmt.Errorf("Try to delete unexist case entry: %s\n", c.Name)
	}

	n, err := tr.Get(feature.Hash(tr.GetFeatureKeyByCase(c)))
	if err != nil {
		panic(err)
	}

	del := -1
	for i, ca := range n.Children {
		if res, ok := ca.Data.(*mcase.Case); ok {
			if res.ID == string(mcase.Hash(tr.CaseKey(c))) {
				del = i
				break
			}
		}
	}

	if del != -1 {
		if Data, ok := n.Data.(*feature.Feature); ok {
			Data.CCount--
			delete(Data.Cases, c.Name)
		}
		n.Children = append(n.Children[0:del], n.Children[del+1:]...)
	}

	tr.Save()
	return nil
}

func (tr *NewCache) DelCaseByID(id string) error {
	c, err := tr.GetCaseByID(id)
	if err != nil {
		return fmt.Errorf("Case %s does not exist!", id)
	}

	return tr.DelCase(c)
}

func (tr *NewCache) GetCase(c *mcase.Case) (*mcase.Case, error) {
	if !tr.isNodeExist(mcase.Hash(tr.CaseKey(c))) {
		return nil, fmt.Errorf("Case %s is not exist", string(tr.CaseKey(c)))
	}

	i, _ := tr.Lookup(mcase.Hash(tr.CaseKey(c)))
	if n, ok := i.(*Node); ok {
		if res, ok := n.Data.(*mcase.Case); ok {
			return res, nil
		} else {
			return nil, fmt.Errorf("Node: %s is not a Case!", string(tr.CaseKey(c)))
		}
	}
	return nil, fmt.Errorf("Invalid Node: %s", string(tr.CaseKey(c)))
}

func (tr *NewCache) AddGroup(g *group.Group) error {
	if tr.isNodeExist(group.Hash(tr.GetKeyByGroup(g))) {
		return fmt.Errorf("Group %s already exist", string(tr.GetKeyByGroup(g)))
	}

	g.ID = string(group.Hash(tr.GetKeyByGroup(g)))
	if g.SubGroups == nil {
		g.SubGroups = make(map[string]*subgroup.SubGroup, 1)
	}
	newnode := &Node{
		Key:      string(tr.GetKeyByGroup(g)),
		ID:       string(group.Hash(tr.GetKeyByGroup(g))),
		Type:     "GROUP",
		Link:     "groupinfo?id=" + string(group.Hash(tr.GetKeyByGroup(g))),
		Children: make([]*Node, 0, 1),
		Data:     g,
	}

	tr.Groups[g.Name] = g
	tr.Insert(group.Hash(tr.GetKeyByGroup(g)), newnode)
	tr.Node.Children = append(tr.Node.Children, newnode)

	return nil
}

func (tr *NewCache) DelGroup(g *group.Group) error {
	if !tr.isNodeExist(group.Hash(tr.GetKeyByGroup(g))) {
		return fmt.Errorf("Group %s is not exist", string(tr.GetKeyByGroup(g)))
	}

	if _, ok := tr.Remove(group.Hash(tr.GetKeyByGroup(g))); !ok {
		return fmt.Errorf("Try to delete unexist group: %s\n", g.Name)
	}

	del := -1
	for i, c := range tr.Node.Children {
		if res, ok := c.Data.(*group.Group); ok {
			if res.ID == string(group.Hash(tr.GetKeyByGroup(g))) {
				del = i
				break
			}
		}
	}

	if del != -1 {
		tr.Node.Children = append(tr.Node.Children[0:del], tr.Node.Children[del+1:]...)
		delete(tr.Groups, g.Name)
		tr.Save()
	}

	return nil
}

func (tr *NewCache) DelGroupByID(id string) error {
	g, err := tr.GetGroupByID(id)
	if err != nil {
		return fmt.Errorf("Group %s is not exist", id)
	}

	return tr.DelGroup(g)
}

func (tr *NewCache) GetGroup(g *group.Group) (*group.Group, error) {
	if !tr.isNodeExist(group.Hash(tr.GetKeyByGroup(g))) {
		return nil, fmt.Errorf("Group %s is not exist", string(tr.GetKeyByGroup(g)))
	}

	i, _ := tr.Lookup(group.Hash(tr.GetKeyByGroup(g)))
	if n, ok := i.(*Node); ok {
		if res, ok := n.Data.(*group.Group); ok {
			return res, nil
		} else {
			return nil, fmt.Errorf("Node: %s is not a Group!", string(tr.GetKeyByGroup(g)))
		}
	}
	return nil, fmt.Errorf("Invalid Node: %s", string(tr.GetKeyByGroup(g)))
}

func (tr *NewCache) AddSubGroup(sg *subgroup.SubGroup) error {
	if tr.isNodeExist(subgroup.Hash(tr.GetKeyBySubGroup(sg))) {
		return fmt.Errorf("SubGroup %s already exist", string(tr.GetKeyBySubGroup(sg)))
	}

	sg.ID = string(subgroup.Hash(tr.GetKeyBySubGroup(sg)))
	if sg.Features == nil {
		sg.Features = make(map[string]*feature.Feature, 1)
	}
	newnode := &Node{
		Key:      string(tr.GetKeyBySubGroup(sg)),
		ID:       string(subgroup.Hash(tr.GetKeyBySubGroup(sg))),
		Type:     "SUBGROUP",
		Link:     "subgroupinfo?id=" + string(subgroup.Hash(tr.GetKeyBySubGroup(sg))),
		Children: make([]*Node, 0, 1),
		Data:     sg,
	}

	tr.Insert(subgroup.Hash(tr.GetKeyBySubGroup(sg)), newnode)

	n, err := tr.Get(group.Hash(tr.GetGroupKeyBySubGroup(sg)))
	if err != nil {
		panic(err)
	}

	if Data, ok := n.Data.(*group.Group); ok {
		Data.SGCount++
		Data.SubGroups[sg.Name] = sg
	} else {
		log.Printf("!!!!!!: %s is not a subgroup\n", tr.GetKeyBySubGroup(sg))
	}
	n.Children = append(n.Children, newnode)

	return nil
}

func (tr *NewCache) DelSubGroup(sg *subgroup.SubGroup) error {
	if !tr.isNodeExist(subgroup.Hash(tr.GetKeyBySubGroup(sg))) {
		return fmt.Errorf("SubGroup %s is not exist", string(tr.GetKeyBySubGroup(sg)))
	}

	if _, ok := tr.Remove(subgroup.Hash(tr.GetKeyBySubGroup(sg))); !ok {
		return fmt.Errorf("Try to delete unexist subgroup: %s\n", sg.Name)
	}

	del := -1
	n, err := tr.Get(group.Hash(tr.GetGroupKeyBySubGroup(sg)))
	if err != nil {
		panic(err)
	}

	for i, k := range n.Children {
		if res, ok := k.Data.(*subgroup.SubGroup); ok {
			if res.ID == string(subgroup.Hash(tr.GetKeyBySubGroup(sg))) {
				del = i
				break
			}
		}
	}

	if del != -1 {
		if Data, ok := n.Data.(*group.Group); ok {
			Data.SGCount--
			delete(Data.SubGroups, sg.Name)
		}

		n.Children = append(n.Children[0:del], n.Children[del+1:]...)
	}

	tr.Save()
	return nil
}

func (tr *NewCache) DelSubGroupByID(id string) error {
	sg, err := tr.GetSubGroupByID(id)
	if err != nil {
		return fmt.Errorf("SubGroup %s is not exist", id)
	}

	return tr.DelSubGroup(sg)
}

func (tr *NewCache) GetSubGroup(sg *subgroup.SubGroup) (*subgroup.SubGroup, error) {
	if !tr.isNodeExist(subgroup.Hash(tr.GetKeyBySubGroup(sg))) {
		return nil, fmt.Errorf("SubGroup %s is not exist", string(tr.GetKeyBySubGroup(sg)))
	}

	i, _ := tr.Lookup(subgroup.Hash(tr.GetKeyBySubGroup(sg)))
	if n, ok := i.(*Node); ok {
		if res, ok := n.Data.(*subgroup.SubGroup); ok {
			return res, nil
		} else {
			return nil, fmt.Errorf("Node: %s is not a SubGroup!", string(tr.GetKeyBySubGroup(sg)))
		}
	}
	return nil, fmt.Errorf("Invalid Node: %s", string(tr.GetKeyBySubGroup(sg)))
}

func (tr *NewCache) AddFeature(f *feature.Feature) error {
	if tr.isNodeExist(feature.Hash(tr.GetKeyByFeature(f))) {
		return fmt.Errorf("Feature %s already exist", string(tr.GetKeyByFeature(f)))
	}

	f.ID = string(feature.Hash(tr.GetKeyByFeature(f)))
	if f.Cases == nil {
		f.Cases = make(map[string]*mcase.Case, 1)
	}
	newnode := &Node{
		Key:      string(tr.GetKeyByFeature(f)),
		ID:       string(feature.Hash(tr.GetKeyByFeature(f))),
		Type:     "FEATURE",
		Link:     "featureinfo?id=" + string(feature.Hash(tr.GetKeyByFeature(f))),
		Children: make([]*Node, 0, 1),
		Data:     f,
	}

	tr.Insert(feature.Hash(tr.GetKeyByFeature(f)), newnode)

	n, err := tr.Get(subgroup.Hash(tr.GetSubGroupKeyByFeature(f)))
	if err != nil {
		panic(err)
	}
	if Data, ok := n.Data.(*subgroup.SubGroup); ok {
		Data.FCount++
		Data.Features[f.Name] = f
	}
	n.Children = append(n.Children, newnode)

	return nil
}

func (tr *NewCache) DelFeature(f *feature.Feature) error {
	if !tr.isNodeExist(feature.Hash(tr.GetKeyByFeature(f))) {
		return fmt.Errorf("Feature %s is not exist", string(tr.GetKeyByFeature(f)))
	}

	if _, ok := tr.Remove(feature.Hash(tr.GetKeyByFeature(f))); !ok {
		return fmt.Errorf("Try to delete unexist feature entry: %s\n", f.Name)
	}

	del := -1
	n, err := tr.Get(subgroup.Hash(tr.GetSubGroupKeyByFeature(f)))
	if err != nil {
		panic(err)
	}

	for i, k := range n.Children {
		if res, ok := k.Data.(*feature.Feature); ok {
			if res.ID == string(feature.Hash(tr.GetKeyByFeature(f))) {
				del = i
				break
			}
		}
	}

	if del != -1 {
		if Data, ok := n.Data.(*subgroup.SubGroup); ok {
			Data.FCount--
			delete(Data.Features, f.Name)
		}
		n.Children = append(n.Children[0:del], n.Children[del+1:]...)
	}

	tr.Save()
	return nil
}

func (tr *NewCache) DelFeatureByID(id string) error {
	f, err := tr.GetFeatureByID(id)
	if err != nil {
		return fmt.Errorf("Feature %s is not exist", id)
	}

	return tr.DelFeature(f)
}

func (tr *NewCache) GetFeature(f *feature.Feature) (*feature.Feature, error) {
	if !tr.isNodeExist(feature.Hash(tr.GetKeyByFeature(f))) {
		return nil, fmt.Errorf("Feature %s is not exist", string(tr.GetKeyByFeature(f)))
	}

	i, _ := tr.Lookup(feature.Hash(tr.GetKeyByFeature(f)))
	if n, ok := i.(*Node); ok {
		if res, ok := n.Data.(*feature.Feature); ok {
			return res, nil
		} else {
			return nil, fmt.Errorf("Node: %s is not a Feature.", string(tr.GetKeyByFeature(f)))
		}
	}
	return nil, fmt.Errorf("Invalid Node: %s", string(tr.GetKeyByFeature(f)))
}

func (tr *NewCache) GetKeyByGroup(g *group.Group) []byte {
	return []byte(g.Name)
}

func (tr *NewCache) GetKeyBySubGroup(sg *subgroup.SubGroup) []byte {
	return []byte(sg.Group + sg.Name)
}

func (tr *NewCache) GetGroupKeyBySubGroup(sg *subgroup.SubGroup) []byte {
	return []byte(sg.Group)
}

func (tr *NewCache) GetKeyByFeature(f *feature.Feature) []byte {
	return []byte(f.Group + f.SubGroup + f.Name)
}

func (tr *NewCache) GetSubGroupKeyByFeature(f *feature.Feature) []byte {
	return []byte(f.Group + f.SubGroup)
}

func (tr *NewCache) GetGroupKeyByCase(c *mcase.Case) []byte {
	return []byte(c.Group)
}

func (tr *NewCache) GetSubGroupKeyByCase(c *mcase.Case) []byte {
	return []byte(c.Group + c.SubGroup)
}

func (tr *NewCache) GetFeatureKeyByCase(c *mcase.Case) []byte {
	return []byte(c.Group + c.SubGroup + c.Feature)
}

func (tr *NewCache) CaseKey(c *mcase.Case) []byte {
	return []byte(c.Group + c.SubGroup + c.Feature + c.Name)
}

func (tr *NewCache) Get(key []byte) (*Node, error) {
	if i, ok := tr.Lookup(key); ok {
		if n, ok := i.(*Node); ok {
			return n, nil
		}
		log.Println("Invalid result when Get node from DB")
	}

	return nil, fmt.Errorf("Cannot find node: %s", string(key))
}

func (tr *NewCache) IsValidCase(c *mcase.Case) bool {
	if c.Group == "" ||
		c.SubGroup == "" ||
		c.Feature == "" ||
		c.Name == "" {
		return false
	}
	return true
}

func (tr *NewCache) isNodeExist(key []byte) bool {
	if _, ok := tr.Lookup(key); ok {
		return true
	}
	return false
}

func (tr *NewCache) GetAllGroup() []*Node {
	groups := make([]*Node, 0, 1)
	for i := range tr.Iterator(nil) {
		if c, ok := i.Value.(*Node); ok {
			if c.Type == "GROUP" {
				groups = append(groups, c)
			}
		}
	}
	return groups
}

func (tr *NewCache) GetAllSubGroup() []*Node {
	sgroups := make([]*Node, 0, 1)
	for i := range tr.Iterator(nil) {
		if c, ok := i.Value.(*Node); ok {
			if c.Type == "SUBGROUP" {
				sgroups = append(sgroups, c)
			}
		}
	}
	return sgroups
}

func (tr *NewCache) GetAllFeature() []*Node {
	features := make([]*Node, 0, 1)
	for i := range tr.Iterator(nil) {
		if c, ok := i.Value.(*Node); ok {
			if c.Type == "FEATURE" {
				features = append(features, c)
			}
		}
	}
	return features
}

func (tr *NewCache) GetAllCase() []*mcase.Case {
	cases := make([]*mcase.Case, 0, 1)
	for i := range tr.Iterator(nil) {
		if c, ok := i.Value.(*Node); ok {
			if c.Type == "CASE" {
				if mc, ok := c.Data.(*mcase.Case); ok {
					cases = append(cases, mc)
				}
			}
		}
	}
	return cases
}

func (tr *NewCache) GetCaseByFeature(f *feature.Feature) ([]*mcase.Case, error) {
	nf, err := tr.GetFeature(f)
	if err != nil {
		return nil, err
	}

	res := make([]*mcase.Case, 0, 1)
	for _, c := range nf.Cases {
		res = append(res, c)
	}
	return res, nil
}

func (tr *NewCache) GetCaseBySubGroup(sg *subgroup.SubGroup) ([]*mcase.Case, error) {
	nsg, err := tr.GetSubGroup(sg)
	if err != nil {
		return nil, err
	}

	res := make([]*mcase.Case, 0, 1)

	for _, f := range nsg.Features {
		cs, err := tr.GetCaseByFeature(f)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		res = append(res, cs...)

	}
	return res, nil
}

func (tr *NewCache) GetCaseByGroup(g *group.Group) ([]*mcase.Case, error) {
	ng, err := tr.GetGroup(g)
	if err != nil {
		return nil, err
	}

	res := make([]*mcase.Case, 0, 1)

	for _, sg := range ng.SubGroups {
		cs, err := tr.GetCaseBySubGroup(sg)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		res = append(res, cs...)

	}
	return res, nil
}

func (tr *NewCache) RunCase(ctx context.Context, c *mcase.Case) <-chan *result.Result {
	res := make(chan *result.Result)
	go func(chan<- *result.Result) {
		message, ok := c.Run(ctx)
		defer close(res)
		res <- &result.Result{
			Group:    c.Group,
			SubGroup: c.SubGroup,
			Feature:  c.Feature,
			Case:     c.Name,
			Pass:     ok,
			Message:  message,
		}
	}(res)

	return res
}

func (tr *NewCache) RunCaseByID(ctx context.Context, id string) <-chan *result.Result {
	res := make(chan *result.Result)
	c, err := tr.GetCaseByID(id)
	if err != nil {
		go func(chan<- *result.Result) {
			defer close(res)
			res <- &result.Result{
				Group:    "Unknown",
				SubGroup: "Unknown",
				Feature:  "Unknown",
				Case:     "Unknown",
				Pass:     false,
				Message:  err.Error(),
			}
		}(res)
		return res
	}
	go func(chan<- *result.Result) {
		message, ok := c.Run(ctx)
		defer close(res)
		res <- &result.Result{
			Group:    c.Group,
			SubGroup: c.SubGroup,
			Feature:  c.Feature,
			Case:     c.Name,
			Pass:     ok,
			Message:  message,
		}
	}(res)

	return res
}

func (tr *NewCache) RunCasesByGroupID(ctx context.Context, id string) <-chan *result.Result {
	res := make(chan *result.Result)
	g, err := tr.GetGroupByID(id)
	if err != nil {
		go func(chan<- *result.Result) {
			defer close(res)
			res <- &result.Result{
				Group:    "Unknown",
				SubGroup: "Unknown",
				Feature:  "Unknown",
				Case:     "Unknown",
				Pass:     false,
				Message:  err.Error(),
			}
		}(res)
		return res
	}

	return tr.RunAllCaseOfGroup(ctx, g)
}

func (tr *NewCache) RunCasesBySubGroupID(ctx context.Context, id string) <-chan *result.Result {
	res := make(chan *result.Result)
	sg, err := tr.GetSubGroupByID(id)
	if err != nil {
		go func(chan<- *result.Result) {
			defer close(res)
			res <- &result.Result{
				Group:    "Unknown",
				SubGroup: "Unknown",
				Feature:  "Unknown",
				Case:     "Unknown",
				Pass:     false,
				Message:  err.Error(),
			}
		}(res)
		return res
	}

	return tr.RunAllCaseOfSubGroup(ctx, sg)
}

func (tr *NewCache) RunCasesByFeatureID(ctx context.Context, id string) <-chan *result.Result {
	res := make(chan *result.Result)
	f, err := tr.GetFeatureByID(id)
	if err != nil {
		go func(chan<- *result.Result) {
			defer close(res)
			res <- &result.Result{
				Group:    "Unknown",
				SubGroup: "Unknown",
				Feature:  "Unknown",
				Case:     "Unknown",
				Pass:     false,
				Message:  err.Error(),
			}
		}(res)
		return res
	}

	return tr.RunAllCaseOfFeature(ctx, f)
}

func (tr *NewCache) RunAllCaseOfFeature(ctx context.Context, f *feature.Feature) <-chan *result.Result {
	res := make(chan *result.Result)
	cases, err := tr.GetCaseByFeature(f)
	if err != nil {
		go func(err error) {
			res <- &result.Result{
				Group:    f.Group,
				SubGroup: f.SubGroup,
				Feature:  f.Name,
				Pass:     false,
				Message:  fmt.Sprintf("Cannot Find Feature: %s\n", f.Name),
			}
			close(res)
		}(err)
		return res
	}

	go func(cs []*mcase.Case) {
		wg := sync.WaitGroup{}
		for _, c := range cases {
			for r := range tr.RunCase(ctx, c) {
				wg.Add(1)
				go func(r *result.Result, rch chan<- *result.Result) {
					rch <- r
					wg.Done()
				}(r, res)
			}
		}
		wg.Wait()
		close(res)
	}(cases)

	return res
}

func (tr *NewCache) RunAllCaseOfSubGroup(ctx context.Context, sg *subgroup.SubGroup) <-chan *result.Result {
	res := make(chan *result.Result)
	cases := make([]*mcase.Case, 0, 1)
	for _, f := range sg.Features {
		cs, err := tr.GetCaseByFeature(f)
		if err != nil {
			go func(err error) {
				res <- &result.Result{
					Group:    f.Group,
					SubGroup: f.SubGroup,
					Feature:  f.Name,
					Pass:     false,
					Message:  fmt.Sprintf("Cannot Find Feature: %s\n", f.Name),
				}
			}(err)
		}
		cases = append(cases, cs...)
	}

	go func(cs []*mcase.Case) {
		wg := sync.WaitGroup{}
		for _, c := range cases {
			for r := range tr.RunCase(ctx, c) {
				wg.Add(1)
				go func(r *result.Result, rch chan<- *result.Result) {
					rch <- r
					wg.Done()
				}(r, res)
			}
		}
		wg.Wait()
		close(res)
	}(cases)

	return res
}

func (tr *NewCache) RunAllCaseOfGroup(ctx context.Context, g *group.Group) <-chan *result.Result {
	res := make(chan *result.Result)
	cases := make([]*mcase.Case, 0, 1)
	for _, sg := range g.SubGroups {
		for _, f := range sg.Features {
			cs, err := tr.GetCaseByFeature(f)
			if err != nil {
				go func(err error) {
					res <- &result.Result{
						Group:    f.Group,
						SubGroup: f.SubGroup,
						Feature:  f.Name,
						Pass:     false,
						Message:  fmt.Sprintf("Cannot Find Feature: %s\n", f.Name),
					}
				}(err)
			}
			cases = append(cases, cs...)
		}
	}

	go func(cs []*mcase.Case) {
		wg := sync.WaitGroup{}
		for _, c := range cases {
			for r := range tr.RunCase(ctx, c) {
				wg.Add(1)
				go func(r *result.Result) {
					log.Printf("%#v\n", r)
					res <- r
					wg.Done()
				}(r)
			}
		}
		wg.Wait()
		defer close(res)
	}(cases)

	return res
}

func (tr *NewCache) AddTask(caseid string, t *task.Task) error {
	c, err := tr.GetCaseByID(caseid)
	if err != nil {
		return err
	}

	t.ID = string(task.Hash([]byte(caseid + t.Name)))
	err = c.AddTask(t)
	if err != nil {
		return err
	}

	tr.Save()
	return nil
}

func (tr *NewCache) DelTask(caseid string, t *task.Task) error {
	c, err := tr.GetCaseByID(caseid)
	if err != nil {
		return err
	}

	defer tr.Save()
	err = c.DelTask(t)
	if err != nil {
		return err
	}

	return nil
}

func (tr *NewCache) DelTaskByID(caseid, taskid string) error {
	t, err := tr.GetTaskByID(caseid, taskid)
	if err != nil {
		return fmt.Errorf("Task %s of case: %s is not exit!", taskid, caseid)
	}

	return tr.DelTask(caseid, t)
}

func (tr *NewCache) GetTaskByID(caseid, taskid string) (*task.Task, error) {
	c, err := tr.GetCaseByID(caseid)
	if err != nil {
		return nil, fmt.Errorf("Case: %d does not exist!", caseid)
	}

	return c.GetTaskByID(taskid)
}

func (tr *NewCache) GetCaseByID(id string) (*mcase.Case, error) {
	if !tr.isNodeExist([]byte(id)) {
		return nil, fmt.Errorf("Case %s is not exist", id)
	}

	i, _ := tr.Lookup([]byte(id))
	if n, ok := i.(*Node); ok {
		if res, ok := n.Data.(*mcase.Case); ok {
			return res, nil
		} else {
			return nil, fmt.Errorf("Node: %s is not a Case!", id)
		}
	}
	return nil, fmt.Errorf("Invalid Node: %s", id)
}

func (tr *NewCache) GetGroupByID(id string) (*group.Group, error) {
	if !tr.isNodeExist([]byte(id)) {
		return nil, fmt.Errorf("Group %s is not exist", id)
	}

	i, _ := tr.Lookup([]byte(id))
	if n, ok := i.(*Node); ok {
		if res, ok := n.Data.(*group.Group); ok {
			return res, nil
		} else {
			return nil, fmt.Errorf("Node: %s is not a Group!", id)
		}
	}
	return nil, fmt.Errorf("Invalid Node: %s", id)
}

func (tr *NewCache) GetSubGroupByID(id string) (*subgroup.SubGroup, error) {
	if !tr.isNodeExist([]byte(id)) {
		return nil, fmt.Errorf("SubGroup %s is not exist", id)
	}

	i, _ := tr.Lookup([]byte(id))
	if n, ok := i.(*Node); ok {
		if res, ok := n.Data.(*subgroup.SubGroup); ok {
			return res, nil
		} else {
			return nil, fmt.Errorf("Node: %s is not a SubGroup!", id)
		}
	}
	return nil, fmt.Errorf("Invalid Node: %s", id)
}

func (tr *NewCache) GetFeatureByID(id string) (*feature.Feature, error) {
	if !tr.isNodeExist([]byte(id)) {
		return nil, fmt.Errorf("Feature %s is not exist", id)
	}

	i, _ := tr.Lookup([]byte(id))
	if n, ok := i.(*Node); ok {
		if res, ok := n.Data.(*feature.Feature); ok {
			return res, nil
		} else {
			return nil, fmt.Errorf("Node: %s is not a Feature!", id)
		}
	}
	return nil, fmt.Errorf("Invalid Node: %s", id)
}

func (tr *NewCache) RunTaskByID(ctx context.Context, caseid, taskid string) <-chan *result.Result {
	res := make(chan *result.Result)
	c, err := tr.GetCaseByID(caseid)
	if err != nil {
		go func(chan<- *result.Result) {
			defer close(res)
			res <- &result.Result{
				Pass:    false,
				Message: err.Error(),
			}
		}(res)
		return res
	}

	t, err := tr.GetTaskByID(caseid, taskid)
	if err != nil {
		go func(chan<- *result.Result) {
			defer close(res)
			res <- &result.Result{
				Group:    c.Group,
				SubGroup: c.SubGroup,
				Feature:  c.Feature,
				Case:     c.Name,
				Pass:     false,
				Message:  err.Error(),
			}
		}(res)
		return res
	}

	go func(chan<- *result.Result) {
		message, ok := c.RunTask(ctx, t)
		defer close(res)
		res <- &result.Result{
			Group:    c.Group,
			SubGroup: c.SubGroup,
			Feature:  c.Feature,
			Case:     c.Name,
			Task:     t.Name,
			Pass:     ok,
			Message:  message,
		}
	}(res)

	return res
}

func (tr *NewCache) GetDUTCountByID(id string) (int, error) {
	n, err := tr.Get([]byte(id))
	if err != nil {
		return -1, fmt.Errorf("Cannot find node by ID: %s", id)
	}

	switch n.Type {
	case "GROUP":
		return -1, fmt.Errorf("Currently not support get Group DUTCount")
	case "SUBGROUP":
		return -1, fmt.Errorf("Currently not support get SubGroup DUTCount")
	case "FEATURE":
		if f, ok := n.Data.(*feature.Feature); ok {
			return f.DUTCount, nil
		}
		return -1, fmt.Errorf("Internal Error: Get an Feature node but data is not feature")
	case "CASE":
		if c, ok := n.Data.(*mcase.Case); ok {
			return c.DUTCount, nil
		}
		return -1, fmt.Errorf("Internal Error: Get an Case node but data is not Case")
	default:
		return -1, fmt.Errorf("Unkown Node type for id: %s", id)
	}
}

func (tr *NewCache) SetDUTsByID(id string, duts []*rut.RUT) error {
	n, err := tr.Get([]byte(id))
	if err != nil {
		return fmt.Errorf("Cannot find node by ID: %s", id)
	}

	switch n.Type {
	case "GROUP":
		return fmt.Errorf("Currently not support set Group DUT")
	case "SUBGROUP":
		return fmt.Errorf("Currently not support set SubGroup DUT")
	case "FEATURE":
		if f, ok := n.Data.(*feature.Feature); ok {
			if f.DUTCount != len(duts) {
				return fmt.Errorf("Not Enough duts for feature:%s, required: %d, current: %d", f.Name, f.DUTCount, len(duts))
			}

			for _, d := range duts {
				if err != f.AddDUT(d) {
					f.ClearDUTs()
					return err
				}
			}
			log.Printf("%q", f.RUTDB)

			return nil
		}
		return fmt.Errorf("Internal Error: Get an Feature node but data is not feature")
	case "CASE":
		if c, ok := n.Data.(*mcase.Case); ok {
			if c.DUTCount != len(duts) {
				return fmt.Errorf("Not Enough duts for case:%s, required: %d, current: %d", c.Name, c.DUTCount, len(duts))
			}

			for _, d := range duts {
				c.AddRUT(d)
			}

			return nil
		}
		return fmt.Errorf("Internal Error: Get an Case node but data is not Case")
	default:
		return fmt.Errorf("Unkown Node type for id: %s", id)
	}
}

func (tr *NewCache) CheckIsReadyForRunByID(id string) (bool, error) {
	n, err := tr.Get([]byte(id))
	if err != nil {
		return false, fmt.Errorf("Cannot find node by ID: %s", id)
	}

	switch n.Type {
	case "GROUP":
		return false, fmt.Errorf("Currently not support run by Group")
	case "SUBGROUP":
		return false, fmt.Errorf("Currently not support run by SubGroup")
	case "FEATURE":
		if f, ok := n.Data.(*feature.Feature); ok {
			if f.DUTCount != len(f.RUTDB) {
				f.ClearDUTs()
				return false, fmt.Errorf("Not Enough duts for feature:%s, required: %d, current: %d", f.Name, f.DUTCount, len(f.RUTDB))
			}

			if f.DUTCount == 0 {
				return false, fmt.Errorf("Invalid feature test case, no DUT exist!")
			}

			log.Println(f.DUTCount, len(f.RUTDB))

			return true, nil
		}
		return false, fmt.Errorf("Internal Error: Get an Feature node but data is not feature")
	case "CASE":
		if c, ok := n.Data.(*mcase.Case); ok {
			if c.DUTCount != len(c.RUTs.DB) {
				c.ClearDUTs()
				return false, nil
			}

			if c.DUTCount == 0 {
				return false, fmt.Errorf("Invalid test case, no DUT exist!")
			}

			log.Println(c.DUTCount, len(c.RUTs.DB))
			log.Printf("%#v: %#v", c.RUTs, c.RUTs.DB)

			return true, nil
		}
		return false, fmt.Errorf("Internal Error: Get an Case node but data is not Case")
	default:
		return false, fmt.Errorf("Unkown Node type for id: %s", id)
	}
}
