package logger

import (
	"github.com/golang/glog"

	"github.com/cnf/stacker/engine"
)

// Logger object
type Logger struct {
	Config *Config
}

// NewReaction returns a new logger Reaction object
func NewReaction() engine.Reaction {
	return &Logger{}
}

// RegisterReaction returns the module name
func RegisterReaction() {
	engine.RegisterReaction("logger", NewReaction)
}

// Setup ...
func (l *Logger) Setup(eng *engine.Engine, cfg engine.Config) error {
	l.Config = cfg.(*Config)
	return nil
}

// Start the logger service
func (l *Logger) Start() error {
	// TODO: ...
	return nil
}

// Stop the logger service
func (l *Logger) Stop() error {
	// TODO: ...
	return nil
}

// Command takes a Command object, and logs it
func (l *Logger) Command(cmd *engine.Command) error {
	// TODO: all of it
	switch cmd.Name {
	case "log":
		glog.Infof(">>>>> %s for %s", cmd.Name, cmd.Source.String())
	}
	return nil
}
