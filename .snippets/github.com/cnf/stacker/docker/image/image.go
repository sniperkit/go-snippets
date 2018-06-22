package image

import (
	"sync"

	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"

	"github.com/cnf/stacker/docker/registry"
	"github.com/cnf/stacker/docker/wrapper"
	"github.com/cnf/stacker/engine"
)

// Image represents a docker image
type Image struct {
	sync.Mutex
	Name   string
	ID     string
	Config *Config
}

// NewImage returns a new image object
func NewImage(name string) *Image {
	return &Image{
		Name: name,
	}
}

// Target returns the namespace of the target of this config
func (i *Image) Target() *engine.NameSpace {
	return &engine.NameSpace{
		Type:   "crt",
		Module: "image",
		ID:     i.Name,
	}
}

// Command takes a Command object and triggers events
func (i *Image) Command(cmd *engine.Command) error {
	i.Lock()
	defer i.Unlock()
	switch cmd.Name {
	case "config":
		i.Config = cmd.Payload.(*Config)
		i.update()
	case "need":
		i.update()
	}
	return nil
}

func (i *Image) update() error {
	if i.Config != nil && i.Config.State == "removed" {
		if err := i.ensureAbsent(); err != nil {
			return err
		}
		return nil
	}
	if err := i.ensurePresent(); err != nil {
		return err
	}
	return nil
}

func (i *Image) pull() error {
	if glog.V(1) {
		glog.Infof("pulling image '%s'", i.Name)
	}
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	auth := registry.GetAuthConfiguration(i.Name)
	opts := getPullImageOptions(i.Name)
	if err := dw.Client.PullImage(*opts, *auth); err != nil {
		glog.Error(err.Error())
		return err
	}
	return nil
}

func (i *Image) ensurePresent() error {
	if glog.V(2) {
		glog.Infof("making sure image '%s' is present", i.Name)
	}
	if err := i.pull(); err != nil {
		return err
	}
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	img, err := dw.Client.InspectImage(i.Name)
	if err != nil {
		return err
	}
	if i.ID != img.ID {
		glog.Infof("updating ID for image %s", i.Name)
	}
	return nil
}

func (i *Image) ensureAbsent() error {
	if glog.V(2) {
		glog.Infof("making sure image '%s' is not present", i.Name)
	}
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	if err := dw.Client.RemoveImage(i.Name); err != nil {
		return err
	}
	i.ID = ""
	return nil
}

func getPullImageOptions(name string) *dockerapi.PullImageOptions {
	repo, tag, err := parseRepoName(name)
	if err != nil {
		glog.Error(err.Error())
		return nil
	}
	return &dockerapi.PullImageOptions{
		Repository: repo,
		// Registry // TODO: own registry support
		Tag: tag,
	}
}

func getAuthConfiguration() *dockerapi.AuthConfiguration {
	return &dockerapi.AuthConfiguration{}
}
