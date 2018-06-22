package registry

import (
	"github.com/golang/glog"

	"github.com/cnf/stacker/engine"
)

var list *List

// List of Registries we know about
type List struct {
	Registries []*Registry
}

// NewList returns a new Registry List object
func NewList() *List {
	if list == nil {
		list = &List{}
	}
	return list
}

// BuildNewList (re)creates the internal registry list
func (rl *List) BuildNewList() error {
	// TODO: ...
	//if err := cl.initializeList(); err != nil {
	//	return err
	//}
	reg := NewRegistry("index.docker.io")
	rl.Registries = append(rl.Registries, reg)
	return nil
}

// Command takes a Command object and triggers events
func (rl *List) Command(cmd *engine.Command) error {
	if cmd.Name != "config" {
		// TODO: proper error
		return nil
	}
	return rl.update(cmd.Payload.(*Config))
}

// update registry from configuration
func (rl *List) update(cfg *Config) error {
	if glog.V(2) {
		glog.Infof("UPDATE Registry: %#v", cfg.Name)
	}
	reg := rl.getByName(cfg.Name)
	if reg != nil {
		reg.Name = cfg.Name
		reg.Username = cfg.Username
		reg.Password = cfg.Password
		reg.Email = cfg.Email
		reg.ServerAddress = cfg.ServerAddress
		return nil
	}
	newreg := NewRegistry(cfg.Name)
	newreg.Username = cfg.Username
	newreg.Password = cfg.Password
	newreg.Email = cfg.Email
	newreg.ServerAddress = cfg.ServerAddress
	rl.Registries = append(rl.Registries, newreg)
	return nil
}

func (rl *List) getByName(name string) *Registry {
	for _, reg := range rl.Registries {
		if reg.Name == name {
			return reg
		}
	}
	return nil
}
