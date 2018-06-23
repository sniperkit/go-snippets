package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Action struct {
	Name     string   `yaml:"name"`
	Commands []string `yaml:"commands"`
}

type Actions map[string]Action

type config struct {
	Actions Actions `yaml:"actions"`
	Rules   *Rules  `yaml:"rules"`
}

func parseConfig(path string) (Actions, *Rules, error) {
	yamlData, err := loadYaml(path, filepath.Dir(path))
	if err != nil {
		return nil, nil, err
	}

	// uncomment for debug
	//fmt.Printf("%s", yamlData)

	config := &config{}

	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return nil, nil, fmt.Errorf("can't decode yaml: %s", err)
	}

	//read name from key of action directive
	for actionName, action := range config.Actions {
		action.Name = actionName
		config.Actions[actionName] = action
	}

	// it can be nil if the configuration file have not 'rules:' directive
	if config.Rules == nil {
		config.Rules = &Rules{}
	}

	// create regexp map, rule here accessible by reference
	for _, rule := range *config.Rules {
		rule.regexps = map[string]*regexp.Regexp{}
	}

	err = config.validate()
	if err != nil {
		return nil, nil, fmt.Errorf("invalid config: %s", err)
	}

	err = config.Rules.Compile()
	if err != nil {
		return nil, nil, err
	}

	return config.Actions, config.Rules, nil
}

func (config *config) validate() error {
	for _, rule := range *config.Rules {
		for _, action := range rule.Workflow {
			if _, ok := config.Actions[action]; !ok {
				return fmt.Errorf("undefined action in workflow: '%s'", action)
			}
		}
	}

	return nil
}
