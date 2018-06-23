package config

import (
	"errors"
	"fmt"
)

func (c Configuration) NewConnectionString() (string, error) {
	if !c.initialized {
		return "", errors.New("Configuration not initialized, call Load() before calling this.")
	}

	return c.Database.Username +
		":" +
		c.Database.Password +
		"@" +
		c.Database.Protocol +
		"(" +
		c.Database.Host +
		":" +
		fmt.Sprintf("%d", c.Database.Port) +
		")/" +
		c.Database.Database +
		"?" +
		c.Database.Options, nil
}
