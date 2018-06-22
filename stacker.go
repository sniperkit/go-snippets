package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	// "time"

	"github.com/golang/glog"

	"github.com/cnf/stacker/engine"
)

var configFile = flag.String("config_file", "/etc/stacker/stacker.toml", "path to our config file")

func main() {
	flag.Parse()

	glog.Infof("started with pid %d", os.Getpid())

	// register all modules
	if err := registerAll(); err != nil {
		glog.Fatal(err)
	}

	// Build up the engine
	eng, err := engine.New()
	if err != nil {
		glog.Fatal(err)
	}
	eng.SetConfigFile(*configFile)
	eng.Initialize()
	if err := eng.Start(); err != nil {
		glog.Fatal(err)
	}
	// eng.State()

	// respond to signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGTERM, os.Interrupt)
	go func() {
		for sig := range sigc {
			switch sig {
			case syscall.SIGHUP:
				glog.Infoln("SIGHUP: Reloading config")
				if err := eng.Restart(); err != nil {
					glog.Fatal(err)
				}
				glog.Flush()
			case os.Interrupt:
				glog.Fatal("SIGINT: Getting killed")
				glog.Flush()
			case syscall.SIGTERM:
				glog.Fatal("SIGTERM: Getting killed")
				glog.Flush()
			}
		}
	}()

	select {}
}
