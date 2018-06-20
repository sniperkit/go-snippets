package lib

import (
	"errors"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//Config for system
var Config = &Configuration{}

//Configuration keep the related configuration options
type Configuration struct {
	Host        string          `yaml:"host"`
	Port        uint            `yaml:"port"`
	Dockerd     *DockerdConfig  `yaml:"dockerd"`
	Harbor      *HarborConfig   `yaml:"harbor"`
	NpmRegistry *RegistryConfig `yaml:"npm_registry"`
	PipRegistry *RegistryConfig `yaml:"pip_registry"`
}

//DockerdConfig is for dockerd
type DockerdConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Admin    string `yaml:"admin"`
	Password string `yaml:"password"`
}

//HarborConfig is for harbor
type HarborConfig struct {
	Host     string `yaml:"host"`
	Protocol string `yaml:"protocol"`
}

//RegistryConfig is for npm registries
type RegistryConfig struct {
	Namespace    string `yaml:"namespace"`
	BaseImage    string `yaml:"base_image"`
	BaseImageTag string `yaml:"base_image_tag"`
}

//Load configurations from yaml file
func (c *Configuration) Load(yamlFile string) error {
	if len(yamlFile) == 0 {
		return errors.New("empty config file")
	}

	fileData, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	if len(fileData) == 0 {
		return errors.New("no any configuration options existing")
	}

	if err := yaml.Unmarshal(fileData, c); err != nil {
		return err
	}

	return c.validate()
}

func (c *Configuration) validate() error {
	if c.Port < 256 {
		return fmt.Errorf("port should be greater than 256, but got %d", c.Port)
	}

	if c.Dockerd == nil {
		return errors.New("dockerd is not configured")
	}

	if err := c.validateDockerd(); err != nil {
		return err
	}

	if c.Harbor == nil {
		return errors.New("Harbor is not configured")
	}

	if err := c.validateHarbor(); err != nil {
		return err
	}

	if c.NpmRegistry == nil {
		return errors.New("npm registry is not configured")
	}

	if err := c.validateNpmRegistry(); err != nil {
		return err
	}

	if c.PipRegistry == nil {
		return errors.New("pip registry is not configured")
	}

	return c.validatePipRegistry()
}

func (c *Configuration) validateDockerd() error {
	if len(c.Dockerd.Host) == 0 {
		return errors.New("dockerd host is not configured")
	}

	if c.Dockerd.Port == 0 {
		return errors.New("dockerd port is not configured")
	}

	return nil
}

func (c *Configuration) validateHarbor() error {
	if len(c.Harbor.Host) == 0 {
		return errors.New("harbor host is not configured")
	}

	if c.Harbor.Protocol != "http" && c.Harbor.Protocol != "https" {
		return fmt.Errorf("harbor protocol is only supporting 'http' or 'https'")
	}

	return nil
}

func (c *Configuration) validateNpmRegistry() error {
	if len(c.NpmRegistry.BaseImage) == 0 {
		return errors.New("npm base image is nil")
	}

	if len(c.NpmRegistry.BaseImageTag) == 0 {
		return errors.New("npm base image tag is nil")
	}

	if len(c.NpmRegistry.Namespace) == 0 {
		return errors.New("no namespace is specified for npm registry")
	}

	return nil
}

func (c *Configuration) validatePipRegistry() error {
	if len(c.PipRegistry.Namespace) == 0 {
		return errors.New("no namespace is specified for pip registry")
	}

	return nil
}
