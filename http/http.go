package http

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(f func() error) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit
	log.Println("Graceful Shutdown Server ...")
	if err := f(); err != nil {
		log.Fatal("Server Shutdown error:", err)
	}
	log.Println("Server exiting")
}
