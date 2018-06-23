package server

import (
	"os"
	"os/signal"
	"syscall"

	"util/run"
)

//
func watchSignal(rch chan<- string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM:
			rch <- "syscall.SIGTERM"

		case syscall.SIGINT:
			rch <- "syscall.SIGINT"
		}

		close(rch)
		return
	}
}

//
func WatchSignal() <-chan string {
	c := make(chan string, 1)

	go run.Exec(true, watchSignal, c)

	return c
}
