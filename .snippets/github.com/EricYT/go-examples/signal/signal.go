package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))

	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a test log entry")

	log.Println("Catch kill signal")
	fmt.Println("Wait a kill signal")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	//signal.Notify(interrupt, syscall.SIGTERM)
	//signal.Notify(interrupt, syscall.SIGHUP)

	select {
	case s := <-interrupt:
		log.Println("get kill signal:", time.Now())
		log.Printf("get singal:%+v", s)
		time.Sleep(time.Second * 10)
		log.Println("stop time:", time.Now())
	}
}
