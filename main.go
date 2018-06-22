package main

import (
	"github.com/alext/tablecloth"
	"github.com/alphagov/entity-extractor/logger"
	"log"
	"os"
	"runtime"
)

var (
	cfg    *Config
	errlog logger.Logger
)

func logInfo(msg ...interface{}) {
	log.Println(msg...)
}

func SetGoMaxProcs() {
	if os.Getenv("GOMAXPROCS") == "" {
		// Use all available cores if not otherwise specified
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	logInfo("using GOMAXPROCS value of", runtime.GOMAXPROCS(0))
}

func SetupTablecloth() {
	// Set working dir for tablecloth if available. This is to allow restarts to
	// pick up new versions.
	// See http://godoc.org/github.com/alext/tablecloth#pkg-variables for details
	if wd := os.Getenv("GOVUK_APP_ROOT"); wd != "" {
		tablecloth.WorkingDir = wd
	}
}

func SetupLoggers(cfg *Config) {
	var err error
	errlog, err = logger.New(cfg.logPath)
	if err != nil {
		log.Fatal(err)
	}
	logInfo("logging JSON to", cfg.logPath)
}

func main() {
	cfg = NewConfig()
	SetupTablecloth()
	SetupLoggers(cfg)
	SetGoMaxProcs()

	extractor := NewExtractor(cfg)

	err := extractor.LoadEntities()
	if err != nil {
		log.Fatal(err)
	}

	extractorApi := NewExtractorAPI(extractor)

	logInfo("listening for requests on", cfg.extractAddress)
	err = tablecloth.ListenAndServe(cfg.extractAddress, extractorApi)
	if err != nil {
		log.Fatal(err)
	}
}
