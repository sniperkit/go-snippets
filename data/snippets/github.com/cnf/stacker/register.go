package main

import (
	"github.com/golang/glog"

	"github.com/cnf/stacker/docker"
	"github.com/cnf/stacker/docker/container"
	"github.com/cnf/stacker/docker/image"
	"github.com/cnf/stacker/docker/registry"
	"github.com/cnf/stacker/modules/consul"
	"github.com/cnf/stacker/modules/cron"
	"github.com/cnf/stacker/modules/file"
	"github.com/cnf/stacker/modules/logger"
)

// registerAll modules with the engine
func registerAll() error {
	glog.Info("registering all modules")
	// Config Parsers
	file.RegisterConfigDecoder()
	registry.RegisterConfigDecoder()
	image.RegisterConfigDecoder()
	container.RegisterConfigDecoder()
	docker.RegisterConfigDecoder()
	consul.RegisterConfigDecoder()
	logger.RegisterConfigDecoder()

	// Docker
	docker.Register()

	// ConfStores
	consul.RegisterConfStore()
	file.RegisterConfStore()

	// Actions
	cron.RegisterAction()

	// Reactions
	consul.RegisterReaction()
	logger.RegisterReaction()
	return nil
}
