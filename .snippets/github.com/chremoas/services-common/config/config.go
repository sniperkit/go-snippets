package config

import (
	"errors"
	"fmt"
	"github.com/micro/go-micro"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config interface {
	Load(filename string) error
	NewConnectionString() (string, error)
	NewService(version, defaultName string) (micro.Service, error)
	AuthServiceName() (string, error)
	LookupService(serviceType string, serviceName string) (serviceFullName string)
}

type Configuration struct {
	initialized bool
	Namespace   string
	Database    struct {
		Driver         string
		Protocol       string
		Host           string
		Port           uint
		Database       string
		Username       string
		Password       string
		Options        string
		MaxConnections int `yaml:"maxConnections"`
	}
	Redis struct {
		Host     string
		Port     uint
		Password string
		Database int
	}
	Bot struct {
		BotToken        string   `yaml:"botToken"`
		DiscordServerId string   `yaml:"discordServerId"`
		BotRole         string   `yaml:"botRole"`
		IgnoredRoles    []string `yaml:"ignoredRoles"`
	}
	OAuth struct {
		ClientId         string `yaml:"clientId"`
		ClientSecret     string `yaml:"clientSecret"`
		CallBackProtocol string `yaml:"callBackProtocol"`
		CallBackHost     string `yaml:"callBackHost"`
		CallBackUrl      string `yaml:"callBackUrl"`
	} `yaml:"oauth"`
	Net struct {
		ListenHost string `yaml:"listenHost"`
		ListenPort int    `yaml:"listenPort"`
	}
	Discord struct {
		InviteUrl string `yaml:"inviteUrl"`
	} `yaml:"discord"`
	Registry struct {
		Hostname         string `yaml:"hostname"`
		Port             int    `yaml:"port"`
		RegisterTTL      int    `yaml:"registerTtl"`
		RegisterInterval int    `yaml:"registerInterval"`
	} `yaml:"registry"`
	Inputs []string `yaml:"inputs"`
	Chat   struct {
		Slack struct {
			Debug bool   `yaml:"debug"`
			Token string `yaml:"token"`
		} `yaml:"slack"`
		Discord struct {
			Token     string   `yaml:"token"`
			WhiteList []string `yaml:"whiteList"`
			Prefix    string   `yaml:"prefix"`
		} `yaml:"discord"`
	} `yaml:"chat"`
	Extensions map[interface{}]interface{} `yaml:"extensions"`
}

func (c *Configuration) Load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("Could not read " + filename + " for configuration data.")
	}

	err = yaml.Unmarshal([]byte(data), c)
	if err != nil {
		return errors.New("Could not unmarshall " + filename + " as yaml")
	}

	// Let's set a default namespace because a lot of people don't care what it actually is
	if c.Namespace == "" {
		c.Namespace = "com.aba-eve"
	}

	c.initialized = true

	return nil
}

func (c *Configuration) IsInitialized() bool {
	return c.initialized
}

func (c *Configuration) LookupService(serviceType string, serviceName string) (serviceFullName string) {
	return fmt.Sprintf("%s.%s.%s", c.Namespace, serviceType, serviceName)
}
