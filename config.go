package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	extractAddress     string
	dbConnectionString string
	logPath            string
}

func NewConfig() *Config {
	cfg := new(Config)

	cfg.extractAddress = getenvDefault("EXTRACTOR_EXTRACT_ADDR", ":3096")
	cfg.dbConnectionString = getenvDefault("EXTRACTOR_DB_CONNECTION_STRING", "host=/var/run/postgresql dbname=entity-extractor_development sslmode=disable")
	cfg.logPath = getenvDefault("EXTRACTOR_LOG_PATH", "STDERR")

	flag.Usage = usage
	flag.Parse()

	return cfg
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
	helpstring := `
The following environment variables and defaults are available:

EXTRACTOR_EXTRACT_ADDR
  - Address on which to serve extraction requests
  - Default: ':3096'

EXTRACTOR_DB_CONNECTION_STRING
  - a postgresql connection string[1] used to connect to the database.
    Entites are read at startup from this database.
  - Default: 'host=/var/run/postgresql dbname=entity-extractor_development sslmode=disable'
  [1] http://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters

EXTRACTOR_ERROR_LOG
  - File to log errors to (in JSON format)
  - Default: 'STDERR'
`
	fmt.Fprintf(os.Stderr, helpstring)
	os.Exit(2)
}

func getenvDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultVal
	}
	return val
}
