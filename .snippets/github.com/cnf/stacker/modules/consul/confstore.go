package consul

import (
	"github.com/golang/glog"
	//"github.com/armon/consul-api"

	"github.com/cnf/stacker/engine"
)

// ConfStore object
type ConfStore struct {
	Config *Config
}

// NewConfStore returns a new consul ConfStore object
func NewConfStore() engine.ConfStore {
	return &ConfStore{}
}

// RegisterConfStore returns the module name
func RegisterConfStore() {
	engine.RegisterConfStore("consul", NewConfStore)
}

// Setup ...
func (cs *ConfStore) Setup(eng *engine.Engine, cfg engine.Config) error {
	cs.Config = cfg.(*Config)
	return nil
}

// Start the ConfStore
func (cs *ConfStore) Start() error {
	if glog.V(1) {
		glog.Infof("starting consul configstore")
	}
	return nil
}

// Stop the ConfStore
func (cs *ConfStore) Stop() error {
	return nil
}
