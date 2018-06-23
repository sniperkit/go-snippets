package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// function to display total execution time
func timer() func() {
	start := time.Now()
	return func() {
		log.Println("Total execution time:", time.Since(start))
	}
}

func main() {
	defer timer()()

	pause := make(chan bool)
	resume := make(chan bool)

	paused := false

	go func(pause chan bool, resume chan bool, paused *bool) {
		for {
			select {
			case <-pause:
				fmt.Println("Received pause signal. Press ^C to resume")
			case <-resume:
				fmt.Println("Resuming..")
			default:
			}
			if !*paused {
				fmt.Println("In loop.. Press ^C to pause or ^\\ to quit")
				time.Sleep(time.Second)
			}
		}
	}(pause, resume, &paused)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func(paused *bool) {
		for {
			select {
			case <-signalCh:
				if !*paused {
					*paused = true
					pause <- true
				} else {
					*paused = false
					resume <- true
				}
			}
		}
	}(&paused)

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGQUIT)
	<-exitSignal

	fmt.Println("Program exited.")
}
