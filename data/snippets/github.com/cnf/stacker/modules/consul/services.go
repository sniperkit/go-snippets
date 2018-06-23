package consul

import (
	"github.com/golang/glog"

	//"github.com/armon/consul-api"
	"github.com/cnf/stacker/engine"

	// "github.com/kr/pretty"
)

// Services object
type Services struct {
	Config *Config
}

// NewReaction returns a new consul Reaction object
func NewReaction() engine.Reaction {
	return &Services{}
}

// RegisterReaction returns the module name
func RegisterReaction() {
	engine.RegisterReaction("consul", NewReaction)
}

// Setup ...
func (s *Services) Setup(eng *engine.Engine, cfg engine.Config) error {
	s.Config = cfg.(*Config)
	return nil
}

// Start the consul service advertiser
func (s *Services) Start() error {
	if glog.V(1) {
		glog.Infof("starting consul service advertiser")
	}
	return nil
}

// Stop the consul service advertiser
func (s *Services) Stop() error {
	if glog.V(1) {
		glog.Infof("stopping consul service advertiser")
	}
	return nil
}

// Command takes a Command object, and executes actions
func (s *Services) Command(cmd *engine.Command) error {
	// TODO: all of it
	if glog.V(2) {
		glog.Infof("consul got trigger for '%s'", cmd.Name)
		// pretty.Print("\n", cmd.Payload, "\n")
	}
	switch cmd.Name {
	case "start":
	case "die":
	}
	return nil
}
