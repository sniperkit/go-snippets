package mcase

import (
	"errors"
	"log"
	"net/url"
	//"strconv"
	"context"
	"crypto/sha1"
	"defaults"
	"encoding/hex"
	"fmt"
	"logger"
	"rut"
	"strings"
	"task"
	"util"
)

type Case struct {
	Group    string
	SubGroup string
	Feature  string
	Name     string
	Tasks    []*task.Task
	RUTs     rut.DB
	ID       string
	DUTCount int
	TCount   int
	Topology string
	Type     string
}

func Hash(name []byte) []byte {
	hash := sha1.New()
	return []byte(hex.EncodeToString(hash.Sum([]byte(name))))
}

func (c *Case) Hash(name []byte) []byte {
	return Hash(name)
}

func (c *Case) GenerateTaskID(t *task.Task) string {
	return string(c.Hash([]byte(c.ID + t.Name)))
}

func (c *Case) String() string {
	var duts string

	for k, v := range c.RUTs.DB {
		duts += fmt.Sprintf("[%s->%s]", k, v.Device)
	}
	return fmt.Sprintf("%s>%s>%s>%s {%s}", c.Group, c.SubGroup, c.Feature, c.Name, duts)
}

func (c *Case) Init() {
	for _, r := range c.RUTs.DB {
		r.Init()
	}
}

func (c *Case) MakeTreeViewKey() string {
	return fmt.Sprintf("%s%s:::%s%s:::%s%s:::%s%s", c.Group, defaults.GroupMark, c.SubGroup, defaults.SubGroupMark, c.Feature, defaults.FeatureMark, c.Name, defaults.CaseMark) //Give case name an identifier
}

func (c *Case) AddRUT(r *rut.RUT) {
	log.Printf("Add DUT %#v\n", r)
	if len(c.RUTs.DB) == 0 {
		c.RUTs.DB = make(map[string]*rut.RUT, 1)
	}
	c.RUTs.DB[r.Name] = r

	log.Printf("%q", c.RUTs.DB)
}

func MakeCaseFromTreeViewKey(key string) (*Case, error) {
	slices := strings.Split(key, ":::")
	if len(slices) == 0 {
		return nil, errors.New("Invalid tree view Key")
	}

	var c Case
	for _, s := range slices {
		if strings.HasSuffix(s, defaults.GroupMark) {
			c.Group = s[:len(s)-len(defaults.GroupMark)]
			continue
		}
		if strings.HasSuffix(s, defaults.SubGroupMark) {
			c.SubGroup = s[:len(s)-len(defaults.SubGroupMark)]
			continue
		}
		if strings.HasSuffix(s, defaults.FeatureMark) {
			c.Feature = s[:len(s)-len(defaults.FeatureMark)]
			continue
		}
		if strings.HasSuffix(s, defaults.CaseMark) {
			c.Feature = s[:len(s)-len(defaults.CaseMark)]
			continue
		}
	}

	return &c, nil
}

func (c *Case) DelRUT(r *rut.RUT) {
	if r, ok := c.RUTs.DB[r.Name]; ok {
		delete(c.RUTs.DB, r.Name)
	}
}

func (c *Case) ClearDUTs() {
	if c.RUTs.DB == nil {
		log.Printf("Clear RUTs called for case: %s when there are ruts in db", c.Name)
	}

	for name, _ := range c.RUTs.DB {
		delete(c.RUTs.DB, name)
	}
}

func (c *Case) Run(ctx context.Context) (string, bool) {
	if len(c.Tasks) == 0 {
		log.Printf("There is no task in case: %s:%s:%s:%s\n", c.Group, c.SubGroup, c.Feature, c.Name)
		message := fmt.Sprintf("There is no task in case: %s:%s:%s:%s\n", c.Group, c.SubGroup, c.Feature, c.Name)
		logger.Push(ctx, message)
		return message, true
	}

	if c.RUTs.DB == nil {
		message := fmt.Sprintf("You should first specify the RUT for case : %s:%s:%s:%s\n", c.Group, c.SubGroup, c.Feature, c.Name)
		logger.Push(ctx, message)
		return message, true
	}

	for _, t := range c.Tasks {
		result := t.Run(ctx, &c.RUTs)
		if !result.Success {
			message := fmt.Sprintf("Run Case: %s>%s>%s>%s's Task: %s failed with: %s", c.Group, c.SubGroup, c.Feature, c.Name, t.Name, result.Message)
			logger.Push(ctx, message)
			return message, false
		}
	}

	return "", true
}

func (c *Case) RunTask(ctx context.Context, t *task.Task) (string, bool) {
	result := t.Run(ctx, &c.RUTs)
	return result.Message, result.Success
}

func IsValidCaseParas(in url.Values) bool {
	if v, ok := in["case_group"]; !ok {
		log.Println(v)
		return false
	} else if v[0] == "" {
		return false
	}

	if v, ok := in["case_sub_group"]; !ok {
		return false
	} else if v[0] == "" {
		return false
	}

	if v, ok := in["case_feature"]; !ok {
		return false
	} else if v[0] == "" {
		return false
	}

	if v, ok := in["case_name"]; !ok {
		return false
	} else if v[0] == "" {
		return false
	}

	for k, _ := range in {
		if strings.HasPrefix(k, "device") {
			//if _, err := strconv.ParseInt(k[3:], 10, 64); err != nil {
			return true
			//	}
		}
	}

	return false
}

func (c *Case) AddTask(t *task.Task) error {
	if c.Tasks == nil {
		c.Tasks = make([]*task.Task, 0, 1)
	}

	t.PreCondition.Name = "Pre-Condition"
	t.PostCondition.Name = "Post-Condition"

	t.ID = c.GenerateTaskID(t) //Necessary
	if c.IsTaskExist(t) {
		log.Printf("Updadte task: %s, ID: %s\n", t.Name, t.ID)
		var index int
		for i, v := range c.Tasks {
			if v.Name == t.Name {
				index = i
				break
			}
		}
		if index < len(c.Tasks) {
			c.Tasks[index] = t
		}
	} else {
		c.Tasks = append(c.Tasks, t)
		c.TCount++
	}

	return nil
}

func (c *Case) DelTask(t *task.Task) error {
	if !c.IsTaskExist(t) {
		return errors.New("Cannot find task :" + t.Name + " in case: " + c.Name)
	}

	for i, v := range c.Tasks {
		if v.Name == t.Name {
			c.Tasks = append(c.Tasks[:i], c.Tasks[i+1:]...)
		}
	}

	c.TCount--
	return nil
}

func (c *Case) GetTask(name string) *task.Task {
	for _, v := range c.Tasks {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func (c *Case) GetTaskByID(id string) (*task.Task, error) {
	for _, v := range c.Tasks {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, fmt.Errorf("Task: %s does not exist!", id)
}

func (c *Case) IsTaskExist(t *task.Task) bool {
	for _, v := range c.Tasks {
		if v.Name == t.Name {
			return true
		}
	}

	return false
}

func (c *Case) SetTopology(name string, content []byte) {
	c.Topology = "asset/web/topology/" + name
	util.SaveToFile(c.Topology, content)
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
