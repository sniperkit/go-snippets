package registry

import (
	"strings"
	"sync"

	"github.com/cnf/stacker/engine"
)

// Registry represents a docker registry
type Registry struct {
	sync.Mutex
	Name          string
	Username      string
	Password      string
	Email         string
	ServerAddress string
}

// NewRegistry returns a new Registry object, with defaults
func NewRegistry(name string) *Registry {
	return &Registry{
		Name: name,
	}
}

// Command takes a Command object and triggers events
// TODO: Delete
func (r *Registry) Command(cmd *engine.Command) error {
	return nil
}

func nameFromServerAddress(url string) (name string, err error) {
	// TODO: Return default?
	stripped := url
	if strings.HasPrefix(url, "http://") {
		stripped = strings.Replace(url, "http://", "", 1)
	} else if strings.HasPrefix(url, "https://") {
		stripped = strings.Replace(url, "https://", "", 1)
	} else if strings.Contains(url, "://") {
		return "", ErrNotARegistryAddress
	}
	nameParts := strings.SplitN(stripped, "/", 2)
	return nameParts[0], nil
}

func nameFromPackageName(repo string) (name string, err error) {
	if strings.Contains(repo, "://") {
		return "", ErrNotAValidRepositoryName
	}
	nameParts := strings.SplitN(repo, "/", 2)
	if len(nameParts) == 1 || (!strings.Contains(nameParts[0], ".") && !strings.Contains(nameParts[0], ":") && nameParts[0] != "localhost") {
		// This is a Docker Index repos (ex: samalba/hipache or ubuntu)
		// TODO: sane response
		name, err := nameFromServerAddress(INDEXSERVER)
		if err != nil {
			return "", err
		}
		return name, nil
	}
	hostname := nameParts[0]
	// reposname := nameParts[1]
	return hostname, nil
}
