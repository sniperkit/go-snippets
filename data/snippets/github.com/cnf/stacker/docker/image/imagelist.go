package image

import (
	"github.com/cnf/stacker/docker/wrapper"
	"github.com/cnf/stacker/engine"
)

// List of docker Images we know of
type List struct {
	Images []*Image
	eng    *engine.Engine
}

// NewList returns a new Image List object
func NewList(eng *engine.Engine) *List {
	return &List{eng: eng}
}

// BuildNewList (re)creates the internal image list
func (il *List) BuildNewList() error {
	if err := il.initializeList(); err != nil {
		return err
	}
	return nil
}

// Command takes a Command object and triggers events
func (il *List) Command(cmd *engine.Command) error {
	img := il.getByName(cmd.Destination.ID)
	if img == nil {
		img = NewImage(cmd.Destination.ID)
		il.Images = append(il.Images, img)
	}
	if err := img.Command(cmd); err != nil {
		return err
	}
	newcmd := &engine.Command{
		Name:   "done",
		Source: img.Target(),
	}
	switch cmd.Name {
	case "config":
		newcmd.Destination = nil
	case "need":
		newcmd.Destination = cmd.Source
	default:
		return nil
	}
	il.eng.Event(newcmd)
	return nil
}

func (il *List) initializeList() error {
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}

	newil := []*Image{}

	// Get a list of images from docker, and iterate
	dil, err := dw.Client.ListImages(false)
	if err != nil {
		return err
	}
	for _, v := range dil {
		for _, rt := range v.RepoTags {
			if entryInList(rt, newil) {
				continue
			}
			img := NewImage(rt)
			img.ID = v.ID
			newil = append(newil, img)
		}
	}

	// pretty.Print(newil)
	il.Images = newil
	return nil
}

func (il *List) getByID(ID string) *Image {
	for _, img := range il.Images {
		if img.ID == ID {
			return img
		}
	}
	return nil
}

func (il *List) getByName(name string) *Image {
	for _, img := range il.Images {
		if img.Name == name {
			return img
		}
	}
	return nil
}

func entryInList(entry string, list []*Image) bool {
	for _, img := range list {
		if entry == img.Name {
			return true
		}
	}
	return false
}
