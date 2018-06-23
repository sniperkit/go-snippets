package registry

import (
	dockerapi "github.com/fsouza/go-dockerclient"
)

// GetAuthConfiguration returns an auth object for the registry.
func GetAuthConfiguration(pkg string) *dockerapi.AuthConfiguration {
	name, err := nameFromPackageName(pkg)
	if err != nil {
		return &dockerapi.AuthConfiguration{}
	}
	// ErrNotARegistryAddress
	reg := list.getByName(name)
	return &dockerapi.AuthConfiguration{
		Username:      reg.Username,
		Password:      reg.Password,
		Email:         reg.Email,
		ServerAddress: reg.ServerAddress,
	}
}
