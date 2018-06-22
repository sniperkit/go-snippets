package main

import (
	"net/http"
	"net/http/httputil"

	"github.com/Ken2mer/knmr/logger"
	"github.com/urfave/cli"
)

var godocCommand = cli.Command{
	Name:   "godoc",
	Usage:  "godoc",
	Action: godocCmd,
}

func godocCmd(clicontext *cli.Context) {
	resp, err := http.Get("http://localhost:6060/")
	if err != nil {
		logger.Errorf("http.Get() error: %v", err)
	}
	defer resp.Body.Close()

	// date := resp.Header.Get("Date")
	// content := resp.Header.Get("Content-Type")

	// logger.Debugf("\n%s", date)
	// logger.Debugf("\n%s", content)

	dump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		logger.Error(err)
	}
	logger.Debugf("\n%s", dump)
}
