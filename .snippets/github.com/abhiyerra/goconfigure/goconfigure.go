package goconfigure

import (
	"errors"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
)

type ParseType int

const (
	EnvParser ParseType = iota
	EtcdParser
)

type Etcd struct {
	Client    *etcd.Client
	Namespace string
}

type Config struct {
	Etcd      Etcd
	ParseType ParseType
	keys      map[string]Key
}

type Key struct {
	Description  string
	DefaultValue string
	EnvKey       string
	EtcdKey      string
}

func (c *Config) SetParser(parseType ParseType) {
	c.ParseType = parseType
}

func (c *Config) Add(name, desc, defaultValue, envKey, etcdKey string) {
	// Initialize the keys map if it isn't set.
	if c.keys == nil {
		c.keys = make(map[string]Key)
	}

	c.keys[name] = Key{
		Description:  desc,
		DefaultValue: defaultValue,
		EnvKey:       envKey,
		EtcdKey:      etcdKey,
	}
}

func (c *Config) Get(key string) (string, error) {
	k := c.keys[key]

	switch c.ParseType {
	case EnvParser:
		value := os.Getenv(k.EnvKey)
		if value == "" {
			return k.DefaultValue, errors.New("Key is empty")
		}
	case EtcdParser:
		resp, err := c.Etcd.Client.Get(
			fmt.Sprintf("/%s/%s", c.Etcd.Namespace, k),
			false,
			false)
		if err != nil {
			return k.DefaultValue, err
		}

		return resp.Node.Value, nil
	}

	return "", errors.New("Invalid Parser")
}
