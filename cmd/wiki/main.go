package main

import (
	"log"

	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/com/config/viper"
	"github.com/gliderlabs/stdcom/daemon"
	stdlog "github.com/gliderlabs/stdcom/log/std"
)

const AppName = "wiki"

func errFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// set up logger based on the standard library logger
	errFatal(stdlog.Register(com.DefaultRegistry))

	// load configuration using the Viper provider
	errFatal(config.Load(com.DefaultRegistry, viper.New(), AppName, []string{"."}))

	// start a daemon to run services
	errFatal(daemon.Run(com.DefaultRegistry, AppName))
}
