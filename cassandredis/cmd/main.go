package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gophergala/cassandredis"
)

var (
	flAddr = flag.String("addr", "0.0.0.0:8765", "The listen address")
)

func printUsage() {
	fmt.Printf("Usage: cassandredis [--addr <addr>] <cassandra hosts> <keyspace>\n")
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		printUsage()
		os.Exit(1)
	}

	hosts := flag.Arg(0)
	keyspace := flag.Arg(1)

	server, err := cassandredis.NewServer(*flAddr, hosts, keyspace)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	if err := server.BootstrapMetadata(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println(server.Run())
}
