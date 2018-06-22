package container

import (
	"strings"

	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"

	"github.com/cnf/stacker/docker/wrapper"
	"github.com/cnf/stacker/engine"
)

// List of Containers we know about
type List struct {
	Containers []*Container
}

// NewList returns a new container List object
func NewList() *List {
	return &List{
		Containers: make([]*Container, 0),
	}
}

// BuildNewList (re)creates the internal container list of config and running
// containers
func (cl *List) BuildNewList(eng *engine.Engine) error {
	if err := cl.initializeList(eng); err != nil {
		return err
	}
	return nil
}

// Add adds a container to the container List
func (cl *List) Add(name string, eng *engine.Engine) (*Container, error) {
	cont := NewContainer(name, eng)
	if len(cl.Containers) == 0 {
		glog.Error("WTF?????")
		return nil, nil
	}
	cl.Containers = append(cl.Containers, cont)
	return cont, nil
}

// React to something
// TODO: RENAME
func (cl *List) React(id string, event string) error {
	glog.Infof("container:react: %s for %s", event, id[:12])
	c := cl.getByID(id)
	if c != nil {
		c.Channel <- event
		// c.event(event)
	}
	return nil
}

// Length of the container list
func (cl *List) Length() int {
	return len(cl.Containers)
}

func (cl *List) GetNameFromID(id string) string {
	c := cl.getByID(id)
	if c != nil {
		return c.Name
	}
	return ""
}

func (cl *List) getByID(id string) *Container {
	for _, c := range cl.Containers {
		if c.ID == id {
			return c
		}
	}
	return nil
}

func (cl *List) getByName(name string) *Container {
	for _, c := range cl.Containers {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (cl *List) closeAllChannels() {
	for _, c := range cl.Containers {
		close(c.Channel)
	}
}

func (cl *List) initializeList(eng *engine.Engine) error {
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	ncl := []*Container{}

	lco := dockerapi.ListContainersOptions{All: true}
	// Get a list of containers from docker, and iterate
	dcl, err := dw.Client.ListContainers(lco)
	if err != nil {
		glog.Infof("container:list: %s", err.Error())
		return err
	}
	// TODO: use Add()
	for _, v := range dcl {
		ci, _ := dw.Client.InspectContainer(v.ID)
		ciName := strings.TrimLeft(ci.Name, "/")
		cHash := ""
		for _, env := range ci.Config.Env {
			evar := strings.SplitN(env, "=", 2)
			if evar[0] == "STACKER_CONTAINER_CFG" {
				cHash = evar[1]
				break
			}
		}
		cont := NewContainer(ciName, eng)
		cont.ID = ci.ID
		cont.DockerState = &ci.State
		cont.Hash = cHash

		ncl = append(ncl, cont)
	}

	cl.Containers = ncl
	// TODO: make container manage its own channel
	//for _, c := range d.Containers {
	//	c.Channel = make(chan string)
	//	go c.Reactor()
	//	// c.Channel <- "update"
	//}
	return nil
}
