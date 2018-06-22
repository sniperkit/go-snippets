package docker

import (
	"github.com/golang/glog"

	"github.com/cnf/stacker/docker/container"
	"github.com/cnf/stacker/docker/image"
	"github.com/cnf/stacker/docker/registry"
	"github.com/cnf/stacker/docker/wrapper"
	"github.com/cnf/stacker/engine"
)

// Docker represents the docker process we are connected to
type Docker struct {
	Registries *registry.List
	Containers *container.List
	Images     *image.List

	stopper chan bool
	eng     *engine.Engine
	ns      *engine.NameSpace
}

// NewObject returns a new Docker object
func NewObject() engine.ContainerRuntime {
	return &Docker{
		ns: &engine.NameSpace{
			Type:   "crt",
			Module: "docker",
		},
	}
}

// Register docker with the engine
func Register() {
	engine.RegisterContainerRuntime("docker", NewObject)
}

// Setup the docker instance
func (d *Docker) Setup(eng *engine.Engine, cfg engine.Config) error {
	d.eng = eng
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	if conf, ok := cfg.(*dockerwrapper.Config); ok {
		if err := dw.Connect(conf); err != nil {
			glog.Errorf(err.Error())
			return err
		}
	}
	return nil
}

// Start docker instance
func (d *Docker) Start() error {
	d.stopper = make(chan bool)
	go d.listener()

	if err := d.buildNewLists(); err != nil {
		return err
	}
	for _, v := range d.Containers.Containers {
		if v.Name != "" {
			d.eng.RegisterContainer(v.Name, v)
		}
	}

	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	if err := dw.StartListener(); err != nil {
		return err
	}
	return nil
}

// Stop docker instance
func (d *Docker) Stop() error {
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	if err := dw.StopListener(); err != nil {
		return err
	}
	close(d.stopper)
	return nil
}

// BuildNewLists of containers, images and volumes
func (d *Docker) buildNewLists() error {
	if glog.V(1) {
		glog.Info("building new lists")
	}
	d.Registries = registry.NewList()
	if err := d.Registries.BuildNewList(); err != nil {
		return err
	}
	d.Images = image.NewList(d.eng)
	if err := d.Images.BuildNewList(); err != nil {
		return err
	}
	d.Containers = container.NewList()
	if err := d.Containers.BuildNewList(d.eng); err != nil {
		return err
	}
	return nil
}

// Command takes a Command object and triggers events
func (d *Docker) Command(cmd *engine.Command) error {
	if glog.V(2) {
		glog.Infof("reacting to event '%s'", cmd.Name)
	}
	if cmd.Source != nil && cmd.Source == d.ns {
		// TODO: don't respond to our own messages
		return nil
	}
	switch cmd.Name {
	case "config":
		if cmd.Destination == nil {
			// TODO: fix
			return nil
		}
		switch cmd.Destination.Module {
		case "container":
			name := cmd.Payload.(*container.Config).Name
			d.newContainer(name)
			cmd.Source = d.ns
			d.eng.Event(cmd)
			return nil
		case "image":
			return d.Images.Command(cmd)
		case "registry":
			return d.Registries.Command(cmd)
		}
	case "need":
		switch cmd.Destination.Module {
		case "container":
			glog.Warningf("container '%s' does not exist", cmd.Destination.ID)
			// TODO: proper error
			return nil
		case "image":
			// TODO:
			return d.Images.Command(cmd)

		}
	case "depcheck":
		d.checkDependencies(cmd)
	}
	return nil
}

func (d *Docker) newContainer(name string) *container.Container {
	c, err := d.Containers.Add(name, d.eng)
	if err != nil {
		// TODO: error handeling
		return nil
	}
	d.eng.RegisterContainer(name, c)
	return c
}
