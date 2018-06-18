package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Server struct {
	wg   *sync.WaitGroup
	quit chan struct{}
}

func OpenServer() *Server {

	srv := &Server{
		wg:   new(sync.WaitGroup),
		quit: make(chan struct{}),
	}

	srv.wg.Add(2)
	go loopServe(srv.wg, srv.quit, "Ping")
	go loopServe(srv.wg, srv.quit, "Pong")

	return srv
}

func (srv *Server) Close() error {

	if srv.quit == nil {
		return errors.New("server is closed")
	}

	close(srv.quit)
	srv.quit = nil

	srv.wg.Wait()

	return nil
}

func loopServe(wg *sync.WaitGroup, quit <-chan struct{}, m string) {

	defer wg.Done()

	defer func() {
		fmt.Println(m + " final task")
		time.Sleep(time.Second) // some task
		fmt.Println(m + " done")
	}()

	ticker := time.NewTicker(time.Millisecond * 1000)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return

		case <-ticker.C:
			fmt.Println(m)
		}
	}
}

func main() {
	srv := OpenServer()
	defer srv.Close()
	time.Sleep(5 * time.Second)
}
