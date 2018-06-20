package config

import (
	"errors"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/prometheus/common/log"
	"strconv"
)

type InitFunc func(configuration *Configuration) error

var svc micro.Service
var confFile string

// Builds and inits a new micro.Service object for use.  The initFunc functor being asked for will be inserted
// into the services options as a BeforeStart which will be called DURING the service.Run invocation but BEFORE the
// service is fully up and operational.  All of your initialization code that you need should go into this initFunc.
// If you don't need init code then feel free to use the NilInit function exported out of this package.
func NewService(version, serviceType string, serviceName string, initFunc InitFunc) micro.Service {
	service := micro.NewService(
		micro.Version(version),
		micro.BeforeStart(
			func() error {
				conf := Configuration{}

				conf.Load(confFile)

				if !conf.initialized {
					err := errors.New("Configuration not initialized, check your yaml format")
					log.Error(err)
					return err
				}

				if serviceType == "" {
					err := errors.New("serviceType is required")
					log.Error(err)
					return err
				}

				if serviceName == "" {
					err := errors.New("serviceName is required")
					log.Error(err)
					return err
				}

				svc.Init(micro.Name(conf.LookupService(serviceType, serviceName)))
				svc.Init(
					micro.Registry(
						registry.NewRegistry(
							registry.Addrs(
								conf.Registry.Hostname + ":" + strconv.Itoa(conf.Registry.Port),
							),
						),
					),
				)

				return initFunc(&conf)
			},
		),
		micro.Flags(
			cli.StringFlag{
				Name:        "configuration_file",
				Usage:       "The yaml configuration file for the service being loaded",
				Value:       "/etc/auth-srv/application.yaml",
				EnvVar:      "CONFIGURATION_FILE",
				Destination: &confFile,
			},
		),
	)

	service.Options()

	service.Init()

	svc = service

	return service
}

func NilInit(conf *Configuration) error {
	return nil
}
