package dislet

import (
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/validation"
	"github.com/maxpowel/goconfig"
    "os"
	"os/signal"
	"syscall"
	"github.com/fatih/color"
	"flag"

)

type Kernel struct {
	event []OnKernelReady
	Config *goconfig.Config
	registry goldi.TypeRegistry
	Container *goldi.Container

}

type OnKernelReady func(k *Kernel)

func (k *Kernel) Subscribe(ready OnKernelReady) {
	k.event = append(k.event, ready)

}


func NewKernel(configPath, paramsPath string, bootstrapModules []func(k *Kernel)) *Kernel {
	k := Kernel{}

	//k.event = make(chan int)
	k.registry = goldi.NewTypeRegistry()
	k.Container = goldi.NewContainer(k.registry, nil)
	k.Container.RegisterType("config", goconfig.NewConfig, configPath, paramsPath)

	//TODO hacer defer db.Close()
	// Config must be created before module bootstraping
	k.Config = k.Container.MustGet("config").(*goconfig.Config)
	for _, f := range bootstrapModules {
		f(&k)
	}

	// Check that container is ok
	validator := validation.NewContainerValidator()
	validator.MustValidate(k.Container)
	// Load configuration
	k.Config.Load()

	k.Container.Config = k.Config.Mapping

	// On kernel ready event

	for _, f := range k.event {
		f(&k)
	}


	return &k
}


func Daemonize(){
	// Daemonize
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		color.Yellow("Interrupt \"%v\" received, exiting", sig)

		done <- true
	}()

	<-done

}


func Boot(bootstrapModules []func(k *Kernel)) {
	color.Green("Starting...")
	// Parse parameters
	configPtr := flag.String("config", "config.yml", "Configuration file")
	parametersPtr := flag.String("parameters", "parameters.yml", "Parameters file")
	flag.Parse()
	// Dependency injection container
	//f := []func(k *dislet.Kernel){apiRestBootstrap}
	/*kernel := */NewKernel(*configPtr, *parametersPtr, bootstrapModules)

	Daemonize()

}
