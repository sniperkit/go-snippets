package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"registry-factory/lib"
	"syscall"
)

func main() {
	yamlFile := flag.String("c", "", "config yaml file")
	flag.Parse()

	//Load config
	if err := lib.Config.Load(*yamlFile); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := lib.NewProxyServer(ctx)
	done := make(chan error, 1)
	go func() {
		if err := s.Start(); err != nil {
			done <- err
		}
	}()

	log.Printf("Server is listening at %s:%d...\n", lib.Config.Host, lib.Config.Port)
	defer log.Println("Server is shutdown")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, os.Kill)

	select {
	case err := <-done:
		log.Fatalf("Server error: %s\n", err)
	case <-ctx.Done():
		log.Println("ctx done!")
	case <-sig:
		log.Println("Gracefully shutting down the server...")
		if err := s.Stop(); err != nil {
			log.Printf("Failed to shutdown server with error: %s\n", err)
		}
	}
}
