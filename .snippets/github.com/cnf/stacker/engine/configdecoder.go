package engine

import (
	"errors"
	"strings"

	"github.com/golang/glog"
)

// ErrNotAConfig ...
var ErrNotAConfig = errors.New("not a config object")

// Config is the interface for every module Config object.
type Config interface {
	Target() *NameSpace
	Version() int
	Hash() string
}

// ConfigDecoder is an interface to parsing and validating configuration data
type ConfigDecoder interface {
	DecodeConfig(cont interface{}) (Config, error)
	// TODO:
	// ValidateConfig(*interface{}) (*interface{}, error)
}

var configDecoderList = make(map[string]ConfigDecoder)

// RegisterConfigDecoder registers a ConfigDecoder with the engine.
func RegisterConfigDecoder(name string, cd ConfigDecoder) {
	name = strings.ToLower(name)
	if glog.V(1) {
		glog.Infof("registering ConfigDecoder: %s", name)
	}

	if configDecoderList[name] != nil {
		glog.Fatalf("a ConfigDecoder named `%s` already exists", name)
	}
	configDecoderList[name] = cd
}
