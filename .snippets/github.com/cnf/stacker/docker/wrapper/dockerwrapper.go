package dockerwrapper

import (
	"github.com/cnf/structhash"
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"

	"github.com/cnf/stacker/engine"
)

var dw *Docker

// Channel is the docker API event channel
var Channel chan *dockerapi.APIEvents

// Docker defines a docker instance used to talk to docker
type Docker struct {
	Client *dockerapi.Client
	// Channel   chan *docker.APIEvents
}

// Config for our docker connection
type Config struct {
	Socket string `mapstructure:"socket" version:"1"`
	Cert   string `mapstructure:"cert" version:"1"`
	Key    string `mapstructure:"key" version:"1"`
	CA     string `mapstructure:"ca" version:"1"`
}

// Target returns the namespace of the target of this config
func (c *Config) Target() *engine.NameSpace {
	return &engine.NameSpace{
		Type:   "crt",
		Module: "docker",
	}
}

// NewConfig returns a Docker Config with defaults
func NewConfig() *Config {
	return &Config{
		Socket: "unix:///var/run/docker.sock",
	}
}

// Version of our config struct.
// TODO: automate
func (c *Config) Version() int {
	return 1
}

// Hash returns the unique hash for this config object
func (c *Config) Hash() string {
	hash, err := structhash.Hash(c, engine.LatestHashVersion)
	if err != nil {
		return ""
	}
	return hash
}

// GetObject returns the docker wrapper object
func GetObject() (*Docker, error) {
	if dw == nil {
		dw = &Docker{}
		Channel = make(chan *dockerapi.APIEvents)
	}
	return dw, nil
}

// Connect to the socket.
func (d *Docker) Connect(cfg *Config) error {
	// TODO: TLS
	client, err := dockerapi.NewClient(cfg.Socket)
	if err == dockerapi.ErrInvalidEndpoint && err == dockerapi.ErrConnectionRefused {
		glog.Errorf("%s: %s", err.Error(), cfg.Socket)
		return err
	}
	if err != nil {
		return err
	}
	dw.Client = client
	return nil
}

// StartListener starts listening to the docker api.
func (d *Docker) StartListener() error {
	if err := d.Client.AddEventListener(Channel); err != nil {
		return err
	}
	return nil
}

// StopListener stops listening to the docker api.
func (d *Docker) StopListener() error {
	if d.Client == nil {
		return nil
	}
	if err := d.Client.RemoveEventListener(Channel); err != nil {
		return err
	}
	return nil
}
