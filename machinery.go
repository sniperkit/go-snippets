package machineryutil

import (
	machinery "github.com/RichardKnop/machinery/v1"
	machineryConfig "github.com/RichardKnop/machinery/v1/config"
)

// MachineryConfig loads the config for machinery
func MachineryConfig(filename string) (machineryConfig.Config, error) {
	var cnf machineryConfig.Config

	bytes, err := machineryConfig.ReadFromFile(filename)

	if err != nil {
		return cnf, err
	}

	err = machineryConfig.ParseYAMLConfig(&bytes, &cnf)

	if err != nil {
		return cnf, err
	}

	return cnf, nil
}

// MachineryServer configures the server for machinery
func MachineryServer(cnf machineryConfig.Config) (*machinery.Server, error) {
	server, err := machinery.NewServer(&cnf)

	if err != nil {
		return nil, err
	}

	return server, err
}
