package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	<-time.Tick(time.Second * 10000)
}

func init() {
	sich := make(chan os.Signal)
	signal.Notify(sich, syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGINT)

	go func() {
		for s := range sich {
			fmt.Println("Received signale")
			fmt.Printf("%#v\n", s)
			fmt.Printf("%s\n", s)
			os.Exit(0)
		}
	}()

}
