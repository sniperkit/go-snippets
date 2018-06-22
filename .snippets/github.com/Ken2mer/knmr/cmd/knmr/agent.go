package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ken2mer/knmr/logger"
	"github.com/urfave/cli"
)

var agentCommand = cli.Command{
	Name:   "agent",
	Usage:  "agent",
	Action: agent,
}

var maxTerminatingInterval = 10 * time.Second

func signalHandler(c chan os.Signal, termCh chan<- struct{}) {
	received := false
	for sig := range c {
		if sig == syscall.SIGHUP {
			logger.Debugf("Received signal '%v'", sig)
		} else {
			if !received {
				received = true
				logger.Infof("Received signal '%v', try graceful shutdown up to %f seconds.",
					sig,
					maxTerminatingInterval.Seconds(),
				)
			} else {
				logger.Infof("Received signal '%v' again, force shutdown.", sig)
			}
			termCh <- struct{}{}
			go func() {
				time.Sleep(maxTerminatingInterval)
				logger.Infof("Timed out. force shutdown.")
				termCh <- struct{}{}
			}()
		}
	}
}

type loopState uint8

const (
	loopStateFirst loopState = iota
	loopStateTerminating
)

func loop(termCh <-chan struct{}) error {
	lState := loopStateFirst

	for {
		initialDelay := 3
		logger.Debugf("wait %d seconds before initial posting.", initialDelay)
		select {
		case <-termCh:
			logger.Debugf("received")
			lState = loopStateTerminating
		case <-time.After(time.Duration(initialDelay) * time.Second):
			logger.Debugf("time.After: %v", time.Duration(initialDelay)*time.Second)
		}

		switch lState {
		case loopStateFirst:
			logger.Debugf("continue...")
		case loopStateTerminating:
			logger.Debugf("terminating...")
			return nil
		}
	}
}

func agent(clicontext *cli.Context) error {
	termCh := make(chan struct{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go signalHandler(c, termCh)

	if err := loop(termCh); err != nil {
		return err
	}
	return nil
}
