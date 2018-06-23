package main

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/Ken2mer/knmr/logger"
)

func TestSignalHandler(t *testing.T) {
	termCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go signalHandler(c, termCh)

	resultCh := make(chan int)

	maxTerminatingInterval = 100 * time.Millisecond
	c <- os.Interrupt
	c <- os.Interrupt

	go func() {
		<-termCh
		<-termCh
		<-termCh
		<-termCh
		resultCh <- 0
	}()

	go func() {
		time.Sleep(time.Second)
		resultCh <- 1
	}()

	if r := <-resultCh; r != 0 {
		t.Errorf("Something went wrong")
	}
}

func TestLoop(t *testing.T) {
	if testing.Verbose() {
		logger.Verbose = true
	}

	termCh := make(chan struct{})
	exitCh := make(chan error)

	go func() {
		exitCh <- loop(termCh)
	}()

	time.Sleep(100 * time.Millisecond)

	termCh <- struct{}{}
	exitErr := <-exitCh
	if exitErr != nil {
		t.Errorf("exitErr should be nil, got: %s", exitErr)
	}
}
